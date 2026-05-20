package ziputil

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/v2/pathutil"
)

// ZipManager provides zip and unzip operations using Go stdlib.
type ZipManager struct {
	pathChecker pathutil.PathChecker
	osProxy     OsProxy
}

// NewZipManager creates a ZipManager backed by the real OS.
func NewZipManager(pathChecker pathutil.PathChecker) *ZipManager {
	return &ZipManager{pathChecker: pathChecker, osProxy: RealOS{}}
}

// ZipDir zips sourceDirPth into destinationZipPth.
// When isContentOnly is false the archive contains basename/...; when true it contains the directory's contents directly.
func (z *ZipManager) ZipDir(sourceDirPth, destinationZipPth string, isContentOnly bool) error {
	if exist, err := z.pathChecker.IsDirExists(sourceDirPth); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("dir (%s) not exist", sourceDirPth)
	}

	baseDir := filepath.Dir(sourceDirPth)
	if isContentOnly {
		baseDir = sourceDirPth
	}

	return z.createZipFromDir(destinationZipPth, sourceDirPth, baseDir)
}

// ZipDirs zips multiple directories into a single archive.
// Each directory appears under its own basename in the archive.
func (z *ZipManager) ZipDirs(sourceDirPths []string, destinationZipPth string) error {
	for _, path := range sourceDirPths {
		if exist, err := z.pathChecker.IsDirExists(path); err != nil {
			return err
		} else if !exist {
			return fmt.Errorf("directory (%s) not exist", path)
		}
	}

	dest, err := z.osProxy.Create(destinationZipPth)
	if err != nil {
		return err
	}
	defer dest.Close() //nolint:errcheck

	zw := zip.NewWriter(dest)
	defer zw.Close() //nolint:errcheck

	for _, sourceDirPth := range sourceDirPths {
		if err := z.addDirToZip(zw, sourceDirPth, filepath.Dir(sourceDirPth)); err != nil {
			return err
		}
	}
	return nil
}

// ZipFile zips a single file into destinationZipPth.
func (z *ZipManager) ZipFile(sourceFilePth, destinationZipPth string) error {
	return z.ZipFiles([]string{sourceFilePth}, destinationZipPth)
}

// ZipFiles zips multiple files into destinationZipPth without preserving directory structure.
// Returns an error if two source files share the same base name.
func (z *ZipManager) ZipFiles(sourceFilePths []string, destinationZipPth string) error {
	seen := make(map[string]bool)
	for _, path := range sourceFilePths {
		if exist, err := z.pathChecker.IsPathExists(path); err != nil {
			return err
		} else if !exist {
			return fmt.Errorf("file (%s) not exist", path)
		}
		baseName := filepath.Base(path)
		if seen[baseName] {
			return fmt.Errorf("duplicate file name %q: files with the same base name cannot be zipped together", baseName)
		}
		seen[baseName] = true
	}

	dest, err := z.osProxy.Create(destinationZipPth)
	if err != nil {
		return err
	}
	defer dest.Close() //nolint:errcheck

	zw := zip.NewWriter(dest)
	defer zw.Close() //nolint:errcheck

	for _, filePath := range sourceFilePths {
		if err := z.addFileToZip(zw, filePath, filepath.Base(filePath)); err != nil {
			return err
		}
	}
	return nil
}

// UnZip extracts the zip archive at zipPth into intoDir, restoring permissions and symlinks.
func (z *ZipManager) UnZip(zipPth, intoDir string) error {
	r, err := zip.OpenReader(zipPth)
	if err != nil {
		return err
	}
	defer r.Close() //nolint:errcheck

	for _, f := range r.File {
		if err := z.extractEntry(f, intoDir); err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipManager) createZipFromDir(destinationZipPth, sourceDirPth, baseDir string) error {
	dest, err := z.osProxy.Create(destinationZipPth)
	if err != nil {
		return err
	}
	defer dest.Close() //nolint:errcheck

	zw := zip.NewWriter(dest)
	defer zw.Close() //nolint:errcheck

	return z.addDirToZip(zw, sourceDirPth, baseDir)
}

func (z *ZipManager) addDirToZip(zw *zip.Writer, sourceDirPth, baseDir string) error {
	return filepath.WalkDir(sourceDirPth, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		if d.Type()&fs.ModeSymlink != 0 {
			return z.addSymlinkToZip(zw, path, relPath)
		}

		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			hdr := &zip.FileHeader{
				Name:   relPath + "/",
				Method: zip.Store,
			}
			hdr.SetMode(info.Mode())
			_, err = zw.CreateHeader(hdr)
			return err
		}

		return z.addFileToZip(zw, path, relPath)
	})
}

func (z *ZipManager) addSymlinkToZip(zw *zip.Writer, path, name string) error {
	target, err := z.osProxy.Readlink(path)
	if err != nil {
		return err
	}

	hdr := &zip.FileHeader{
		Name:   name,
		Method: zip.Store,
	}
	hdr.SetMode(os.ModeSymlink | 0777)

	w, err := zw.CreateHeader(hdr)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(target))
	return err
}

func (z *ZipManager) addFileToZip(zw *zip.Writer, path, name string) error {
	info, err := z.osProxy.Lstat(path)
	if err != nil {
		return err
	}

	hdr, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	hdr.Name = name
	hdr.Method = zip.Deflate

	src, err := z.osProxy.Open(path)
	if err != nil {
		return err
	}
	defer src.Close() //nolint:errcheck

	w, err := zw.CreateHeader(hdr)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, src)
	return err
}

func (z *ZipManager) extractEntry(f *zip.File, intoDir string) error {
	destPath, err := sanitizeExtractPath(f.Name, intoDir)
	if err != nil {
		return err
	}

	if f.Mode()&os.ModeSymlink != 0 {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close() //nolint:errcheck

		target, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		if err := z.osProxy.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		return z.osProxy.Symlink(string(target), destPath)
	}

	if f.FileInfo().IsDir() {
		return z.osProxy.MkdirAll(destPath, f.Mode().Perm())
	}

	if err := z.osProxy.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close() //nolint:errcheck

	dest, err := z.osProxy.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode().Perm())
	if err != nil {
		return err
	}
	defer dest.Close() //nolint:errcheck

	_, err = io.Copy(dest, rc)
	return err
}

// sanitizeExtractPath guards against zip-slip: entries whose name would escape intoDir are rejected.
func sanitizeExtractPath(name, destDir string) (string, error) {
	cleanDest := filepath.Clean(destDir)
	destPath := filepath.Join(cleanDest, name)
	cleanPath := filepath.Clean(destPath)
	sep := string(os.PathSeparator)
	if cleanPath != cleanDest && !strings.HasPrefix(cleanPath, cleanDest+sep) {
		return "", fmt.Errorf("illegal path in zip entry: %s", name)
	}
	return destPath, nil
}

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
// isContentOnly=false: archive contains the directory under its own basename (e.g. "mydir/file.txt").
// isContentOnly=true: archive contains the directory's contents directly (e.g. "file.txt").
// Directory modification times are preserved. Symlinks are stored as symlinks (not followed).
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

// ZipDirs zips multiple directories into a single archive, each under its own basename.
// When two entries in sourceDirPths share the same basename (e.g. "/a/shared" and "/b/shared"),
// their contents are merged: files unique to either directory are preserved, and conflicting
// files (same relative path) resolve in favour of the last directory in the list.
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
// Each file is stored under its base name only.
// If two source files share the same base name, an error is returned before the destination
// file is created.
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
// Entries with path-traversal components (e.g. "../../escape") are rejected and an error is
// returned immediately; no further entries are extracted after a traversal is detected.
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
				Name:     relPath + "/",
				Method:   zip.Store,
				Modified: info.ModTime(),
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
		targetStr := string(target)
		if filepath.IsAbs(targetStr) {
			return fmt.Errorf("symlink target %q is absolute", targetStr)
		}
		resolvedTarget := filepath.Clean(filepath.Join(filepath.Dir(destPath), targetStr))
		cleanDest := filepath.Clean(intoDir)
		sep := string(os.PathSeparator)
		if resolvedTarget != cleanDest && !strings.HasPrefix(resolvedTarget, cleanDest+sep) {
			return fmt.Errorf("symlink target %q escapes extraction directory", targetStr)
		}
		if err := z.osProxy.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		return z.osProxy.Symlink(targetStr, destPath)
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

// sanitizeExtractPath guards against zip-slip: entries whose cleaned path would escape destDir
// are rejected. Returns the cleaned safe path, or an error if the entry would escape destDir.
func sanitizeExtractPath(name, destDir string) (string, error) {
	cleanDest := filepath.Clean(destDir)
	destPath := filepath.Join(cleanDest, name)
	cleanPath := filepath.Clean(destPath)
	sep := string(os.PathSeparator)
	if cleanPath != cleanDest && !strings.HasPrefix(cleanPath, cleanDest+sep) {
		return "", fmt.Errorf("illegal path in zip entry: %s", name)
	}
	return cleanPath, nil
}

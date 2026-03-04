package fileutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

// CopyFile copies a single file from src to dst.
func (fm fileManager) CopyFile(src, dst string) error {
	srcDir := filepath.Dir(src)
	fsys := fm.osProxy.DirFS(srcDir)

	return fm.copyFileFS(fsys, filepath.Base(src), dst)
}

// CopyFileFS is the excerpt from fs.CopyFS that copies a single file from fs.FS to dst path.
func (fm fileManager) copyFileFS(fsys fs.FS, src, dst string) error {
	r, err := fsys.Open(src)
	if err != nil {
		return err
	}
	defer r.Close() // nolint:errcheck
	info, err := r.Stat()
	if err != nil {
		return err
	}
	w, err := fm.osProxy.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	defer w.Close() // nolint:errcheck
	if _, err := io.Copy(w, r); err != nil {
		return &fs.PathError{Op: "Copy", Path: dst, Err: err}
	}
	if err := w.Sync(); err != nil {
		return &fs.PathError{Op: "Sync", Path: dst, Err: err}
	}
	if err := fm.copyOwner(info, dst); err != nil {
		return &fs.PathError{Op: "copyOwner", Path: dst, Err: err}
	}
	if err := fm.copyMode(info, dst); err != nil {
		return &fs.PathError{Op: "copyMode", Path: dst, Err: err}
	}
	if err := fm.copyTimes(info, dst); err != nil {
		return &fs.PathError{Op: "copyTimes", Path: dst, Err: err}
	}

	return nil
}

// CopyDir is a convenience method for copying a directory from src to dst.
//
// A copy of os.CopyFS because it messes up permissions when copying files
// from fs.FS
//
// CopyFS copies the file system fsys into the directory dir,
// creating dir if necessary.
//
// Preserves permissions and ownership when possible.
//
// CopyFS will not overwrite existing files. If a file name in fsys
// already exists in the destination, CopyFS will return an error
// such that errors.Is(err, fs.ErrExist) will be true.
//
// Symbolic links in dir are followed.
//
// New files added to fsys (including if dir is a subdirectory of fsys)
// while CopyFS is running are not guaranteed to be copied.
//
// Copying stops at and returns the first error encountered.
// Note: symlinks are preserved during the copy operation
func (fm fileManager) CopyDir(src, dst string) error {
	fsys := fm.osProxy.DirFS(src)
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		newPath := filepath.Join(dst, path)
		info, err := d.Info()
		if err != nil {
			return err
		}

		// This is not exhausetive in the original implementation either.
		// nolint:exhaustive
		switch d.Type() {
		case os.ModeDir:
			if err := fm.osProxy.MkdirAll(newPath, 0777); err != nil {
				return err
			}
			if err := fm.copyOwner(info, newPath); err != nil {
				return err
			}
			if err := fm.copyMode(info, newPath); err != nil {
				return err
			}
			return fm.copyTimes(info, newPath)

		case os.ModeSymlink:
			srcPath := filepath.Join(src, path)
			target, err := fm.osProxy.Readlink(srcPath)
			if err != nil {
				return err
			}
			if err := fm.osProxy.Symlink(target, newPath); err != nil {
				return err
			}
			if err := fm.copyOwner(info, newPath); err != nil {
				return err
			}
			return fm.copyTimes(info, newPath)

		// "normal" file
		case 0:
			return fm.copyFileFS(fsys, path, newPath)

		default:
			return &os.PathError{Op: "CopyFS", Path: path, Err: os.ErrInvalid}
		}
	})
}

// lchown ...
func (fm fileManager) lchown(path string, uid, gid int) error {
	return fm.osProxy.Lchown(path, uid, gid)
}

// copyOwner invokes lchown to copy ownership from srcInfo to dstPath.
func (fm fileManager) copyOwner(srcInfo os.FileInfo, dstPath string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	stat, ok := srcInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("missing Stat_t for symlink %s", dstPath)
	}
	// os.Lchown affects the link itself when given the link path
	if err := fm.lchown(dstPath, int(stat.Uid), int(stat.Gid)); err != nil {
		return fmt.Errorf("lchown(symlink) %s: %w", dstPath, err)
	}
	return nil
}

// chtimes ...
func (fm fileManager) chtimes(path string, atime, mtime time.Time) error {
	return fm.osProxy.Chtimes(path, atime, mtime)
}

// copyTimes invokes chtimes to copy access and modification times from srcInfo to dstPath.
func (fm fileManager) copyTimes(srcInfo os.FileInfo, dstPath string) error {
	if runtime.GOOS == "windows" {
		// On Windows we only set mod time (atime setting optional)
		if err := fm.chtimes(dstPath, srcInfo.ModTime(), srcInfo.ModTime()); err != nil {
			// ignore or return depending on policy
			return fmt.Errorf("chtimes %s: %w", dstPath, err)
		}

	} else {
		if stat, ok := srcInfo.Sys().(*syscall.Stat_t); ok {
			// set times (for non-symlink paths we use os.chtimes)
			if srcInfo.Mode()&os.ModeSymlink == 0 {
				atime := atimeFromStat(stat)
				mtime := srcInfo.ModTime()
				if err := fm.chtimes(dstPath, atime, mtime); err != nil {
					return fmt.Errorf("chtimes %s: %w", dstPath, err)
				}
			}
		}
	}
	return nil
}

// chmod ...
func (fm fileManager) chmod(path string, mode os.FileMode) error {
	return fm.osProxy.Chmod(path, mode)
}

// copyMode invokes chmod to copy file mode from srcInfo to dstPath.
func (fm fileManager) copyMode(srcInfo os.FileInfo, dstPath string) error {
	return fm.chmod(dstPath, srcInfo.Mode())
}

package fileutil

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"testing/fstest"

	"github.com/bitrise-io/go-utils/v2/fileutil/mocks/osproxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	assert.NoError(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1", nil))
}

func TestCopyFile_WithOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	// Pre-create the destination file so overwrite is exercised
	require.NoError(t, os.WriteFile(filepath.Join(dstDir, "file1"), []byte("old"), 0644))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing with TRUNC flag (overwrite)
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	assert.NoError(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1", &CopyOptions{Overwrite: true}))
}

func TestCopyFile_GivenDstFileOpenFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(nil, os.ErrPermission).
		Once()

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1", nil), os.ErrPermission.Error())
}

// TestCopyFile_OwnershipPermissionErrorIsTolerated: only root may chown to
// another user, so a non-root process copying a file owned by someone else
// (e.g. a build artifact written by a root container) always gets EPERM from
// Lchown — after the content was copied. The copy must succeed owned by the
// caller instead of failing.
func TestCopyFile_OwnershipPermissionErrorIsTolerated(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Ownership copy is denied, the rest of the copy must still complete
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(os.ErrPermission).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	assert.NoError(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1", nil))
}

// Non-permission ownership errors still fail the copy.
func TestCopyFile_GivenDstFileOwnershipChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(os.ErrInvalid).Once()

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1", nil), os.ErrInvalid.Error())
}

// TestCopyFile_SourceOwnedByAnotherUser reproduces the failure on the real
// filesystem: /etc/passwd is root-owned and world-readable, so for a non-root
// test process the byte copy succeeds and the ownership copy gets EPERM.
func TestCopyFile_SourceOwnedByAnotherUser(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix ownership semantics")
	}
	if os.Geteuid() == 0 {
		t.Skip("root may chown freely; the regression needs a non-root copier")
	}

	src := "/etc/passwd"
	info, err := os.Stat(src)
	require.NoError(t, err)
	stat, ok := info.Sys().(*syscall.Stat_t)
	require.True(t, ok)
	require.NotEqual(t, uint32(os.Geteuid()), stat.Uid, "test premise: source must be owned by another user")

	dst := filepath.Join(t.TempDir(), "copy")
	require.NoError(t, NewFileManager().CopyFile(src, dst, nil),
		"a non-root copy of another user's file must succeed with best-effort ownership")

	copied, err := os.ReadFile(dst)
	require.NoError(t, err)
	original, err := os.ReadFile(src)
	require.NoError(t, err)
	assert.Equal(t, original, copied)
}

func TestCopyFile_GivenDstFileModeChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(os.ErrPermission).Once()

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1", nil), os.ErrPermission.Error())
}

func TestCopyFile_GivenDstFileTimesChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(os.ErrPermission).Once()

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1", nil), os.ErrPermission.Error())
}

func TestCopyFileFS(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	assert.NoError(t, sut.CopyFileFS(os.DirFS(srcDir), "file1", filepath.Join(dstDir, "file1"), nil))
}

func TestCopyFileFS_WithOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	// Pre-create the destination file so overwrite is exercised
	require.NoError(t, os.WriteFile(filepath.Join(dstDir, "file1"), []byte("old"), 0644))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect dst file open for writing with TRUNC flag (overwrite)
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	assert.NoError(t, sut.CopyFileFS(os.DirFS(srcDir), "file1", filepath.Join(dstDir, "file1"), &CopyOptions{Overwrite: true}))
}

func TestCopyFileFS_GivenInMemoryFS_SkipsOwnershipPreservation(t *testing.T) {
	dstDir := t.TempDir()

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// fstest.MapFS does not expose a *syscall.Stat_t, so ownership and times
	// cannot be preserved. Lchown and Chtimes must not be called, while the file
	// mode is still set.
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()

	srcFS := fstest.MapFS{"file1": {Data: []byte("hello")}}

	assert.NoError(t, sut.CopyFileFS(srcFS, "file1", filepath.Join(dstDir, "file1"), nil))
}

func TestCopyFileFS_GivenMissingSrc_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	assert.Error(t, sut.CopyFileFS(os.DirFS(srcDir), "missing", filepath.Join(dstDir, "missing"), nil))
}

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(dstDir, mock.Anything, mock.Anything).Return(nil).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	// Expect symlink copy expectations for link
	osProxy.EXPECT().Readlink(filepath.Join(srcDir, "link")).Return(linkTarget, nil).Once()
	osProxy.EXPECT().Symlink(linkTarget, filepath.Join(dstDir, "link")).Return(nil).Once()
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "link"), mock.Anything, mock.Anything).Return(nil).Once()

	assert.NoError(t, sut.CopyDir(srcDir, dstDir, nil))
}

func TestCopyDir_WithOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	// Pre-create the destination file so overwrite is exercised
	require.NoError(t, os.WriteFile(filepath.Join(dstDir, "file1"), []byte("old"), 0644))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(dstDir, mock.Anything, mock.Anything).Return(nil).Once()

	// Expect file copy with TRUNC flag (overwrite)
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0o777))).
		Once()
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	assert.NoError(t, sut.CopyDir(srcDir, dstDir, &CopyOptions{Overwrite: true}))
}

func TestCopyDir_GivenDstDirCreationFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(os.ErrPermission).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir, nil), os.ErrPermission.Error())
}

func TestCopyDir_GivenDstOwnershipChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	// permission errors are tolerated (best-effort ownership); others fail
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(os.ErrInvalid).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir, nil), os.ErrInvalid.Error())
}

func TestCopyDir_GivenDstModeChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(os.ErrPermission).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir, nil), os.ErrPermission.Error())
}

func TestCopyDir_GivenDstTimesChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(dstDir, mock.Anything, mock.Anything).Return(os.ErrPermission).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir, nil), os.ErrPermission.Error())
}

func TestCopyDir_GivenReadLinkFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(dstDir, mock.Anything, mock.Anything).Return(nil).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	// Expect symlink copy expectations for link
	osProxy.EXPECT().Readlink(filepath.Join(srcDir, "link")).Return("", os.ErrPermission).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir, nil), os.ErrPermission.Error())
}

func TestCopyDir_SymlinkFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(dstDir, mock.Anything, mock.Anything).Return(nil).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	// Expect symlink copy expectations for link
	osProxy.EXPECT().Readlink(filepath.Join(srcDir, "link")).Return(linkTarget, nil).Once()
	osProxy.EXPECT().Symlink(linkTarget, filepath.Join(dstDir, "link")).Return(os.ErrPermission).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir, nil), os.ErrPermission.Error())
}

func TestCopyDir_SymlinkLChownFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := osproxy.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(dstDir, mock.Anything, mock.Anything).Return(nil).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(filepath.Join(dstDir, "file1"), mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(nil).Once()

	// Expect symlink copy expectations for link
	osProxy.EXPECT().Readlink(filepath.Join(srcDir, "link")).Return(linkTarget, nil).Once()
	osProxy.EXPECT().Symlink(linkTarget, filepath.Join(dstDir, "link")).Return(nil).Once()
	// permission errors are tolerated (best-effort ownership); others fail
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "link"), mock.Anything, mock.Anything).Return(os.ErrInvalid).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir, nil), os.ErrInvalid.Error())
}

// ---------------------------
// Helpers
// ---------------------------

func createSrcDirWithFiles(t *testing.T, baseDir string, fileNames []string) string {
	t.Helper()
	srcDir := filepath.Join(baseDir, "src-dir")
	require.NoError(t, os.MkdirAll(srcDir, 0755))
	for _, name := range fileNames {
		sourceFile := filepath.Join(srcDir, name)
		require.NoError(t, os.WriteFile(sourceFile, []byte(name), 0755))
	}
	return srcDir
}

package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/v2/fileutil/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

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

	assert.NoError(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1"))
}

func TestCopyFile_GivenDstFileOpenFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(nil, os.ErrPermission).
		Once()

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1"), os.ErrPermission.Error())
}

func TestCopyFile_GivenDstFileOwnershipChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect srcDir check
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	// Expect dst file open for writing
	osProxy.EXPECT().
		OpenFile(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).
		Return(os.OpenFile(filepath.Join(dstDir, "file1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(0o777))).
		Once()

	// Expect dst file ownership, permissions and times to be set
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "file1"), mock.Anything, mock.Anything).Return(os.ErrPermission).Once()

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1"), os.ErrPermission.Error())
}

func TestCopyFile_GivenDstFileModeChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

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

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1"), os.ErrPermission.Error())
}

func TestCopyFile_GivenDstFileTimesChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

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

	assert.ErrorContains(t, sut.CopyFile(srcDir+"/file1", dstDir+"/file1"), os.ErrPermission.Error())
}

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := mocks.NewOsProxy(t)

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

	assert.NoError(t, sut.CopyDir(srcDir, dstDir))
}

func TestCopyDir_GivenDstDirCreationFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(os.ErrPermission).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir), os.ErrPermission.Error())
}

func TestCopyDir_GivenDstOwnershipChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(os.ErrPermission).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir), os.ErrPermission.Error())
}

func TestCopyDir_GivenDstModeChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(os.ErrPermission).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir), os.ErrPermission.Error())
}

func TestCopyDir_GivenDstTimesChangeFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))

	osProxy := mocks.NewOsProxy(t)

	sut := fileManager{osProxy: osProxy}

	// Expect changes for dstDir
	osProxy.EXPECT().MkdirAll(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Lchown(dstDir, mock.Anything, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chmod(dstDir, mock.Anything).Return(nil).Once()
	osProxy.EXPECT().Chtimes(dstDir, mock.Anything, mock.Anything).Return(os.ErrPermission).Once()

	// Expect file copy expectations for file1
	osProxy.EXPECT().DirFS(srcDir).Return(os.DirFS(srcDir)).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir), os.ErrPermission.Error())
}

func TestCopyDir_GivenReadLinkFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := mocks.NewOsProxy(t)

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

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir), os.ErrPermission.Error())
}

func TestCopyDir_SymlinkFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := mocks.NewOsProxy(t)

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

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir), os.ErrPermission.Error())
}

func TestCopyDir_SymlinkLChownFailure_WillFail(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := createSrcDirWithFiles(t, t.TempDir(), []string{"file1"})
	dstDir := filepath.Join(tmpDir, "dst-dir")
	assert.NoError(t, os.MkdirAll(dstDir, 0755))
	linkTarget := filepath.Join(srcDir, "file1")
	assert.NoError(t, os.Symlink(linkTarget, filepath.Join(srcDir, "link")))

	osProxy := mocks.NewOsProxy(t)

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
	osProxy.EXPECT().Lchown(filepath.Join(dstDir, "link"), mock.Anything, mock.Anything).Return(os.ErrPermission).Once()

	assert.ErrorContains(t, sut.CopyDir(srcDir, dstDir), os.ErrPermission.Error())
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

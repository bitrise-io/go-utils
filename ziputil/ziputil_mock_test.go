package ziputil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-io/go-utils/v2/ziputil/mocks/osproxy"
	"github.com/stretchr/testify/require"
)

var errBoom = errors.New("boom")

// managerWithMock builds a ZipManager whose OS boundary is mocked, while path existence checks
// still hit the real filesystem. This lets the tests drive OS error paths that cannot be
// triggered deterministically against the real OS.
func managerWithMock(m *osproxy.OsProxy) *ZipManager {
	return &ZipManager{pathChecker: pathutil.NewPathChecker(), osProxy: m}
}

// TestZipFilesCreateFailureSurfaces verifies that a failure creating the temporary archive is
// returned and that no cleanup is attempted (nothing was created yet).
func TestZipFilesCreateFailureSurfaces(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src.txt")
	require.NoError(t, os.WriteFile(src, []byte("x"), 0644))
	dest := filepath.Join(tmpDir, "out.zip")

	m := osproxy.NewOsProxy(t)
	m.EXPECT().Create(dest+".tmp").Return(nil, errBoom).Once()

	// No Remove/Rename expectations: the mock asserts neither is called.
	require.ErrorIs(t, managerWithMock(m).ZipFiles([]string{src}, dest), errBoom)
}

// TestZipFilesWriteFailureCleansUpTemp verifies that when writing an entry fails, the error is
// surfaced and the temporary archive is removed through the OS proxy.
func TestZipFilesWriteFailureCleansUpTemp(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src.txt")
	require.NoError(t, os.WriteFile(src, []byte("x"), 0644))
	dest := filepath.Join(tmpDir, "out.zip")
	tmpPath := dest + ".tmp"

	realTmp, err := os.Create(tmpPath)
	require.NoError(t, err)
	defer realTmp.Close() //nolint:errcheck

	m := osproxy.NewOsProxy(t)
	m.EXPECT().Create(tmpPath).Return(realTmp, nil).Once()
	m.EXPECT().Lstat(src).Return(nil, errBoom).Once() // addFileToZip fails before any rename
	m.EXPECT().Remove(tmpPath).Return(nil).Once()     // partial temp archive is cleaned up

	require.ErrorIs(t, managerWithMock(m).ZipFiles([]string{src}, dest), errBoom)
}

// TestZipFilesRenameFailureCleansUpTemp verifies that when the final rename onto the destination
// fails, the error is surfaced and the temporary archive is removed, leaving any pre-existing
// destination untouched.
func TestZipFilesRenameFailureCleansUpTemp(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src.txt")
	require.NoError(t, os.WriteFile(src, []byte("x"), 0644))
	dest := filepath.Join(tmpDir, "out.zip")
	tmpPath := dest + ".tmp"

	realTmp, err := os.Create(tmpPath)
	require.NoError(t, err)
	defer realTmp.Close() //nolint:errcheck

	info, err := os.Lstat(src)
	require.NoError(t, err)
	realSrc, err := os.Open(src)
	require.NoError(t, err)
	defer realSrc.Close() //nolint:errcheck

	m := osproxy.NewOsProxy(t)
	m.EXPECT().Create(tmpPath).Return(realTmp, nil).Once()
	m.EXPECT().Lstat(src).Return(info, nil).Once()
	m.EXPECT().Open(src).Return(realSrc, nil).Once()
	m.EXPECT().Rename(tmpPath, dest).Return(errBoom).Once()
	m.EXPECT().Remove(tmpPath).Return(nil).Once()

	require.ErrorIs(t, managerWithMock(m).ZipFiles([]string{src}, dest), errBoom)
}

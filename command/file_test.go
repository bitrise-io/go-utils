package command

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyFileErrorIfDirectory(t *testing.T) {
	t.Log("It fails if source is a directory")
	{
		cmd := exec.Command("mkdir", "source_dir")
		cmd.Run()
		err := CopyFile("source_dir", "./nothing/whatever")
		require.Error(t, err)
		cmd = exec.Command("rm", "-rf", "./source_dir/")
		cmd.Run()
	}
}

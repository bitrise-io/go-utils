package gems

import (
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/rubycommand"
)

// InstallBundlerCommand returns a command to install a specific bundler version
func InstallBundlerCommand(gemfileLockVersion Version) (*command.Model, error) {
	installBundlerCmdParams := []string{"gem", "install", "bundler", "--force", "--no-document"}
	if gemfileLockVersion.Found {
		installBundlerCmdParams = append(installBundlerCmdParams, []string{"-v", gemfileLockVersion.Version}...)
	}

	return command.NewFromSlice(installBundlerCmdParams)
}

// BundleInstallCommand returns a command to install a bundle using bundler
func BundleInstallCommand(gemfileLockVersion Version) (*command.Model, error) {
	bundleInstallCmdParams := []string{"bundle"}
	if gemfileLockVersion.Found {
		bundleInstallCmdParams = append(bundleInstallCmdParams, "_"+gemfileLockVersion.Version+"_")
	}
	bundleInstallCmdParams = append(bundleInstallCmdParams, []string{"install", "--jobs", "20", "--retry", "5"}...)

	return rubycommand.NewFromSlice(bundleInstallCmdParams)
}

// BundleExecPrefix returns a slice containing: "bundle [_verson_] exec"
func BundleExecPrefix(bundlerVersion Version) []string {
	bundleExec := []string{"bundle"}
	if bundlerVersion.Found {
		bundleExec = append(bundleExec, "_"+bundlerVersion.Version+"_")
	}
	return append(bundleExec, "exec")
}

// RbenvVersionsCommand retruns a command to print used and available ruby versions if rbenv is installed
func RbenvVersionsCommand() (*command.Model, error) {
	if _, err := command.New("which", "rbenv").RunAndReturnTrimmedCombinedOutput(); err != nil {
		return nil, err
	}

	return command.New("rbenv", "versions"), nil
}

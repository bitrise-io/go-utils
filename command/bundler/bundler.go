package bundler

import (
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/rubycommand"
)

// InstallBundlerCommand returns a command to install a specific bundler version
func InstallBundlerCommand(gemfileLockVersion GemVersion) (*command.Model, error) {
	installBundlerCmdParams := []string{"gem", "install", "bundler", "--force", "--no-document"}
	if gemfileLockVersion.found {
		installBundlerCmdParams = append(installBundlerCmdParams, []string{"-v", gemfileLockVersion.version}...)
	}

	return command.NewFromSlice(installBundlerCmdParams)
}

// BundleInstallCommand returns a command to install a bundle using bundler
func BundleInstallCommand(gemfileLockVersion GemVersion) (*command.Model, error) {
	bundleInstallCmdParams := []string{"bundle"}
	if gemfileLockVersion.found {
		bundleInstallCmdParams = append(bundleInstallCmdParams, "_"+gemfileLockVersion.version+"_")
	}
	bundleInstallCmdParams = append(bundleInstallCmdParams, []string{"install", "--jobs", "20", "--retry", "5"}...)

	return rubycommand.NewFromSlice(bundleInstallCmdParams)
}

// RbenvVersionsCommand retruns a command to print used and available ruby versions if rbenv is installed
func RbenvVersionsCommand() (*command.Model, error) {
	if _, err := command.New("which", "rbenv").RunAndReturnTrimmedCombinedOutput(); err != nil {
		return nil, err
	}

	return command.New("rbenv", "versions"), nil
}

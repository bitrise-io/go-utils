package rubycommand

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmdExist(t *testing.T) {
	t.Log("exist")
	{
		require.Equal(t, true, cmdExist("ls"))
	}

	t.Log("not exist")
	{
		require.Equal(t, false, cmdExist("__not_existing_command__"))
	}
}

func TestSudoNeeded(t *testing.T) {
	t.Log("sudo NOT need")
	{
		require.Equal(t, false, sudoNeeded(Unkown, "ls"))
		require.Equal(t, false, sudoNeeded(SystemRuby, "ls"))
		require.Equal(t, false, sudoNeeded(BrewRuby, "ls"))
		require.Equal(t, false, sudoNeeded(RVMRuby, "ls"))
		require.Equal(t, false, sudoNeeded(RbenvRuby, "ls"))
	}

	t.Log("sudo needed for SystemRuby in case of gem list management command")
	{
		require.Equal(t, false, sudoNeeded(Unkown, "gem", "install", "fastlane"))
		require.Equal(t, true, sudoNeeded(SystemRuby, "gem", "install", "fastlane"))
		require.Equal(t, false, sudoNeeded(BrewRuby, "gem", "install", "fastlane"))
		require.Equal(t, false, sudoNeeded(RVMRuby, "gem", "install", "fastlane"))
		require.Equal(t, false, sudoNeeded(RbenvRuby, "gem", "install", "fastlane"))

		require.Equal(t, false, sudoNeeded(Unkown, "gem", "uninstall", "fastlane"))
		require.Equal(t, true, sudoNeeded(SystemRuby, "gem", "uninstall", "fastlane"))
		require.Equal(t, false, sudoNeeded(BrewRuby, "gem", "uninstall", "fastlane"))
		require.Equal(t, false, sudoNeeded(RVMRuby, "gem", "uninstall", "fastlane"))
		require.Equal(t, false, sudoNeeded(RbenvRuby, "gem", "uninstall", "fastlane"))

		require.Equal(t, false, sudoNeeded(Unkown, "bundle", "install"))
		require.Equal(t, false, sudoNeeded(Unkown, "bundle", "_2.0.2_", "install"))
		require.Equal(t, true, sudoNeeded(SystemRuby, "bundle", "install"))
		require.Equal(t, true, sudoNeeded(SystemRuby, "bundle", "_2.0.2_", "install"))
		require.Equal(t, false, sudoNeeded(SystemRuby, "bundle", "_2.0.2_"))
		require.Equal(t, false, sudoNeeded(BrewRuby, "bundle", "install"))
		require.Equal(t, false, sudoNeeded(RVMRuby, "bundle", "install"))
		require.Equal(t, false, sudoNeeded(RbenvRuby, "bundle", "install"))

		require.Equal(t, false, sudoNeeded(Unkown, "bundle", "update"))
		require.Equal(t, false, sudoNeeded(Unkown, "bundle", "_2.0.2_", "update"))
		require.Equal(t, true, sudoNeeded(SystemRuby, "bundle", "update"))
		require.Equal(t, true, sudoNeeded(SystemRuby, "bundle", "_2.0.2_", "update"))
		require.Equal(t, false, sudoNeeded(BrewRuby, "bundle", "update"))
		require.Equal(t, false, sudoNeeded(RVMRuby, "bundle", "update"))
		require.Equal(t, false, sudoNeeded(RbenvRuby, "bundle", "update"))
	}
}

func TestFindGemInList(t *testing.T) {
	t.Log("finds gem")
	{
		gemList := `
*** LOCAL GEMS ***

addressable (2.5.0, 2.4.0, 2.3.8)
activesupport (5.0.0.1, 4.2.7.1, 4.2.6, 4.2.5, 4.1.16, 4.0.13)
angularjs-rails (1.5.8)`

		found, err := findGemInList(gemList, "activesupport", "")
		require.NoError(t, err)
		require.Equal(t, true, found)
	}

	t.Log("finds gem with version")
	{
		gemList := `
*** LOCAL GEMS ***

addressable (2.5.0, 2.4.0, 2.3.8)
activesupport (5.0.0.1, 4.2.7.1, 4.2.6, 4.2.5, 4.1.16, 4.0.13)
angularjs-rails (1.5.8)`

		found, err := findGemInList(gemList, "activesupport", "4.2.5")
		require.NoError(t, err)
		require.Equal(t, true, found)
	}

	t.Log("gem version not found in list")
	{
		gemList := `
*** LOCAL GEMS ***

addressable (2.5.0, 2.4.0, 2.3.8)
activesupport (5.0.0.1, 4.2.7.1, 4.2.6, 4.2.5, 4.1.16, 4.0.13)
angularjs-rails (1.5.8)`

		found, err := findGemInList(gemList, "activesupport", "0.9.0")
		require.NoError(t, err)
		require.Equal(t, false, found)
	}

	t.Log("gem not found in list")
	{
		gemList := `
*** LOCAL GEMS ***

addressable (2.5.0, 2.4.0, 2.3.8)
activesupport (5.0.0.1, 4.2.7.1, 4.2.6, 4.2.5, 4.1.16, 4.0.13)
angularjs-rails (1.5.8)`

		found, err := findGemInList(gemList, "fastlane", "")
		require.NoError(t, err)
		require.Equal(t, false, found)
	}

	t.Log("gem with version not found in list")
	{
		gemList := `
*** LOCAL GEMS ***

addressable (2.5.0, 2.4.0, 2.3.8)
activesupport (5.0.0.1, 4.2.7.1, 4.2.6, 4.2.5, 4.1.16, 4.0.13)
angularjs-rails (1.5.8)`

		found, err := findGemInList(gemList, "fastlane", "2.70")
		require.NoError(t, err)
		require.Equal(t, false, found)
	}
}

func Test_isSpecifiedRbenvRubyInstalled(t *testing.T) {

	t.Log("RBENV_VERSION installed -  2.3.5 (set by RBENV_VERSION environment variable)")
	{
		message := "2.3.5 (set by RBENV_VERSION environment variable)"
		installed, version, err := isSpecifiedRbenvRubyInstalled(message)
		require.NoError(t, err)
		require.Equal(t, true, installed)
		require.Equal(t, "2.3.5", version)
	}

	t.Log("RBENV_VERSION not installed - rbenv: version `2.34.0' is not installed (set by RBENV_VERSION environment variable)")
	{
		message := "rbenv: version `2.34.0' is not installed (set by RBENV_VERSION environment variable)"
		installed, version, err := isSpecifiedRbenvRubyInstalled(message)
		require.NoError(t, err)
		require.Equal(t, false, installed)
		require.Equal(t, "2.34.0", version)
	}

	t.Log("Global ruby installed - 2.3.5 (set by /Users/Vagrant/.rbenv/version)")
	{

		message := "2.3.5 (set by /Users/Vagrant/.rbenv/version)"
		installed, version, err := isSpecifiedRbenvRubyInstalled(message)
		require.NoError(t, err)
		require.Equal(t, true, installed)
		require.Equal(t, "2.3.5", version)
	}

	t.Log("Global ruby not installed - rbenv: version `2.4.2' is not installed (set by /Users/Vagrant/.rbenv/version)")
	{

		message := "rbenv: version `2.4.2' is not installed (set by /Users/Vagrant/.rbenv/version)"
		installed, version, err := isSpecifiedRbenvRubyInstalled(message)
		require.NoError(t, err)
		require.Equal(t, false, installed)
		require.Equal(t, "2.4.2", version)
	}

	t.Log(".ruby-version not installed - rbenv: version `2.89.2' is not installed (set by /Users/Vagrant/.ruby-version)")
	{

		message := "rbenv: version `2.89.2' is not installed (set by /Users/Vagrant/.ruby-version)"
		installed, version, err := isSpecifiedRbenvRubyInstalled(message)
		require.NoError(t, err)
		require.Equal(t, false, installed)
		require.Equal(t, "2.89.2", version)
	}

	t.Log(".ruby-version installed 2.3.5 (set by /Users/Vagrant/.ruby-version)")
	{

		message := "2.3.5 (set by /Users/Vagrant/.ruby-version)"
		installed, version, err := isSpecifiedRbenvRubyInstalled(message)
		require.NoError(t, err)
		require.Equal(t, true, installed)
		require.Equal(t, "2.3.5", version)
	}
}

package colorstring

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddColor(t *testing.T) {
	/*
	  blackColor   Color = "\x1b[30;1m"
	  resetColor   Color = "\x1b[0m"
	*/

	t.Log("colored_string = color + string + reset_color")
	{
		desiredColored := "\x1b[30;1m" + "test" + "\x1b[0m"
		colored := AddColor(blackColor, "test")
		require.Equal(t, desiredColored, colored)
	}
}

func TestBlack(t *testing.T) {
	t.Log("Simple string can be blacked")
	{
		desiredColored := "\x1b[30;1m" + "test" + "\x1b[0m"
		colored := Black("test")
		require.Equal(t, desiredColored, colored)
	}

	t.Log("Multiple strings can be blacked")
	{
		desiredColored := "\x1b[30;1m" + "Hello Bitrise !" + "\x1b[0m"
		colored := Black("Hello ", "Bitrise ", "!")
		require.Equal(t, desiredColored, colored)
	}
}

func TestBlackf(t *testing.T) {
	t.Log("Simple format can be blacked")
	{
		desiredColored := "\x1b[30;1m" + fmt.Sprintf("Hello %s", "bitrise") + "\x1b[0m"
		colored := Blackf("Hello %s", "bitrise")
		require.Equal(t, desiredColored, colored)
	}

	t.Log("Complex format can be blacked")
	{
		desiredColored := "\x1b[30;1m" + fmt.Sprintf("Hello %s %s", "bitrise", "!") + "\x1b[0m"
		colored := Blackf("Hello %s %s", "bitrise", "!")
		require.Equal(t, desiredColored, colored)
	}
}

func TestGitError(t *testing.T) {
	fmt.Printf(`
Error:
  fetch failed:
    remote: Support for password authentication was removed on August 13, 2021.
    remote: Please see https://docs.github.com/en/get-started/getting-started-with-git/about-remote-repositories#cloning-with-https-urls for information on currently recommended modes of authentication.
    fatal: Authentication failed for 'https://github.com/Intan90002/empty.git/'
`)

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf(`
%s:
  fetch failed:
    remote: Support for password authentication was removed on August 13, 2021.
    remote: Please see https://docs.github.com/en/get-started/getting-started-with-git/about-remote-repositories#cloning-with-https-urls for information on currently recommended modes of authentication.
	fatal: %s
`,
		Red("Error"),
		Redf("Authentication failed for '%s'", Cyan("https://github.com/Intan90002/empty.git/")),
	)
}

func TestXcodeError(t *testing.T) {
	fmt.Printf(`
Command failed with exit status 70 (xcodebuild "-exportArchive" "-archivePath" "./code-sign-test.xcarchive" "-exportPath" "./exported" "-exportOptionsPlist" "./export_options.plist"):
  error: exportArchive: "share-extension.appex" requires a provisioning profile.
  error: exportArchive: "code-sign-test.app" requires a provisioning profile.
  error: exportArchive: "watchkit-app.app" requires a provisioning profile.
  error: exportArchive: "watchkit-app Extension.appex" requires a provisioning profile.
`)

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf(`
Command failed with exit status 70 (%s):
  error: exportArchive: %s
  error: exportArchive: %s
  error: exportArchive: %s
  error: exportArchive: %s
`,
		Cyan(`xcodebuild "-exportArchive" "-archivePath" "./code-sign-test.xcarchive" "-exportPath" "./exported" "-exportOptionsPlist" "./export_options.plist"`),
		Red(`"share-extension.appex" requires a provisioning profile.`),
		Red(`"code-sign-test.app" requires a provisioning profile.`),
		Red(`"watchkit-app.app" requires a provisioning profile.`),
		Red(`"watchkit-app Extension.appex" requires a provisioning profile.`),
	)
}

package pathutil_test

import (
	"fmt"
	"strings"
	"testing/fstest"

	"github.com/bitrise-io/go-utils/pathutil"
)

func Examplev2() {
	// root := "/usr/local/go/bin"
	// fileSystem := os.DirFS(root)

	fsys := fstest.MapFS{
		"/app2.apk":       {},
		"dir/MyApp.ipa":   {},
		"/app1.apk":       {},
		"dir/file.go":     {},
		"dir/subdir/x.go": {},
	}

	pths, err := pathutil.FilterPathsV2(
		fsys,
		pathutil.ExtensionFilter(".apk", true),
	)
	if err != nil {
		panic(err)
	}

	// pths, err = pathutil.SortPathsByComponents(pths)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(strings.Join(pths, ", "))
	// Output: /app1.apk, /app2.apk
}

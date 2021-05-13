package pathutil_test

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

func Example() {
	apkFilter := func(s string) (bool, error) {
		return filepath.Ext(s) == ".apk", nil
	}
	pths, err := pathutil.FilterPaths([]string{"/app2.apk", "/MyApp.ipa", "/app1.apk"}, apkFilter)
	if err != nil {
		panic(err)
	}

	pths, err = pathutil.SortPathsByComponents(pths)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(pths, ", "))
	// Output: /app1.apk, /app2.apk
}

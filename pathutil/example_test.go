package pathutil_test

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

func Example() {
	pths, err := pathutil.FilterPaths(
		[]string{"/app2.apk", "/MyApp.ipa", "/app1.apk"},
		pathutil.ExtensionFilter(".apk", true),
	)
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

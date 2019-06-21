package gems

import (
	"fmt"
	"regexp"
	"strings"
)

// Version contains bundler or fastlane version
type Version struct {
	Version string
	Found   bool
}

// ParseVersionFromBundle  returns the specified grm version parsed from a Gemfile.lock on a best effort basis, for logging purposes only.
//
// for "fastlane" and the following Gemfile.lock example, it returns: ">= 2.0)"
//   specs:
//     CFPropertyList (3.0.0)
//     addressable (2.6.0)
//       public_suffix (>= 2.0.2, < 4.0)
//     atomos (0.1.3)
//     babosa (1.0.2)
//     badge (0.8.5)
//       curb (~> 0.9)
//       fastimage (>= 1.6)
//       fastlane (>= 2.0)
//       mini_magick (>= 4.5)
//     claide (1.0.2)
func ParseVersionFromBundle(gemName string, gemfileLockContent string) (gemVersion Version, err error) {
	relevantLines := []string{}
	lines := strings.Split(gemfileLockContent, "\n")

	specsStart := false
	for _, line := range lines {
		trimmed := strings.Trim(line, " ")
		if trimmed == "" {
			specsStart = false
		}

		if strings.Contains(line, "specs:") {
			specsStart = true
			continue
		}

		if specsStart {
			relevantLines = append(relevantLines, line)
		}
	}

	//     fastlane (1.109.0)
	exp := regexp.MustCompile(fmt.Sprintf(`^%s \((.+)\)`, regexp.QuoteMeta(gemName)))
	for _, line := range relevantLines {
		match := exp.FindStringSubmatch(strings.TrimSpace(line))
		if match == nil {
			continue
		}
		if len(match) != 2 {
			return Version{}, fmt.Errorf("unexpected regexp match: %v", match)
		}
		return Version{
			Version: match[1],
			Found:   true,
		}, nil
	}

	return Version{}, nil
}

// ParseBundlerVersion returns the bundler version used to create the bundle
func ParseBundlerVersion(gemfileLockContent string) (gemVersion Version, err error) {
	/*
		BUNDLED WITH
			1.17.1
	*/
	bundlerRegexp := regexp.MustCompile(`(?m)^BUNDLED WITH\n\s+(\S+)`)
	match := bundlerRegexp.FindStringSubmatch(gemfileLockContent)
	if match == nil {
		return Version{}, nil
	}
	if len(match) != 2 {
		return Version{}, fmt.Errorf("unexpected regexp match: %v", match)
	}

	return Version{
		Version: match[1],
		Found:   true,
	}, nil
}

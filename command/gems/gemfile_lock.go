package gems

import (
	"fmt"
	"regexp"
	"strings"
)

// GemVersion contains bundler or fastlane version
type GemVersion struct {
	Version string
	Found   bool
}

// ParseFastlaneVersion  returns the fastlane version parsed from a Gemfile.lock on a best effort basis, for logging purposes only.
//
// for the following Gemfile.lock example, it returns: ">= 2.0)"
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
func ParseFastlaneVersion(gemfileLockContent string) (gemVersion GemVersion, err error) {
	relevantLines := []string{}
	lines := strings.Split(gemfileLockContent, "\n")

	specsStart := false
	for _, line := range lines {
		if strings.Contains(line, "specs:") {
			specsStart = true
		}

		trimmed := strings.Trim(line, " ")
		if trimmed == "" {
			specsStart = false
		}

		if specsStart {
			relevantLines = append(relevantLines, line)
		}
	}

	//     fastlane (1.109.0)
	exp := regexp.MustCompile(fmt.Sprintf(`^%s \((.+)\)`, regexp.QuoteMeta("fastlane")))
	for _, line := range relevantLines {
		match := exp.FindStringSubmatch(strings.TrimSpace(line))
		if match == nil {
			continue
		}
		if len(match) != 2 {
			return GemVersion{}, fmt.Errorf("unexpected regexp match: %v", match)
		}
		return GemVersion{
			Version: match[1],
			Found:   true,
		}, nil
	}

	return GemVersion{}, nil
}

// ParseBundlerVersion returns the bundler version used to create the bundle
func ParseBundlerVersion(gemfileLockContent string) (gemVersion GemVersion, err error) {
	/*
		BUNDLED WITH
			1.17.1
	*/
	bundlerRegexp := regexp.MustCompile(`(?m)^BUNDLED WITH\n\s+(\S+)`)
	match := bundlerRegexp.FindStringSubmatch(gemfileLockContent)
	if match == nil {
		return GemVersion{}, nil
	}
	if len(match) != 2 {
		return GemVersion{}, fmt.Errorf("unexpected regexp match: %v", match)
	}

	return GemVersion{
		Version: match[1],
		Found:   true,
	}, nil
}

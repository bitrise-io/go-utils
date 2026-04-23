package urlutil

import (
	"errors"
	"net/url"
	"strings"
)

// Join concatenates URL path elements with slash separators. The first
// element must be an absolute URL carrying both scheme and host; subsequent
// elements are appended as path segments. Leading and trailing slashes on
// intermediate elements are normalized away. A trailing slash on the last
// element is preserved.
func Join(elems ...string) (string, error) {
	if len(elems) < 1 {
		return "", errors.New("No elements defined to Join")
	}

	u, err := url.Parse(elems[0])
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		return "", errors.New("No Scheme defined")
	}
	if u.Host == "" {
		return "", errors.New("No Host defined")
	}

	var b strings.Builder
	b.WriteString(u.Scheme)
	b.WriteString("://")
	b.WriteString(u.Host)
	if firstPath := strings.Trim(u.Path, "/"); firstPath != "" {
		b.WriteByte('/')
		b.WriteString(firstPath)
	}

	lastIdx := len(elems) - 1
	for i := 1; i < len(elems); i++ {
		seg := strings.TrimLeft(elems[i], "/")
		if i != lastIdx {
			seg = strings.TrimRight(seg, "/")
		}
		b.WriteByte('/')
		b.WriteString(seg)
	}

	return b.String(), nil
}

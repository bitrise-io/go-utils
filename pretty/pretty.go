package pretty

import (
	"encoding/json"
	"fmt"
)

// Object returns the indented JSON representation of o. If o cannot be
// marshaled as JSON, Object falls back to Go's default fmt format.
func Object(o any) string {
	b, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		return fmt.Sprint(o)
	}
	return string(b)
}

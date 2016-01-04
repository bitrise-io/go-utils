package sliceutil

// IndexOfStringInSlice ...
func IndexOfStringInSlice(searchFor string, searchIn []string) int {
	for idx, anItm := range searchIn {
		if anItm == searchFor {
			return idx
		}
	}
	return -1
}

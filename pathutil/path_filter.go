package pathutil

// FilterFunc ...
type FilterFunc func(string) (bool, error)

// FilterPaths ...
func FilterPaths(fileList []string, filters ...FilterFunc) ([]string, error) {
	filtered := []string{}

	for _, pth := range fileList {
		allowed := true
		for _, filter := range filters {
			if allows, err := filter(pth); err != nil {
				return []string{}, err
			} else if !allows {
				allowed = false
				break
			}
		}
		if allowed {
			filtered = append(filtered, pth)
		}
	}

	return filtered, nil
}

package pkg

func mergeMaps(src, target map[string]string) map[string]string {
	if target == nil {
		target = map[string]string{}
	}
	for key, value := range src {
		target[key] = value
	}
	return target
}

func mergeArrays(src, target []string) []string {
	if target == nil {
		target = []string{}
	}
	for _, value := range src {
		if !contains(target, value) {
			target = append(target, value)
		}
	}
	return target
}

func contains(arr []string, entry string) bool {
	for _, value := range arr {
		if value == entry {
			return true
		}
	}
	return false
}

func appendStrings(arr []string, elems ...string) []string {
	if arr == nil {
		arr = []string{}
	}
	return append(arr, elems...)
}

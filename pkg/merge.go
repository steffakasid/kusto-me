package pkg

func mergeMaps(src, target map[string]string) map[string]string {
	for key, value := range src {
		target[key] = value
	}
	return target
}

func mergeArrays(src, target []string) []string {
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

package pkg

func removeFromArray(arr []string, elem string) []string {
	a := []string{}

	for _, c := range arr {
		if c != elem {
			a = append(a, c)
		}
	}
	return a
}

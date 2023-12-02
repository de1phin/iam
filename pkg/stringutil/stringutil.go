package stringutil

func SliceContains(slice []string, entry string) bool {
	for _, x := range slice {
		if x == entry {
			return true
		}
	}
	return false
}

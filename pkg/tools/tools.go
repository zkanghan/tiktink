package tools

func SliceIntToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{}, len(slice))
	for _, v := range slice {
		set[v] = struct{}{}
	}
	return set
}

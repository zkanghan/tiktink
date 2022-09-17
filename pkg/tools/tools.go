package tools

func SliceIntToSet(slice []int64) map[int64]struct{} {
	set := make(map[int64]struct{}, len(slice))
	for _, v := range slice {
		set[v] = struct{}{}
	}
	return set
}

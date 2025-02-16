package fmi2

func Transform[To, From any](source []From, f func(int, From) To) []To {
	vsm := make([]To, 0, len(source))
	for i, v := range source {
		vsm = append(vsm, f(i, v))
	}
	return vsm
}

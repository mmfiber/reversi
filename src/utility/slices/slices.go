package slices

func Filter[T any](slice []T, f func(T) bool) []T {
	var res []T
	for _, v := range slice {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func Map[Input any, Output any](slice []Input, f func(Input) Output) []Output {
	res := make([]Output, len(slice))
	for i, v := range slice {
		res[i] = f(v)
	}
	return res
}

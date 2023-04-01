package utils

func Remove[T comparable](l []T, item T) []T {
	for i, el := range l {
		if el == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func InsertBetween[T comparable](l []T, item T, prevItem T, nextItem T) []T {
	for i, el := range l {
		if el == prevItem {
			ret := append(l[:i+1], l[i:]...)
			ret[i] = item
			return ret
		} else if el == nextItem {
			if i-1 < 0 {
				return append([]T{item}, l...)
			}
			ret := append(l[:i], l[i-1:]...)
			ret[i] = item
			return ret
		}
	}
	return l
}

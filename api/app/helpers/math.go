package helpers

import "golang.org/x/exp/constraints"

func Max[T constraints.Integer | constraints.Float](numbers ...T) T {
	max := numbers[0]

	for _, v := range numbers[1:] {
		if max < v {
			max = v
		}
	}

	return max
}

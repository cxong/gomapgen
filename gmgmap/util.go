package gmgmap

import "math/rand"

type rect struct {
	x, y, w, h int
}

func randomWalk(x, y, w, h int) (int, int) {
	for {
		// Choose random direction, up/right/down/left
		switch rand.Intn(4) {
		case 0:
			// up
			if y > 0 {
				return x, y - 1
			}
		case 1:
			// right
			if x < w-1 {
				return x + 1, y
			}
		case 2:
			// down
			if y < h-1 {
				return x, y + 1
			}
		case 3:
			// left
			if x > 0 {
				return x - 1, y
			}
		}
	}
}

func imin(i1, i2 int) int {
	if i1 < i2 {
		return i1
	}
	return i2
}

func imax(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

func iclamp(v, min, max int) int {
	switch {
	case v < min:
		return min
	case v > max:
		return max
	default:
		return v
	}
}

func irand(min, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}

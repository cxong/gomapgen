package gmgmap

import "math/rand"

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

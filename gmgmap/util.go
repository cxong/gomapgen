package gmgmap

import (
	"errors"
	"math/rand"
)

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

type bspRoom struct {
	r              rect
	parent         int
	child1, child2 int
	level          int
}

func bspRoomRoot(width, height int) bspRoom {
	return bspRoom{rect{0, 0, width, height}, -1, -1, -1, 0}
}

func split(room *bspRoom, i, minRoomSize int) (bspRoom, bspRoom, error) {
	// If more than 3:2, split the long dimension, otherwise randomise
	if room.r.w*3 > room.r.h*2 || (room.r.h*3 < room.r.w*2 && rand.Intn(2) == 0) {
		// Split horizontally
		r := room.r.w - minRoomSize*2
		var x int
		if r < 0 {
			return bspRoom{}, bspRoom{}, errors.New("room too small")
		} else if r == 0 {
			x = minRoomSize
		} else {
			x = rand.Intn(r) + minRoomSize
		}
		return bspRoom{rect{room.r.x, room.r.y, x, room.r.h}, i, -1, -1, room.level + 1},
			bspRoom{rect{room.r.x + x, room.r.y, room.r.w - x, room.r.h}, i, -1, -1, room.level + 1},
			nil
	}
	// Split vertically
	r := room.r.h - minRoomSize*2
	var y int
	if r < 0 {
		return bspRoom{}, bspRoom{}, errors.New("room too small")
	} else if r == 0 {
		y = minRoomSize
	} else {
		y = rand.Intn(r) + minRoomSize
	}
	return bspRoom{rect{room.r.x, room.r.y, room.r.w, y}, i, -1, -1, room.level + 1},
		bspRoom{rect{room.r.x, room.r.y + y, room.r.w, room.r.h - y}, i, -1, -1, room.level + 1},
		nil
}

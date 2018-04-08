package gmgmap

import (
	"errors"
	"math/rand"
)

type rect struct {
	x, y, w, h int
}

func rectIsAdjacent(r1, r2 rect, overlapSize int) bool {
	// If left/right edges adjacent
	if r1.x-(r2.x+r2.w) == 0 || r2.x-(r1.x+r1.w) == 0 {
		return r1.y+overlapSize < r2.y+r2.h && r2.y+overlapSize < r1.y+r1.h
	}
	if r1.y-(r2.y+r2.h) == 0 || r2.y-(r1.y+r1.h) == 0 {
		return r1.x+overlapSize < r2.x+r2.w && r2.x+overlapSize < r1.x+r1.w
	}
	return false
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

func bspSplit(room *bspRoom, i, minRoomSize, maxRoomSize int) (bspRoom, bspRoom, error) {
	// If the room is too small, then don't split
	if room.r.w-minRoomSize*2 < 0 && room.r.h-minRoomSize < 0 {
		return bspRoom{}, bspRoom{}, errors.New("room too small")
	}
	// If the room is small enough already, consider not splitting
	if room.r.w <= maxRoomSize && room.r.h <= maxRoomSize && rand.Intn(2) == 0 {
		return bspRoom{}, bspRoom{}, errors.New("room is small enough")
	}
	// If more than 2:1, split the long dimension, otherwise randomise
	if room.r.w*2 > room.r.h || (room.r.h*2 < room.r.w && rand.Intn(2) == 0) {
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

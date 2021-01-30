package gmgmap

import (
	"errors"
	"math"
	"math/rand"
)

type vec2 struct {
	x, y int
}

type rect struct {
	x, y, w, h int
}

func (r rect) IsAdjacent(r2 rect, overlapSize int) bool {
	// If left/right edges adjacent
	if r.x-(r2.x+r2.w) == 0 || r2.x-(r.x+r.w) == 0 {
		return r.y+overlapSize < r2.y+r2.h && r2.y+overlapSize < r.y+r.h
	}
	if r.y-(r2.y+r2.h) == 0 || r2.y-(r.y+r.h) == 0 {
		return r.x+overlapSize < r2.x+r2.w && r2.x+overlapSize < r.x+r.w
	}
	return false
}

func (r rect) Overlaps(r2 rect) bool {
	return r.x < r2.x+r2.w && r.x+r.w > r2.x &&
		r.y < r2.y+r2.h && r.y+r.h > r2.y
}

func (r rect) isIn(x, y int) bool {
	return x >= r.x && x < r.x+r.w && y >= r.y && y < r.y+r.h
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
	horizontal     bool
}

func bspRoomRoot(width, height int) bspRoom {
	return bspRoom{rect{0, 0, width, height}, -1, -1, -1, 0, false}
}

func (room *bspRoom) Split(i, minRoomSize, maxRoomSize int) (bspRoom, bspRoom, error) {
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
		return room.SplitHorizontal(i, minRoomSize)
	}
	return room.SplitVertical(i, minRoomSize)
}

// Split rooms horizontally (left + right children)
func (room *bspRoom) SplitHorizontal(i, minRoomSize int) (bspRoom, bspRoom, error) {
	r := room.r.w - minRoomSize*2
	var x int
	if r < 0 {
		return bspRoom{}, bspRoom{}, errors.New("room too small")
	} else if r == 0 {
		x = minRoomSize
	} else {
		x = rand.Intn(r) + minRoomSize
	}
	return bspRoom{rect{room.r.x, room.r.y, x, room.r.h}, i, -1, -1, room.level + 1, true},
		bspRoom{rect{room.r.x + x, room.r.y, room.r.w - x, room.r.h}, i, -1, -1, room.level + 1, true},
		nil
}

// Split rooms horizontally (top + bottom children)
func (room *bspRoom) SplitVertical(i, minRoomSize int) (bspRoom, bspRoom, error) {
	r := room.r.h - minRoomSize*2
	var y int
	if r < 0 {
		return bspRoom{}, bspRoom{}, errors.New("room too small")
	} else if r == 0 {
		y = minRoomSize
	} else {
		y = rand.Intn(r) + minRoomSize
	}
	return bspRoom{rect{room.r.x, room.r.y, room.r.w, y}, i, -1, -1, room.level + 1, false},
		bspRoom{rect{room.r.x, room.r.y + y, room.r.w, room.r.h - y}, i, -1, -1, room.level + 1, false},
		nil
}

func (room *bspRoom) IsLeaf() bool {
	return room.child1 < 0 && room.child2 < 0
}

// Abs - absolute value, integer
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattanDistance(x1, y1, x2, y2 int) int {
	return Abs(x1-x2) + Abs(y1-y2)
}

func euclideanDistance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2))
}

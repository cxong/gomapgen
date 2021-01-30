package gmgmap

import (
	"math/rand"
)

type bspArea struct {
	bspRoom
	isStreet    bool
	isConnected bool
}

func (s bspArea) dAlong() vec2 {
	if s.horizontal {
		return vec2{1, 0}
	}
	return vec2{0, 1}
}

func (s bspArea) dAcross() vec2 {
	if s.horizontal {
		return vec2{0, 1}
	}
	return vec2{1, 0}
}

// NewBSPInterior - Create new BSP interior map
// Implementation of https://gamedev.stackexchange.com/questions/47917/procedural-house-with-rooms-generator/48216#48216
func NewBSPInterior(width, height, splits, minRoomSize, corridorWidth int) *Map {
	corridorLevelDiffBlock := 1
	m := NewMap(width, height)
	var areas []bspArea

	// Split the map for a number of iterations, choosing alternating axis and random location
	hcount := rand.Intn(2)
	areas = append(areas, bspArea{bspRoomRoot(width, height), false, false})
	for i := 0; i < len(areas); i++ {
		if areas[i].level == splits {
			break
		}
		var r1, r2 bspRoom
		var err error = nil
		// Alternate splitting direction per level
		horizontal := ((hcount + areas[i].level) % 2) == 1
		if horizontal {
			r1, r2, err = areas[i].SplitHorizontal(i, minRoomSize+corridorWidth/2)
		} else {
			r1, r2, err = areas[i].SplitVertical(i, minRoomSize+corridorWidth/2)
		}
		if err == nil {
			// Resize rooms to allow space for street
			for j := 0; j < corridorWidth; j++ {
				if horizontal {
					if j%2 == 0 {
						r1.r.w--
					} else {
						r2.r.x++
						r2.r.w--
					}
				} else {
					if j%2 == 0 {
						r1.r.h--
					} else {
						r2.r.y++
						r2.r.h--
					}
				}
			}
			// Replace current area with a street
			areas[i].isStreet = true
			if horizontal {
				areas[i].r = rect{r1.r.x + r1.r.w, r1.r.y, corridorWidth, r1.r.h}
			} else {
				areas[i].r = rect{r1.r.x, r1.r.y + r1.r.h, r1.r.w, corridorWidth}
			}
			areas[i].horizontal = !horizontal
			areas[i].child1 = len(areas)
			areas = append(areas, bspArea{r1, false, false})
			areas[i].child2 = len(areas)
			areas = append(areas, bspArea{r2, false, false})
		}
	}
	// Try to split leaf rooms into more rooms, by longest axis
	for i := 0; i < len(areas); i++ {
		if areas[i].isStreet {
			continue
		}

		var r1, r2 bspRoom
		var err error = nil
		if areas[i].r.w > areas[i].r.h {
			r1, r2, err = areas[i].SplitHorizontal(i, minRoomSize)
		} else {
			r1, r2, err = areas[i].SplitVertical(i, minRoomSize)
		}
		if err == nil {
			// Resize rooms so they share a splitting wall
			if r1.horizontal {
				r1.r.w++
			} else {
				r1.r.h++
			}
			areas[i].child1 = len(areas)
			areas = append(areas, bspArea{r1, false, false})
			areas[i].child2 = len(areas)
			areas = append(areas, bspArea{r2, false, false})
		}
	}

	g := m.Layer("Ground")
	s := m.Layer("Structures")

	// Find deepest leaf going down both branches; place stairs
	// This represents longest path
	deepestRoom1 := findDeepestRoomFrom(areas, areas[0].child1)
	placeInsideRoom(s, deepestRoom1.r, stairsUp)
	deepestRoom2 := findDeepestRoomFrom(areas, areas[0].child2)
	placeInsideRoom(s, deepestRoom2.r, stairsDown)

	// Fill rooms
	for i := range areas {
		// Skip non-leaves
		if !areas[i].IsLeaf() {
			continue
		}
		r := areas[i].r
		g.rectangleFilled(rect{r.x + 1, r.y + 1, r.w - 2, r.h - 2}, room)
		s.rectangleUnfilled(r, wall2)
	}

	// Add doors leading to the closest street in the hierarchy
	for i := range areas {
		// Skip non-leaves
		if !areas[i].IsLeaf() {
			continue
		}
		r := areas[i].r
		// Add doors leading to hallways
		street := areas[i]
		for {
			if street.isStreet {
				break
			}
			street = areas[street.parent]
		}
		for j := 0; j < 4; j++ {
			if areas[i].isConnected {
				break
			}
			doorPos := vec2{r.x + r.w/2, r.y + r.h/2}
			var outsideDoor vec2
			if j == 0 {
				// top
				doorPos.y = r.y
				outsideDoor = vec2{doorPos.x, doorPos.y - 1}
			} else if j == 1 {
				// right
				doorPos.x = r.x + r.w - 1
				outsideDoor = vec2{doorPos.x + 1, doorPos.y}
			} else if j == 2 {
				// bottom
				doorPos.y = r.y + r.h - 1
				outsideDoor = vec2{doorPos.x, doorPos.y + 1}
			} else {
				// left
				doorPos.x = r.x
				outsideDoor = vec2{doorPos.x - 1, doorPos.y}
			}
			if street.r.isIn(outsideDoor.x, outsideDoor.y) {
				g.setTile(doorPos.x, doorPos.y, room)
				s.setTile(doorPos.x, doorPos.y, door)
				areas[i].isConnected = true
				break
			}
		}
	}
	// For every room, connect it to a random shallower room
	// Keep going until all rooms are connected
	for {
		numUnconnected := 0
		for i := range areas {
			if areas[i].isConnected || areas[i].isStreet || !areas[i].IsLeaf() {
				continue
			}
			numUnconnected++
			r := areas[i].r
			// Shrink rectangles by 1 to determine overlap
			r.w--
			r.h--
			overlapSize := 1
			for j := range areas {
				if !areas[j].IsLeaf() {
					continue
				}
				if i == j {
					continue
				}
				roomOther := areas[j]
				// Only connect to a room that is also connected
				if !roomOther.isConnected {
					continue
				}
				rOther := roomOther.r
				// Shrink rectangles by 1 to determine overlap
				rOther.w--
				rOther.h--
				if !r.IsAdjacent(rOther, overlapSize) {
					continue
				}
				// Rooms are adjacent; pick the cell that's in the middle of the
				// adjacent area and turn into a door
				minOverlapX := imin(
					areas[i].r.x+areas[i].r.w, roomOther.r.x+roomOther.r.w)
				maxOverlapX := imax(areas[i].r.x, roomOther.r.x)
				minOverlapY := imin(
					areas[i].r.y+areas[i].r.h, roomOther.r.y+roomOther.r.h)
				maxOverlapY := imax(areas[i].r.y, roomOther.r.y)
				overlapX := (minOverlapX + maxOverlapX) / 2
				overlapY := (minOverlapY + maxOverlapY) / 2
				g.setTile(overlapX, overlapY, room2)
				s.setTile(overlapX, overlapY, door)
				areas[i].isConnected = true
				numUnconnected--
				break
			}
		}
		if numUnconnected == 0 {
			break
		}
	}

	// Fill streets
	for i := range areas {
		if !areas[i].isStreet {
			continue
		}
		g.rectangleFilled(areas[i].r, room2)
		// Check ends of street - if next to much older street, block off with wall
		end1 := vec2{areas[i].r.x, areas[i].r.y}
		end2 := vec2{areas[i].r.x + areas[i].r.w - 1, areas[i].r.y + areas[i].r.h - 1}
		capStreet(g, s, areas, areas[i], end1, areas[i].dAcross(), areas[i].dAlong(), corridorWidth, corridorLevelDiffBlock)
		capStreet(g, s, areas, areas[i], end2, vec2{-areas[i].dAcross().x, -areas[i].dAcross().y}, vec2{-areas[i].dAlong().x, -areas[i].dAlong().y}, corridorWidth, corridorLevelDiffBlock)
	}

	return m
}

func capStreet(g, s *Layer, streets []bspArea, st bspArea, end, dAcross, dAlong vec2, corridorWidth, corridorLevelDiffBlock int) {
	// Check ends of street - if outside map, or next to much older street, block off with wall
	outside := vec2{end.x - dAlong.x, end.y - dAlong.y}
	doCap := false
	if !g.isIn(outside.x, outside.y) {
		doCap = true
	} else {
		for i := range streets {
			if streets[i].r.isIn(outside.x, outside.y) && st.level-streets[i].level > corridorLevelDiffBlock {
				doCap = true
				break
			}
		}
	}
	if doCap {
		for i := 0; i < corridorWidth; i++ {
			g.setTile(end.x+dAcross.x*i, end.y+dAcross.y*i, nothing)
			s.setTile(end.x+dAcross.x*i, end.y+dAcross.y*i, wall2)
		}
	}
}

func findDeepestRoomFrom(areas []bspArea, child int) *bspArea {
	var pathStack []bspArea
	pathStack = append(pathStack, areas[child])
	var deepestChild *bspArea = nil
	maxDepth := 0
	for len(pathStack) > 0 {
		r := pathStack[len(pathStack)-1]
		pathStack = pathStack[:len(pathStack)-1]
		if r.IsLeaf() {
			if r.level > maxDepth {
				maxDepth = r.level
				deepestChild = &r
			}
		}
		if r.child1 >= 0 {
			pathStack = append(pathStack, areas[r.child1])
		}
		if r.child2 >= 0 {
			pathStack = append(pathStack, areas[r.child2])
		}
	}
	return deepestChild
}

func placeInsideRoom(s *Layer, r rect, t rune) {
	s.setTile(r.x+r.w/2, r.y+r.h/2, t)
}

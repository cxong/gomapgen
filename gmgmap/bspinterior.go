package gmgmap

import (
	"math/rand"
)

type street struct {
	r          rect
	horizontal bool
	level      int
}

func (s street) dAlong() vec2 {
	if s.horizontal {
		return vec2{1, 0}
	}
	return vec2{0, 1}
}

func (s street) dAcross() vec2 {
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

	// Split the map for a number of iterations, choosing alternating axis and random location
	var areas []bspRoom
	var streets []street
	hcount := rand.Intn(2)
	areas = append(areas, bspRoomRoot(width, height))
	for i := 0; i < len(areas); i++ {
		if areas[i].level == splits {
			break
		}
		var r1, r2 bspRoom
		var err error = nil
		// Alternate splitting direction per level
		horizontal := ((hcount + areas[i].level) % 2) == 1
		if horizontal {
			r1, r2, err = bspSplitHorizontal(&areas[i], i, minRoomSize+corridorWidth/2)
		} else {
			r1, r2, err = bspSplitVertical(&areas[i], i, minRoomSize+corridorWidth/2)
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
			areas[i].child1 = len(areas)
			areas = append(areas, r1)
			areas[i].child2 = len(areas)
			areas = append(areas, r2)
			var s street
			if horizontal {
				s.r = rect{r1.r.x + r1.r.w, r1.r.y, corridorWidth, r1.r.h}
			} else {
				s.r = rect{r1.r.x, r1.r.y + r1.r.h, r1.r.w, corridorWidth}
			}
			s.level = r1.level
			s.horizontal = !horizontal
			streets = append(streets, s)
		}
	}

	g := m.Layer("Ground")
	s := m.Layer("Structures")
	// Turn the leaves into rooms
	for i := 0; i < len(areas); i++ {
		// Only place rooms in leaf nodes
		if areas[i].child1 >= 0 || areas[i].child2 >= 0 {
			continue
		}

		// Try to split into more rooms, length-wise
		var r1, r2 bspRoom
		var err error = nil
		if !areas[i].horizontal {
			r1, r2, err = bspSplitHorizontal(&areas[i], i, minRoomSize)
		} else {
			r1, r2, err = bspSplitVertical(&areas[i], i, minRoomSize)
		}
		if err == nil {
			// Resize rooms so they share a splitting wall
			if r1.horizontal {
				r1.r.w++
			} else {
				r1.r.h++
			}
			areas[i].child1 = len(areas)
			areas = append(areas, r1)
			areas[i].child2 = len(areas)
			areas = append(areas, r2)
			continue
		}

		var r rect
		r.w = areas[i].r.w
		r.x = areas[i].r.x
		r.h = areas[i].r.h
		r.y = areas[i].r.y
		g.rectangleFilled(rect{r.x + 1, r.y + 1, r.w - 2, r.h - 2}, room)
		s.rectangleUnfilled(r, wall2)
		// Add doors leading to hallways
		for j := 0; j < 4; j++ {
			doorPos := vec2{areas[i].r.x + areas[i].r.w/2, areas[i].r.y + areas[i].r.h/2}
			var outsideDoor vec2
			if j == 0 {
				// top
				doorPos.y = areas[i].r.y
				outsideDoor = vec2{doorPos.x, doorPos.y - 1}
			} else if j == 1 {
				// right
				doorPos.x = areas[i].r.x + areas[i].r.w - 1
				outsideDoor = vec2{doorPos.x + 1, doorPos.y}
			} else if j == 2 {
				// bottom
				doorPos.y = areas[i].r.y + areas[i].r.h - 1
				outsideDoor = vec2{doorPos.x, doorPos.y + 1}
			} else {
				// left
				doorPos.x = areas[i].r.x
				outsideDoor = vec2{doorPos.x - 1, doorPos.y}
			}
			for i := range streets {
				if streets[i].r.isIn(outsideDoor.x, outsideDoor.y) {
					g.setTile(doorPos.x, doorPos.y, room)
					s.setTile(doorPos.x, doorPos.y, door)
					break
				}
			}
		}
	}

	// Fill streets
	for i := range streets {
		g.rectangleFilled(streets[i].r, room2)
		// Check ends of street - if next to much older street, block off with wall
		end1 := vec2{streets[i].r.x, streets[i].r.y}
		end2 := vec2{streets[i].r.x + streets[i].r.w - 1, streets[i].r.y + streets[i].r.h - 1}
		capStreet(g, s, streets, streets[i], end1, streets[i].dAcross(), streets[i].dAlong(), corridorWidth, corridorLevelDiffBlock)
		capStreet(g, s, streets, streets[i], end2, vec2{-streets[i].dAcross().x, -streets[i].dAcross().y}, vec2{-streets[i].dAlong().x, -streets[i].dAlong().y}, corridorWidth, corridorLevelDiffBlock)
	}

	// Place stairs going up at end of first (main) street
	s.setTile(streets[0].r.x+streets[0].dAlong().x, streets[0].r.y+streets[0].dAlong().y, stairsUp)
	// Place stairs going down in last room
	lastRoomRect := areas[len(areas)-1].r
	s.setTile(lastRoomRect.x+lastRoomRect.w/2, lastRoomRect.y+lastRoomRect.h/2, stairsDown)

	return m
}

func capStreet(g, s *Layer, streets []street, st street, end, dAcross, dAlong vec2, corridorWidth, corridorLevelDiffBlock int) {
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

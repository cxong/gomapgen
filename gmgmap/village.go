package gmgmap

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/beefsack/go-astar"
)

type building struct {
	r          rect
	importance int
}

func (b building) addNPC(c *Layer) {
	// Try to place a random NPC somewhere inside the building
	c.setTileInAreaIfEmpty(rect{b.r.x + 1, b.r.y + 1, b.r.w - 2, b.r.h - 2}, player)
}

// NewVillage - create a village, made up of multiple buildings
func NewVillage(width, height, buildingPadding int) *Map {
	m := NewMap(width, height)
	g := m.Layer("Ground")
	s := m.Layer("Structures")
	f := m.Layer("Furniture")

	// Grass
	g.fill(grass)

	buildings := genBuildings(width, height, buildingPadding)
	assignBuildingImportance(buildings)
	placeBuildings(g, s, f, buildings)
	addPaths(g, s, buildings)
	c := m.Layer("Characters")
	placeNPCs(c, buildings)

	return m
}

func genBuildings(width, height, buildingPadding int) []building {
	buildings := make([]building, 0)
	// Keep placing buildings for a while
	for i := 0; i < 500; i++ {
		w := rand.Intn(3) + 5
		h := rand.Intn(3) + 5
		x := rand.Intn(width - w)
		y := rand.Intn(height - h)
		if x < 0 || y < 0 {
			continue
		}
		// Check if it overlaps with any existing buildings
		overlaps := false
		for _, b := range buildings {
			// Add a bit of padding between the buildings
			if b.r.Overlaps(
				rect{
					x - buildingPadding,
					y - buildingPadding,
					w + buildingPadding*2,
					h + buildingPadding*2}) {
				overlaps = true
				break
			}
		}
		if overlaps {
			continue
		}
		buildings = append(buildings, building{rect{x, y, w, h}, 0})
	}
	return buildings
}

func assignBuildingImportance(buildings []building) {
	for i := range buildings {
		imp := int(math.Pow(float64(rand.Float32()*3.0+1), 2))
		buildings[i].importance = imp
	}
}

func placeBuildings(g, s, f *Layer, buildings []building) {
	for _, building := range buildings {
		imp := building.importance
		// Use tiles based on importance
		tileRoom, tileWall := room, wall
		if imp > 10 {
			tileRoom, tileWall = room2, wall2
		}
		hasSign := imp > 5
		addBuilding(g, s, f, building.r, tileRoom, tileWall, hasSign)
	}
}

func addPaths(g, s *Layer, buildings []building) {
	// Draw paths between random pairs of entrances via importance
	// Ensure at least one path exists for all buildings

	impSum := 0
	for _, building := range buildings {
		impSum += building.importance
	}

	world := newVillageWorld(g.Width, g.Height, s)
	buildingsWithPaths := map[int]bool{}
	numPaths := len(buildings) * 3
	for i := 0; i < numPaths || len(buildingsWithPaths) < len(buildings); i++ {
		for {
			// Check for path valid and exists
			building1 := rand.Intn(len(buildings))
			// randomly select second building by importance
			impFact := rand.Intn(impSum)
			building2 := 0
			impFactSum := 0
			for j, b2 := range buildings {
				impFactSum += b2.importance
				if impFactSum > impFact {
					building2 = j
					break
				}
			}
			if building1 == building2 {
				continue
			}
			buildingsWithPaths[building1] = true
			buildingsWithPaths[building2] = true
			// TODO: find entrance and start/end paths there
			b1 := buildings[building1]
			b2 := buildings[building2]
			startX := b1.r.x + b1.r.w/2
			startY := b1.r.y + b1.r.h
			endX := b2.r.x + b2.r.w/2
			endY := b2.r.y + b2.r.h
			path, _, found := world.addPath(startX, startY, endX, endY)
			if !found {
				fmt.Println("Could not find path")
			} else {
				for _, t := range path {
					world.incUsage(t.(*villageTile).x, t.(*villageTile).y)
				}
			}
			break
		}
	}

	placePaths(g, s, world, tree, grass, road, road2)
}

func placePaths(g, s *Layer, world villageWorld, usage0, usage1, usage2, usage3 rune) {
	// Draw paths based on how well they're used
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			usage := world.getUsage(x, y)
			if usage == 0 {
				if g.getTile(x, y) == grass && s.getTile(x, y) == nothing {
					s.setTile(x, y, usage0)
				}
			} else if usage <= 3 {
				g.setTile(x, y, usage1)
			} else if usage <= 6 {
				g.setTile(x, y, usage2)
			} else {
				g.setTile(x, y, usage3)
			}
		}
	}
}

func placeNPCs(c *Layer, buildings []building) {
	// Place NPCs based on importance
	for _, building := range buildings {
		for i := 0; i < building.importance/2; i++ {
			building.addNPC(c)
		}
	}
}

func addBuilding(g, s, f *Layer, r rect, tileRoom, tileWall rune, hasSign bool) {
	// Perimeter
	s.rectangle(r, tileWall, false)
	// Floor
	g.rectangle(rect{r.x + 1, r.y + 1, r.w - 2, r.h - 2}, tileRoom, true)
	// Entrance
	entranceX := r.x + r.w/2
	entranceY := r.y + r.h - 1
	g.setTile(entranceX, entranceY, tileRoom)
	s.setTile(entranceX, entranceY, door)
	if hasSign {
		f.setTile(entranceX-1, entranceY, sign)
	}
}

// Special tile types for A* to find paths via erosion
type villageTile struct {
	x, y  int
	s     *Layer
	w     villageWorld
	usage int
}

func (t *villageTile) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.s.getTile(t.x+offset[0], t.y+offset[1]); n == nothing {
			neighbors = append(neighbors, t.w.tile(t.x+offset[0], t.y+offset[1]))
		}
	}
	return neighbors
}

func (t *villageTile) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*villageTile)
	// Max cost 1.5, min cost 1 (asymptote)
	return 0.5/float64(toT.usage+1) + 1
}

func (t *villageTile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*villageTile)
	return euclideanDistance(t.x, t.y, toT.x, toT.y)
}

type villageWorld map[int]map[int]*villageTile

func newVillageWorld(w, h int, s *Layer) villageWorld {
	world := villageWorld{}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			world.setTile(&villageTile{x, y, s, world, 0}, x, y)
		}
	}
	return world
}

func (w villageWorld) tile(x, y int) *villageTile {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

func (w villageWorld) setTile(t *villageTile, x, y int) {
	if w[x] == nil {
		w[x] = map[int]*villageTile{}
	}
	w[x][y] = t
	t.x = x
	t.y = y
	t.w = w
	t.usage = 0
}

func (w villageWorld) incUsage(x, y int) {
	w[x][y].usage++
}

func (w villageWorld) getUsage(x, y int) int {
	return w[x][y].usage
}

// Use A* to find and return a path between two points
// A* will avoid any tiles where there's something in the structure (s) layer
func (w villageWorld) addPath(x1, y1, x2, y2 int) (path []astar.Pather, distance float64, found bool) {
	return astar.Path(w.tile(x1, y1), w.tile(x2, y2))
}

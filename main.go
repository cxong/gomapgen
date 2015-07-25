package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/cxong/gomapgen/gmgmap"
)

func main() {
	algo := flag.String("algo", "rogue", "generation algorithm: walk/rogue")
	width := flag.Int("width", 32, "map width")
	height := flag.Int("height", 32, "map height")
	iterations := flag.Int("iterations", 3000, "number of iterations for random walk algo")
	gridWidth := flag.Int("gridwidth", 3, "grid size, for rogue algo")
	gridHeight := flag.Int("gridheight", 3, "grid size, for rogue algo")
	minRoomPct := flag.Int("minroompct", 50, "percent of rooms per grid, for rogue algo")
	maxRoomPct := flag.Int("maxroompct", 100, "percent of rooms per grid, for rogue algo")
	seed := flag.Int64("seed", time.Now().UTC().UnixNano(), "random seed")
	flag.Parse()
	// make map
	fmt.Println("Using seed", *seed)
	rand.Seed(*seed)
	m := gmgmap.NewMap(*width, *height)
	switch *algo {
	case "walk":
		m = gmgmap.NewRandomWalk(*width, *height, *iterations)
	case "rogue":
		m = gmgmap.NewRogue(*width, *height, *gridWidth, *gridHeight,
			*minRoomPct, *maxRoomPct)
	}
	// print
	m.Print()
	// export TMX
	template := gmgmap.DawnLikeTemplate
	m.ToTMX(&template)
}

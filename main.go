package main

import (
	"flag"

	"github.com/cxong/gomapgen/gmgmap"
)

func main() {
	algo := flag.String("algo", "rogue", "generation algorithm: walk/rogue")
	width := flag.Int("width", 32, "map width")
	height := flag.Int("height", 32, "map height")
	iterations := flag.Int("iterations", 3000, "number of iterations for random walk algo")
	minRoomSize := flag.Int("minroomsize", 3, "minimum room size, for rogue algo")
	maxRoomSize := flag.Int("maxroomsize", 15, "maximum room size, for rogue algo")
	flag.Parse()
	// make map
	m := gmgmap.NewMap(*width, *height)
	switch *algo {
	case "walk":
		m = gmgmap.NewRandomWalk(*width, *height, *iterations)
	case "rogue":
		m = gmgmap.NewRogue(*width, *height, *minRoomSize, *maxRoomSize)
	}
	// print
	m.Print()
	// export TMX
	template := gmgmap.DawnLikeTemplate
	m.ToTMX(&template)
}

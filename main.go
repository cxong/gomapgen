package main

import (
	"flag"

	"github.com/cxong/gomapgen/gmgmap"
)

func main() {
	width := flag.Int("width", 32, "map width")
	height := flag.Int("height", 32, "map height")
	//iterations := flag.Int("iterations", 3000, "number of iterations for random walk algo")
	flag.Parse()
	// make map
	//m := gmgmap.NewRandomWalk(*width, *height, *iterations)
	m := gmgmap.NewRogue(*width, *height, 3, 15)
	// print
	m.Print()
	// export TMX
	template := gmgmap.DawnLikeTemplate
	m.ToTMX(&template)
}

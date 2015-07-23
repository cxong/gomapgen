package main

import "github.com/cxong/gomapgen/gmgmap"

func main() {
	// make map
	m := gmgmap.NewRandomWalk(80, 24, 3000)
	// print
	m.Print()
	// export TMX
	template := gmgmap.DawnLikeTemplate
	m.ToTMX(&template)
}

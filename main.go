package main

import "github.com/cxong/gomapgen/gmgmap"

func main() {
	// make map
	m := gmgmap.NewMap(80, 24)
	// print
	m.Print()
}

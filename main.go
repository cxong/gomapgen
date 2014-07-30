package main

import "github.com/cxong/gomapgen/gmgmap"

func main() {
  // make map
  var m gmgmap.Map
  m.Width = 80
  m.Height = 24
  for y := 0; y < m.Height; y++ {
    for x := 0; x < m.Width; x++ {
      m.Tiles = append(m.Tiles, ' ')
    }
  }
  // print
  m.Print()
}
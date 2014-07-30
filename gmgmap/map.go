package gmgmap

import "fmt"

type Map struct {
  Tiles  []byte
  Width  int
  Height int
}

func (m Map) Print() {
  for y := 0; y < m.Height; y++ {
    // Upper frame
    if y == 0 {
      fmt.Print("+")
      for x := 0; x < m.Width; x++ {
        fmt.Print("-")
      }
      fmt.Print("+")
      fmt.Println()
    }
    
    // Left of frame
    fmt.Print("|")
    
    // Interior cells
    for x := 0; x < m.Width; x++ {
      var i = y * m.Width + x
      fmt.Printf("%c", m.Tiles[i])
    }
    
    // Right of frame
    fmt.Print("|")
    
    // Bottom frame
    if y == m.Height - 1 {
      fmt.Println()
      fmt.Print("+")
      for x := 0; x < m.Width; x++ {
        fmt.Print("-")
      }
      fmt.Print("+")
    }
    
    fmt.Println()
  }
}
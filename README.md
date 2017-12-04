GoMapGen
========

A 2d map generator written in Go

    $ go run main.go --algo=rogue
	Using seed 1512389956399933000
	+--------------------------------+
	|  WWWWWWW   WWWW                |
	|  W.....W   W..W                |
	|  W.....+## W..W    WWWWWWW     |
	|  WWW+WWW ##+..+##  W.....W     |
	|     #      W..W ###+.....W     |
	|     #      W..W    WWW+WWW     |
	|     #      WW+W       #        |
	|     #        #        #        |
	|     #        #        #        |
	|     #        #        #        |
	|     #      WW+W      ##        |
	|     #      W..W      #         |
	|  WWW+WWW   W..W      #         |
	|  W.....W   W..W      #         |
	|  W.....W ##+.<+###   #         |
	|  W.....+## W..W  #   #         |
	|  W.....W   W..W  #####         |
	|  WWW+WWW   W..W      #         |
	|    ##      WW+W      #         |
	|    #         #       #         |
	|WWWW+WWWW  ####       ###       |
	|W.......W  #            #       |
	|W.......W  #            #       |
	|W.......W  #            #       |
	|W.......+###            #       |
	|W.......W            WWW+WWW    |
	|W.......W            W.....W    |
	|WWWWWWWWW            W..>..W    |
	|                     WWWWWWW    |
	|                                |
	|                                |
	|                                |
	+--------------------------------+

This map generator implements a number of different algorithms and can output to ASCII, CSV and TMX tile map.

See `main.go` for all the options.

## --algo=rogue --width=30 --height=18

![rogue](https://raw.githubusercontent.com/cxong/gomapgen/master/examples/rogue.png)

## --algo=shop --width=16 --height=13

![shop](https://raw.githubusercontent.com/cxong/gomapgen/master/examples/shop.png)

## --algo=bsp --width=24 --height=20

![bsp](https://raw.githubusercontent.com/cxong/gomapgen/master/examples/bsp.png)

## --algo=walk --width=16 --height=16 --iterations=500

![walk](https://raw.githubusercontent.com/cxong/gomapgen/master/examples/walk.png)

## --algo=cell --width=24 --height=20 --template=kenney

![cell](https://raw.githubusercontent.com/cxong/gomapgen/master/examples/cell.png)




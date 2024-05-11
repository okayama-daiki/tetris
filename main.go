package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/okayama-daiki/tetris/tetris/game"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("EbiTetris")

	minoBag := game.MinoBag{}
	currentMino := minoBag.Next()
	holdingMino := game.HoldingMino{Available: true}
	board := game.NewBoard()

	var game = game.Game{
		MinoBag:              minoBag,
		CurrentMino:          currentMino,
		HoldingMino:          holdingMino,
		Board:                board,
		NormalDroppingSpeed:  60,
		CurrentDroppingSpeed: 60,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

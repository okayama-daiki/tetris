package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/okayama-daiki/tetris/tetris/game"
)

func main() {
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("EbiTetris")

	var minoBag = game.MinoBag{}
	var currentMino = game.Mino{IsGrounded: true}
	var holdingMino = game.HoldingMino{Available: true}
	var board = game.Board{}
	board.Init()

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
}

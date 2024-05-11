package game

import (
	"image/color"
)

var (
	Wall             = color.RGBA{108, 122, 137, 255}
	BACKGROUND_COLOR = color.RGBA{5, 5, 5, 255}
	LINE_COLOR       = color.RGBA{75, 75, 75, 255}
	BORDER_COLOR     = color.RGBA{240, 240, 240, 255}
)

type Board [OUTER_HEIGHT][OUTER_WIDTH]color.Color

func NewBoard() Board {
	board := Board{}
	for y := range OUTER_HEIGHT - 1 {
		board[y][0] = Wall
		board[y][OUTER_WIDTH-1] = Wall
	}
	for x := range OUTER_WIDTH {
		board[OUTER_HEIGHT-1][x] = Wall
	}
	return board
}

// Return true if the y-th row is filled
func (b *Board) IsFilled(y int) bool {
	for _, cell := range b[y] {
		if cell == nil {
			return false
		}
	}
	return true
}

// Return true if the mino is collided with the board
func (b *Board) isCollided(mino Mino) bool {
	for dy := range len(mino.Shape()) {
		for dx := range len(mino.Shape()[dy]) {
			if mino.Shape()[dy][dx] == 0 {
				continue
			}
			ny, nx := mino.Y+dy, mino.X+dx
			if ny < 0 || ny >= OUTER_HEIGHT || nx < 0 || nx >= OUTER_WIDTH {
				return true
			}
			if b[ny][nx] != nil || b[ny][nx] == Wall {
				return true
			}
		}
	}
	return false
}

// Write the color of mino to the board at each position
func (b *Board) Fix(mino *Mino) {
	for dy := range len(mino.Shape()) {
		for dx := range len(mino.Shape()[dy]) {
			if mino.Shape()[dy][dx] == 0 {
				continue
			}
			b[mino.Y+dy][mino.X+dx] = mino.Color
		}
	}
}

// Clear the filled lines and return the number of cleared lines
func (b *Board) ClearLines() (clearedLines []int, clearedColors [][12]color.Color) {
	clearedLines = []int{}
	clearedColors = [][12]color.Color{}

	newBoard := NewBoard()

	for y := MARGIN + INNER_HEIGHT - SENTINEL; y >= 0; y-- {
		if b.IsFilled(y) {
			clearedLines = append(clearedLines, y)
			clearedColors = append(clearedColors, b[y])
			continue
		}
		newBoard[y+len(clearedLines)] = b[y]
	}
	for y := range OUTER_HEIGHT {
		b[y] = newBoard[y]
	}
	return
}

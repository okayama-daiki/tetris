package game

import "image/color"

var (
	Wall             = color.RGBA{108, 122, 137, 255}
	BACKGROUND_COLOR = color.RGBA{5, 5, 5, 255}
	LINE_COLOR       = color.RGBA{75, 75, 75, 255}
	BORDER_COLOR     = color.RGBA{240, 240, 240, 255}
)

type Row = [WIDTH + 2]color.Color
type Board [HEIGHT + 4]Row

func IsFilled(row Row) bool {
	for _, cell := range row {
		if cell == nil {
			return false
		}
	}
	return true
}

func (b *Board) Init() {
	for y := range HEIGHT + 4 {
		b[y][0] = Wall
		b[y][WIDTH+1] = Wall
	}
	for i := range WIDTH + 2 {
		b[HEIGHT+3][i] = Wall
	}
}

func (b *Board) ClearLines() (clearedLines int) {
	clearedLines = 0
	for y := HEIGHT + 2; y > 0; y-- {
		if IsFilled(b[y]) {
			clearedLines++
			for i := y; i > 2; i-- {
				b[i] = b[i-1]
			}
			b[0] = Row{}
			b[0][0] = Wall
			b[0][WIDTH+1] = Wall
		}
	}
	return
}

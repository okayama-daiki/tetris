package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	HEIGHT    = 20
	WIDTH     = 10
	CELL_SIZE = 20
	MARGIN    = 0.5
)

type Game struct {
	PutPieces            int
	ClearedLines         int
	FrameCount           int
	NormalDroppingSpeed  int
	CurrentDroppingSpeed int
	Board                Board
	CurrentMino          Mino
	HoldingMino          HoldingMino
	MinoBag              MinoBag
}

func isCollided(board Board, mino Mino) bool {
	for dy := 0; dy < len(mino.Shape); dy++ {
		for dx := 0; dx < len(mino.Shape[dy]); dx++ {
			if mino.Shape[dy][dx] == 0 {
				continue
			}
			if board[mino.Y+dy][mino.X+dx] != nil || board[mino.Y+dy][mino.X+dx] == Wall {
				return true
			}
		}
	}
	return false
}

func (g *Game) Update() error {
	g.FrameCount++
	g.CurrentMino.FrameCount++
	g.CurrentDroppingSpeed = g.NormalDroppingSpeed

	switch {

	// Hold
	case inpututil.IsKeyJustPressed(ebiten.KeyC) && g.HoldingMino.Available:
		if g.HoldingMino.Mino.Name == "" {
			g.HoldingMino.Mino = g.MinoBag.Next()
		}
		g.CurrentMino.Y, g.CurrentMino.X = 0, 4
		g.HoldingMino.Mino, g.CurrentMino = g.CurrentMino, g.HoldingMino.Mino
		g.HoldingMino.Available = false

	// Hard drop
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		for nextMino := g.CurrentMino.MoveDown(); !isCollided(g.Board, nextMino); nextMino = nextMino.MoveDown() {
			g.CurrentMino = nextMino
		}
		g.CurrentMino.IsGrounded = true
		for dy := 0; dy < len(g.CurrentMino.Shape); dy++ {
			for dx := 0; dx < len(g.CurrentMino.Shape[dy]); dx++ {
				if g.CurrentMino.Shape[dy][dx] == 0 {
					continue
				}
				g.Board[g.CurrentMino.Y+dy][g.CurrentMino.X+dx] = g.CurrentMino.Color
			}
		}
		g.Board.ClearLines()
		g.CurrentMino = g.MinoBag.Next()
		g.HoldingMino.Available = true

	// Move Left
	case inpututil.KeyPressDuration(ebiten.KeyLeft) > 9 && inpututil.KeyPressDuration(ebiten.KeyLeft)%2 == 0 || inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		nextMino := g.CurrentMino.MoveLeft()
		if !isCollided(g.Board, nextMino) {
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			g.CurrentMino = nextMino
		}

	// Move Right
	case inpututil.KeyPressDuration(ebiten.KeyRight) > 9 && inpututil.KeyPressDuration(ebiten.KeyRight)%2 == 0 || inpututil.IsKeyJustPressed(ebiten.KeyRight):
		nextMino := g.CurrentMino.MoveRight()
		if !isCollided(g.Board, nextMino) {
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			g.CurrentMino = nextMino
		}

	// Rotate right
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyX):
		nextMino := g.CurrentMino.RotateRight()
		if !isCollided(g.Board, nextMino) {
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			g.CurrentMino = nextMino
		}

	// Rotate left
	case inpututil.IsKeyJustPressed(ebiten.KeyZ):
		nextMino := g.CurrentMino.RotateLeft()
		if !isCollided(g.Board, nextMino) {
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			g.CurrentMino = nextMino
		}

	// Soft drop
	case inpututil.KeyPressDuration(ebiten.KeyDown) > 0:
		g.CurrentDroppingSpeed = g.NormalDroppingSpeed / 20

	}

	switch {

	case g.CurrentMino.IsGrounded && g.CurrentMino.BacklashFrame == 0:
		for dy := 0; dy < len(g.CurrentMino.Shape); dy++ {
			for dx := 0; dx < len(g.CurrentMino.Shape[dy]); dx++ {
				if g.CurrentMino.Shape[dy][dx] == 0 {
					continue
				}
				g.Board[g.CurrentMino.Y+dy][g.CurrentMino.X+dx] = g.CurrentMino.Color
			}
		}
		g.Board.ClearLines()
		g.CurrentMino = g.MinoBag.Next()
		g.HoldingMino.Available = true

	case g.CurrentMino.IsGrounded:
		g.CurrentMino.BacklashFrame--

	case g.CurrentMino.FrameCount%g.CurrentDroppingSpeed == 0:
		nextMino := g.CurrentMino.MoveDown()
		if isCollided(g.Board, nextMino) {
			g.CurrentMino.IsGrounded = true
		} else {
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			g.CurrentMino = nextMino
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for y := 0; y < HEIGHT+4; y++ {
		for x := 0; x < WIDTH+2; x++ {
			if y == 0 && (x == 0 || x == WIDTH+1) || y == 1 && (x == 0 || x == WIDTH+1) {
				continue
			}
			c := g.Board[y][x]
			if c != nil {
				vector.DrawFilledRect(
					screen,
					float32(x*CELL_SIZE)+MARGIN,
					float32(y*CELL_SIZE)+MARGIN,
					CELL_SIZE-MARGIN*2,
					CELL_SIZE-MARGIN*2,
					c,
					true,
				)
			}
		}
	}

	// Ghost mino
	ghostMino := g.CurrentMino
	for ; !isCollided(g.Board, ghostMino.MoveDown()); ghostMino = ghostMino.MoveDown() {
	}
	for dy := 0; dy < len(ghostMino.Shape); dy++ {
		for dx := 0; dx < len(ghostMino.Shape[dy]); dx++ {
			if ghostMino.Shape[dy][dx] == 0 {
				continue
			}
			vector.DrawFilledRect(
				screen,
				float32((ghostMino.X+dx)*CELL_SIZE)+MARGIN,
				float32((ghostMino.Y+dy)*CELL_SIZE)+MARGIN,
				CELL_SIZE-MARGIN*2,
				CELL_SIZE-MARGIN*2,
				color.RGBA{50, 50, 50, 128},
				true,
			)
		}
	}

	// Dropping mino
	for dy := 0; dy < len(g.CurrentMino.Shape); dy++ {
		for dx := 0; dx < len(g.CurrentMino.Shape[dy]); dx++ {
			if g.CurrentMino.Shape[dy][dx] == 0 {
				continue
			}
			vector.DrawFilledRect(
				screen,
				float32((g.CurrentMino.X+dx)*CELL_SIZE)+MARGIN,
				float32((g.CurrentMino.Y+dy)*CELL_SIZE)+MARGIN,
				CELL_SIZE-MARGIN*2,
				CELL_SIZE-MARGIN*2,
				g.CurrentMino.Color,
				true,
			)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return CELL_SIZE * (WIDTH + 2), CELL_SIZE * (HEIGHT + 4)
}

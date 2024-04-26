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
	CELL_SIZE = 25
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
	for dy := range len(mino.Shape) {
		for dx := range len(mino.Shape[dy]) {
			if mino.Shape[dy][dx] == 0 {
				continue
			}
			ny, nx := mino.Y+dy, mino.X+dx
			if ny < 0 || ny >= HEIGHT+4 || nx < 0 || nx >= WIDTH+2 {
				return true
			}
			if board[ny][nx] != nil || board[ny][nx] == Wall {
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
		for dy := range len(g.CurrentMino.Shape) {
			for dx := range len(g.CurrentMino.Shape[dy]) {
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
		for _, nextMino := range g.CurrentMino.RotateRightSRS() {
			if !isCollided(g.Board, nextMino) {
				nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
				g.CurrentMino = nextMino
				break
			}
		}

	// Rotate left
	case inpututil.IsKeyJustPressed(ebiten.KeyZ):
		for _, nextMino := range g.CurrentMino.RotateLeftSSR() {
			if !isCollided(g.Board, nextMino) {
				nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
				g.CurrentMino = nextMino
				break
			}
		}

	// Soft drop
	case inpututil.KeyPressDuration(ebiten.KeyDown) > 0:
		g.CurrentDroppingSpeed = g.NormalDroppingSpeed / 20

	}

	switch {

	case g.CurrentMino.IsGrounded && g.CurrentMino.BacklashFrame == 0:
		for dy := range len(g.CurrentMino.Shape) {
			for dx := range len(g.CurrentMino.Shape[dy]) {
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
	var (
		boardWidth  = CELL_SIZE * (WIDTH + 2)
		boardHeight = CELL_SIZE * (HEIGHT + 4)
	)

	boardScreen := ebiten.NewImage(boardWidth, boardHeight)
	boardScreen.Fill(color.Black)

	for y := range HEIGHT + 4 {
		for x := range WIDTH + 2 {
			if y == 0 && (x == 0 || x == WIDTH+1) || y == 1 && (x == 0 || x == WIDTH+1) {
				continue
			}
			c := g.Board[y][x]
			if c != nil {
				vector.DrawFilledRect(
					boardScreen,
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
	for dy := range len(ghostMino.Shape) {
		for dx := range len(ghostMino.Shape[dy]) {
			if ghostMino.Shape[dy][dx] == 0 {
				continue
			}
			vector.DrawFilledRect(
				boardScreen,
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
	for dy := range len(g.CurrentMino.Shape) {
		for dx := range len(g.CurrentMino.Shape[dy]) {
			if g.CurrentMino.Shape[dy][dx] == 0 {
				continue
			}
			vector.DrawFilledRect(
				boardScreen,
				float32((g.CurrentMino.X+dx)*CELL_SIZE)+MARGIN,
				float32((g.CurrentMino.Y+dy)*CELL_SIZE)+MARGIN,
				CELL_SIZE-MARGIN*2,
				CELL_SIZE-MARGIN*2,
				g.CurrentMino.Color,
				true,
			)
		}
	}

	// Holding mino
	holdingMinoScreen := ebiten.NewImage(CELL_SIZE*6, CELL_SIZE*6)
	vector.DrawFilledRect(holdingMinoScreen, 0, 0, CELL_SIZE*6, CELL_SIZE*6, color.RGBA{0, 0, 0, 255}, true)

	if g.HoldingMino.Mino.Name != "" {
		for dy := range len(g.HoldingMino.Mino.Shape) {
			for dx := range len(g.HoldingMino.Mino.Shape[dy]) {
				if g.HoldingMino.Mino.Shape[dy][dx] == 0 {
					continue
				}
				var c color.Color = color.RGBA{30, 30, 30, 255}
				if g.HoldingMino.Available {
					c = g.HoldingMino.Mino.Color
				}
				vector.DrawFilledRect(
					holdingMinoScreen,
					float32(dx*CELL_SIZE)+MARGIN,
					float32(dy*CELL_SIZE)+MARGIN,
					CELL_SIZE-MARGIN*2,
					CELL_SIZE-MARGIN*2,
					c,
					true,
				)
			}
		}
	}

	// Next minos
	nextMinosScreen := ebiten.NewImage(CELL_SIZE*6, CELL_SIZE*(HEIGHT+4))
	for i, mino := range g.MinoBag.Sniff(6) {
		for dy := range len(mino.Shape) {
			for dx := range len(mino.Shape[dy]) {
				if mino.Shape[dy][dx] == 0 {
					continue
				}
				vector.DrawFilledRect(
					nextMinosScreen,
					float32(dx*CELL_SIZE)+MARGIN,
					float32((dy+i*3)*CELL_SIZE)+MARGIN,
					CELL_SIZE-MARGIN*2,
					CELL_SIZE-MARGIN*2,
					mino.Color,
					true,
				)
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 2*CELL_SIZE)
	screen.DrawImage(holdingMinoScreen, op)
	op.GeoM.Translate(6*CELL_SIZE, -2*CELL_SIZE)
	screen.DrawImage(boardScreen, op)
	op.GeoM.Translate(float64(boardWidth), 2*CELL_SIZE)
	screen.DrawImage(nextMinosScreen, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 600
}

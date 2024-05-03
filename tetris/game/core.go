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

func drawBlock(screen *ebiten.Image, x, y int, c color.Color, size, margin float32) {
	r, g, b, _ := c.RGBA()

	vector.DrawFilledRect(
		screen,
		float32(x)*size+margin,
		float32(y)*size+margin,
		size-2*margin,
		size-2*margin,
		color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 255},
		true,
	)
	vector.DrawFilledRect(
		screen,
		float32(x)*size+10*margin,
		float32(y)*size+10*margin,
		size-20*margin,
		size-20*margin,
		color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 230},
		true,
	)
	vector.DrawFilledRect(
		screen,
		float32(x)*size+13*margin,
		float32(y)*size+13*margin,
		size-26*margin,
		size-26*margin,
		color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 255},
		true,
	)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(BACKGROUND_COLOR)

	var (
		boardWidth  = CELL_SIZE * (WIDTH + 2)
		boardHeight = CELL_SIZE * (HEIGHT + 4)
	)

	boardScreen := ebiten.NewImage(boardWidth, boardHeight)
	boardScreen.Fill(BACKGROUND_COLOR)

	// Horizontal Lines
	for y := 2; y < HEIGHT+3; y++ {
		vector.StrokeLine(
			boardScreen,
			CELL_SIZE,
			float32(y*CELL_SIZE)+2,
			float32(WIDTH+1)*CELL_SIZE,
			float32(y*CELL_SIZE)+2,
			0.5,
			LINE_COLOR,
			true,
		)
	}

	// Vertical Lines
	for x := 0; x < WIDTH+2; x++ {
		vector.StrokeLine(
			boardScreen,
			float32(x*CELL_SIZE),
			2*CELL_SIZE,
			float32(x*CELL_SIZE),
			float32(HEIGHT+3)*CELL_SIZE,
			0.5,
			LINE_COLOR,
			true,
		)
	}

	// Border
	vector.StrokeLine(
		boardScreen,
		CELL_SIZE,
		2*CELL_SIZE,
		CELL_SIZE,
		float32(HEIGHT+3)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)
	vector.StrokeLine(
		boardScreen,
		float32(WIDTH+1)*CELL_SIZE,
		2*CELL_SIZE,
		float32(WIDTH+1)*CELL_SIZE,
		float32(HEIGHT+3)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)
	vector.StrokeLine(
		boardScreen,
		CELL_SIZE,
		float32(HEIGHT+3)*CELL_SIZE,
		float32(WIDTH+1)*CELL_SIZE,
		float32(HEIGHT+3)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)

	// Fixed minos
	for y := 2; y < HEIGHT+3; y++ {
		for x := 1; x < WIDTH+1; x++ {
			c := g.Board[y][x]
			if c != nil {
				drawBlock(boardScreen, x, y, c, CELL_SIZE, MARGIN)
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
			drawBlock(boardScreen, ghostMino.X+dx, ghostMino.Y+dy, GHOST_COLOR, CELL_SIZE, MARGIN)
		}
	}

	// Dropping mino
	for dy := range len(g.CurrentMino.Shape) {
		for dx := range len(g.CurrentMino.Shape[dy]) {
			if g.CurrentMino.Shape[dy][dx] == 0 {
				continue
			}
			drawBlock(boardScreen, g.CurrentMino.X+dx, g.CurrentMino.Y+dy, g.CurrentMino.Color, CELL_SIZE, MARGIN)
		}
	}

	// Holding mino
	holdingMinoScreen := ebiten.NewImage(CELL_SIZE*6, CELL_SIZE*6)
	holdingMinoScreen.Fill(BACKGROUND_COLOR)

	if g.HoldingMino.Mino.Name != "" {
		for dy := range len(g.HoldingMino.Mino.Shape) {
			for dx := range len(g.HoldingMino.Mino.Shape[dy]) {
				if g.HoldingMino.Mino.Shape[dy][dx] == 0 {
					continue
				}
				var c color.Color = GHOST_COLOR
				if g.HoldingMino.Available {
					c = g.HoldingMino.Mino.Color
				}
				drawBlock(holdingMinoScreen, dx+2, dy, c, CELL_SIZE, MARGIN)
			}
		}
	}

	// Next minos
	nextMinosScreen := ebiten.NewImage(CELL_SIZE*6, CELL_SIZE*(HEIGHT+4))
	nextMinosScreen.Fill(BACKGROUND_COLOR)

	for i, mino := range g.MinoBag.Sniff(6) {
		for dy := range len(mino.Shape) {
			for dx := range len(mino.Shape[dy]) {
				if mino.Shape[dy][dx] == 0 {
					continue
				}
				drawBlock(nextMinosScreen, dx, dy+i*3, mino.Color, CELL_SIZE, MARGIN)
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

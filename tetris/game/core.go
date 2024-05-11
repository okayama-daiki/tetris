package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	INNER_HEIGHT = 20
	INNER_WIDTH  = 10
	SENTINEL     = 1
	MARGIN       = 3
	OUTER_HEIGHT = MARGIN + INNER_HEIGHT + SENTINEL
	OUTER_WIDTH  = SENTINEL + INNER_WIDTH + SENTINEL
)

const (
	CELL_SIZE = 25
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
	Fragments            [OUTER_HEIGHT][OUTER_WIDTH]Fragment
}

func (g *Game) restart() {
	g.HoldingMino = HoldingMino{Available: true}
	g.Board = NewBoard()
	g.MinoBag = MinoBag{}
	g.CurrentMino = g.MinoBag.Next()
}

func (g *Game) Update() error {
	g.FrameCount++
	g.CurrentMino.FrameCount++
	g.CurrentMino.LockDown.UpdateTimer()
	g.CurrentDroppingSpeed = g.NormalDroppingSpeed

	if inpututil.KeyPressDuration(ebiten.KeyR) == 30 {
		g.restart()
	}

	// Hold
	if inpututil.IsKeyJustPressed(ebiten.KeyC) && g.HoldingMino.Available {
		if g.HoldingMino.Mino.Name == "" {
			g.HoldingMino.Mino = g.MinoBag.Next()
		}
		g.CurrentMino.Y, g.CurrentMino.X = 0, 4
		g.HoldingMino.Mino, g.CurrentMino = g.CurrentMino, g.HoldingMino.Mino
		g.HoldingMino.Available = false
	}
	// Hard drop
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		for nextMino := g.CurrentMino.MoveDown(); !g.Board.isCollided(nextMino); nextMino = nextMino.MoveDown() {
			g.CurrentMino = nextMino
		}
		g.CurrentMino.IsGrounded = true
		g.Board.Fix(&g.CurrentMino)
		clearedLines, clearedColors := g.Board.ClearLines()
		for i, y := range clearedLines {
			g.ClearedLines++
			for x := range OUTER_WIDTH {
				g.Fragments[y][x] = NewFragment(clearedColors[i][x], x, y)
			}
		}
		g.CurrentMino = g.MinoBag.Next()
		g.HoldingMino.Available = true
	}

	// Move Left
	if inpututil.KeyPressDuration(ebiten.KeyLeft) > 9 && inpututil.KeyPressDuration(ebiten.KeyLeft)%2 == 0 || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		nextMino := g.CurrentMino.MoveLeft()
		if !g.Board.isCollided(nextMino) {
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			nextMino.IsGrounded = false
			nextMino.LockDown.UpdateCounter()
			g.CurrentMino = nextMino
		}
	}

	// Move Right
	if inpututil.KeyPressDuration(ebiten.KeyRight) > 9 && inpututil.KeyPressDuration(ebiten.KeyRight)%2 == 0 || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		nextMino := g.CurrentMino.MoveRight()
		if !g.Board.isCollided(nextMino) {
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			nextMino.IsGrounded = false
			nextMino.LockDown.UpdateCounter()
			g.CurrentMino = nextMino
		}
	}

	// Rotate right
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		for _, nextMino := range g.CurrentMino.RotateRightSRS() {
			if !g.Board.isCollided(nextMino) {
				nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
				nextMino.IsGrounded = false
				nextMino.LockDown.UpdateCounter()
				g.CurrentMino = nextMino
				break
			}
		}
	}

	// Rotate left
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		for _, nextMino := range g.CurrentMino.RotateLeftSSR() {
			if !g.Board.isCollided(nextMino) {
				nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
				nextMino.IsGrounded = false
				nextMino.LockDown.UpdateCounter()
				g.CurrentMino = nextMino
				break
			}
		}
	}

	// Soft drop
	if inpututil.KeyPressDuration(ebiten.KeyDown) > 0 {
		g.CurrentDroppingSpeed = g.NormalDroppingSpeed / 20
	}

	switch {

	case g.CurrentMino.LockDown.IsFixed():
		for nextMino := g.CurrentMino.MoveDown(); !g.Board.isCollided(nextMino); nextMino = nextMino.MoveDown() {
			g.CurrentMino = nextMino
		}
		g.Board.Fix(&g.CurrentMino)
		clearedLines, clearedColors := g.Board.ClearLines()
		for i, y := range clearedLines {
			g.ClearedLines++
			for x := range OUTER_WIDTH {
				g.Fragments[y][x] = NewFragment(clearedColors[i][x], x, y)
			}
		}
		g.CurrentMino = g.MinoBag.Next()
		g.HoldingMino.Available = true

	case g.CurrentMino.FrameCount%g.CurrentDroppingSpeed == 0:
		nextMino := g.CurrentMino.MoveDown()
		if !g.Board.isCollided(nextMino) {
			nextMino.IsGrounded = false
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			nextMino.LockDown.Reset()
			g.CurrentMino = nextMino
		} else {
			g.CurrentMino.LockDown.Activate()
		}
	}

	return nil
}

func drawBlock(screen *ebiten.Image, x, y int, c color.Color, size float32) {
	var margin float32 = 0.5
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
		boardWidth  = CELL_SIZE * OUTER_WIDTH
		boardHeight = CELL_SIZE * OUTER_HEIGHT
	)

	boardScreen := ebiten.NewImage(boardWidth, boardHeight)
	boardScreen.Fill(BACKGROUND_COLOR)

	// Animation
	for y := range OUTER_HEIGHT {
		for x := range OUTER_WIDTH {
			if g.Fragments[y][x].Frame > 0 {
				g.Fragments[y][x].Frame--
				posX, posY := g.Fragments[y][x].Position()
				vector.DrawFilledRect(boardScreen, posX, posY, CELL_SIZE/2, CELL_SIZE/2, g.Fragments[y][x].Color(), true)
			}
		}
	}

	// Horizontal Lines
	for y := MARGIN; y < OUTER_HEIGHT; y++ {
		vector.StrokeLine(
			boardScreen,
			CELL_SIZE,
			float32(y*CELL_SIZE)+2,
			float32(INNER_WIDTH+SENTINEL)*CELL_SIZE,
			float32(y*CELL_SIZE)+2,
			0.5,
			LINE_COLOR,
			true,
		)
	}

	// Vertical Lines
	for x := range OUTER_WIDTH {
		vector.StrokeLine(
			boardScreen,
			float32(x*CELL_SIZE),
			MARGIN*CELL_SIZE,
			float32(x*CELL_SIZE),
			float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
			0.5,
			LINE_COLOR,
			true,
		)
	}

	// Border
	vector.StrokeLine(
		boardScreen,
		CELL_SIZE,
		MARGIN*CELL_SIZE,
		CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)
	vector.StrokeLine(
		boardScreen,
		float32(SENTINEL+INNER_WIDTH)*CELL_SIZE,
		MARGIN*CELL_SIZE,
		float32(SENTINEL+INNER_WIDTH)*CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)
	vector.StrokeLine(
		boardScreen,
		CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		float32(INNER_WIDTH+SENTINEL)*CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)

	// Fixed minos
	for y := 0; y < MARGIN+INNER_HEIGHT; y++ {
		for x := SENTINEL; x < INNER_WIDTH+SENTINEL; x++ {
			c := g.Board[y][x]
			if c != nil {
				drawBlock(boardScreen, x, y, c, CELL_SIZE)
			}
		}
	}

	// Ghost mino
	ghostMino := g.CurrentMino
	for ; !g.Board.isCollided(ghostMino.MoveDown()); ghostMino = ghostMino.MoveDown() {
	}
	for dy := range len(ghostMino.Shape()) {
		for dx := range len(ghostMino.Shape()[dy]) {
			if ghostMino.Shape()[dy][dx] == 0 {
				continue
			}
			drawBlock(boardScreen, ghostMino.X+dx, ghostMino.Y+dy, GHOST_COLOR, CELL_SIZE)
		}
	}

	// Dropping mino
	for dy := range len(g.CurrentMino.Shape()) {
		for dx := range len(g.CurrentMino.Shape()[dy]) {
			if g.CurrentMino.Shape()[dy][dx] == 0 {
				continue
			}
			drawBlock(boardScreen, g.CurrentMino.X+dx, g.CurrentMino.Y+dy, g.CurrentMino.Color, CELL_SIZE)
		}
	}

	// Holding mino
	holdingMinoScreen := ebiten.NewImage(CELL_SIZE*6, CELL_SIZE*6)
	holdingMinoScreen.Fill(BACKGROUND_COLOR)

	if g.HoldingMino.Mino.Name != "" {
		for dy := range len(g.HoldingMino.Mino.Shape()) {
			for dx := range len(g.HoldingMino.Mino.Shape()[dy]) {
				if g.HoldingMino.Mino.Shape()[dy][dx] == 0 {
					continue
				}
				var c color.Color = GHOST_COLOR
				if g.HoldingMino.Available {
					c = g.HoldingMino.Mino.Color
				}
				drawBlock(holdingMinoScreen, dx+2, dy, c, CELL_SIZE)
			}
		}
	}

	// Next minos
	nextMinosScreen := ebiten.NewImage(CELL_SIZE*6, CELL_SIZE*(OUTER_HEIGHT))
	nextMinosScreen.Fill(BACKGROUND_COLOR)

	for i, mino := range g.MinoBag.Sniff(6) {
		for dy := range len(mino.Shape()) {
			for dx := range len(mino.Shape()[dy]) {
				if mino.Shape()[dy][dx] == 0 {
					continue
				}
				drawBlock(nextMinosScreen, dx, dy+i*3, mino.Color, CELL_SIZE)
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %f\ntps: %f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 600
}

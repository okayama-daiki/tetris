package game

import (
	_ "embed"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/okayama-daiki/tetris/tetris/audio"
)

const (
	INNER_HEIGHT  = 20
	INNER_WIDTH   = 10
	SENTINEL_SIZE = 1
	CELL_SIZE     = 25
	MARGIN        = 3
	OUTER_HEIGHT  = MARGIN + INNER_HEIGHT + SENTINEL_SIZE
	OUTER_WIDTH   = SENTINEL_SIZE + INNER_WIDTH + SENTINEL_SIZE
)

const (
	KEY_LONG_PRESS_WAIT_TIME = 9
	KEY_PRESS_DURATION       = 2
)

const (
	MAX_LEVEL = 110
)

var fontFace = text.NewGoXFace(bitmapfont.Face)

func NewGame(audioPlayer *audio.Player) *Game {
	minoBag := MinoBag{}
	return &Game{
		MinoBag:              minoBag,
		Board:                NewBoard(),
		CurrentMino:          minoBag.Next(),
		HoldingMino:          HoldingMino{Available: true},
		AudioPlayer:          audioPlayer,
		CurrentDroppingSpeed: 60,
		NormalDroppingSpeed:  60,
	}
}

type Game struct {
	PutPieces            int
	ClearedLines         int
	FrameCount           int
	NormalDroppingSpeed  int
	CurrentDroppingSpeed int
	level                int
	Board                Board
	CurrentMino          Mino
	HoldingMino          HoldingMino
	MinoBag              MinoBag
	Fragments            [OUTER_HEIGHT][OUTER_WIDTH]Fragment
	AudioPlayer          *audio.Player
}

func (g *Game) restart() {
	g.AudioPlayer.PlayClear()
	g.PutPieces = 0
	g.ClearedLines = 0
	g.FrameCount = 0
	g.Fragments = [OUTER_HEIGHT][OUTER_WIDTH]Fragment{}
	for y := range OUTER_HEIGHT {
		for x := range OUTER_WIDTH {
			if g.Board[y][x] != nil {
				g.Fragments[y][x] = NewFragment(g.Board[y][x], x, y)
			}
		}
	}
	g.HoldingMino = HoldingMino{Available: true}
	g.Board = NewBoard()
	g.MinoBag = MinoBag{}
	g.CurrentMino = g.MinoBag.Next()
}

func (g *Game) Update() error {
	g.AudioPlayer.Update()
	g.FrameCount++
	g.CurrentMino.FrameCount++
	g.CurrentMino.LockDown.UpdateTimer()
	g.level = min(g.ClearedLines/10+1, MAX_LEVEL)
	g.CurrentDroppingSpeed = max(int((0.8-float64(g.level-1)*0.05)*60), 1)

	if inpututil.KeyPressDuration(ebiten.KeyR) == 30 {
		g.restart()
	}

	// Hold
	if inpututil.IsKeyJustPressed(ebiten.KeyC) && g.HoldingMino.Available {
		if g.HoldingMino.Mino.Name == "" {
			g.HoldingMino.Mino = g.MinoBag.Next()
		}
		g.AudioPlayer.PlayHold()
		g.CurrentMino.Y, g.CurrentMino.X = 0, 4
		g.CurrentMino.Angle = 0
		g.HoldingMino.Mino, g.CurrentMino = g.CurrentMino, g.HoldingMino.Mino
		g.HoldingMino.Available = false
	}

	// Hard drop
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.AudioPlayer.PlayHardDrop()
		for nextMino := g.CurrentMino.MoveDown(); !g.Board.isCollided(nextMino); nextMino = nextMino.MoveDown() {
			g.CurrentMino = nextMino
		}
		g.CurrentMino.IsGrounded = true
		g.Board.Fix(&g.CurrentMino)
		clearedLines, clearedColors := g.Board.ClearLines()
		if len(clearedLines) > 0 {
			g.AudioPlayer.PlayClear()
		}
		for i, y := range clearedLines {
			g.ClearedLines++
			for x := range OUTER_WIDTH {
				g.Fragments[y][x] = NewFragment(clearedColors[i][x], x, y)
			}
		}
		g.PutPieces++
		g.CurrentMino = g.MinoBag.Next()
		if g.IsGameOver() {
			g.restart()
		}
		g.HoldingMino.Available = true

	}

	// Move Left
	if inpututil.KeyPressDuration(ebiten.KeyLeft) > KEY_LONG_PRESS_WAIT_TIME &&
		inpututil.KeyPressDuration(ebiten.KeyLeft)%KEY_PRESS_DURATION == 0 || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		nextMino := g.CurrentMino.MoveLeft()
		if !g.Board.isCollided(nextMino) {
			g.AudioPlayer.PlayMove()
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			nextMino.IsGrounded = false
			nextMino.LockDown.UpdateCounter()
			g.CurrentMino = nextMino
		}
	}

	// Move Right
	if inpututil.KeyPressDuration(ebiten.KeyRight) > KEY_LONG_PRESS_WAIT_TIME &&
		inpututil.KeyPressDuration(ebiten.KeyRight)%KEY_PRESS_DURATION == 0 || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		nextMino := g.CurrentMino.MoveRight()
		if !g.Board.isCollided(nextMino) {
			g.AudioPlayer.PlayMove()
			nextMino.BacklashFrame = DEFAULT_BACKLASH_FRAME
			nextMino.IsGrounded = false
			nextMino.LockDown.UpdateCounter()
			g.CurrentMino = nextMino
		}
	}

	// Rotate right
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		for nextMino := range g.CurrentMino.RotateRightSRS() {
			if !g.Board.isCollided(nextMino) {
				g.AudioPlayer.PlayRotate()
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
		for nextMino := range g.CurrentMino.RotateLeftSSR() {
			if !g.Board.isCollided(nextMino) {
				g.AudioPlayer.PlayRotate()
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
		g.CurrentDroppingSpeed = max(g.NormalDroppingSpeed/20, 1)
	}

	switch {

	case g.CurrentMino.LockDown.IsFixed():
		for nextMino := g.CurrentMino.MoveDown(); !g.Board.isCollided(nextMino); nextMino = nextMino.MoveDown() {
			g.CurrentMino = nextMino
		}
		g.Board.Fix(&g.CurrentMino)
		clearedLines, clearedColors := g.Board.ClearLines()
		if len(clearedLines) > 0 {
			g.AudioPlayer.PlayClear()
		}
		for i, y := range clearedLines {
			g.ClearedLines++
			for x := range OUTER_WIDTH {
				g.Fragments[y][x] = NewFragment(clearedColors[i][x], x, y)
			}
		}
		g.PutPieces++
		g.CurrentMino = g.MinoBag.Next()
		if g.IsGameOver() {
			g.restart()
		}
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

func (g *Game) IsGameOver() bool {
	return g.CurrentMino.Y == 0 && g.Board.isCollided(g.CurrentMino)
}

func MakeDrawFilledRect(offsetX, offsetY float32) func(screen *ebiten.Image, x, y, width, height float32, clr color.Color, antialias bool) {
	return func(screen *ebiten.Image, x, y, width, height float32, clr color.Color, antialias bool) {
		vector.DrawFilledRect(screen, x+offsetX, y+offsetY, width, height, clr, antialias)
	}
}

func MakeStrokeLine(offsetX, offsetY float32) func(dst *ebiten.Image, x0, y0, x1, y1, strokeWidth float32, clr color.Color, antialias bool) {
	return func(dst *ebiten.Image, x0, y0, x1, y1, strokeWidth float32, clr color.Color, antialias bool) {
		vector.StrokeLine(dst, x0+offsetX, y0+offsetY, x1+offsetX, y1+offsetY, strokeWidth, clr, antialias)
	}
}

func MakeDrawBlock(offsetX, offsetY float32) func(screen *ebiten.Image, x, y int, c color.Color, size float32) {
	return func(screen *ebiten.Image, x, y int, c color.Color, size float32) {
		var padding float32 = 0.5
		r, g, b, _ := c.RGBA()

		vector.DrawFilledRect(
			screen,
			float32(x)*size+padding+offsetX,
			float32(y)*size+padding+offsetY,
			size-2*padding,
			size-2*padding,
			color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 255},
			true,
		)
		vector.DrawFilledRect(
			screen,
			float32(x)*size+10*padding+offsetX,
			float32(y)*size+10*padding+offsetY,
			size-20*padding,
			size-20*padding,
			color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 230},
			true,
		)
		vector.DrawFilledRect(
			screen,
			float32(x)*size+13*padding+offsetX,
			float32(y)*size+13*padding+offsetY,
			size-26*padding,
			size-26*padding,
			color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 255},
			true,
		)
	}
}

func (g *Game) drawGameBoard(screen *ebiten.Image, offsetX, offsetY float32) {
	drawFilledRect := MakeDrawFilledRect(offsetX, offsetY)
	strokeLine := MakeStrokeLine(offsetX, offsetY)
	drawBlock := MakeDrawBlock(offsetX, offsetY)

	// Animation
	for y := range OUTER_HEIGHT {
		for x := range OUTER_WIDTH {
			if g.Fragments[y][x].Frame > 0 {
				g.Fragments[y][x].Frame--
				posX, posY := g.Fragments[y][x].Position()
				drawFilledRect(
					screen,
					posX,
					posY,
					CELL_SIZE/2,
					CELL_SIZE/2,
					g.Fragments[y][x].Color(),
					true,
				)
			}
		}
	}

	// Horizontal Lines
	for y := MARGIN; y < OUTER_HEIGHT; y++ {
		strokeLine(
			screen,
			CELL_SIZE,
			float32(y*CELL_SIZE)+2,
			float32(INNER_WIDTH+SENTINEL_SIZE)*CELL_SIZE,
			float32(y*CELL_SIZE)+2,
			0.5,
			LINE_COLOR,
			true,
		)
	}

	// Vertical Lines
	for x := SENTINEL_SIZE; x < OUTER_WIDTH; x++ {
		strokeLine(
			screen,
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
	strokeLine(
		screen,
		CELL_SIZE,
		MARGIN*CELL_SIZE,
		CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)
	strokeLine(
		screen,
		float32(SENTINEL_SIZE+INNER_WIDTH)*CELL_SIZE,
		MARGIN*CELL_SIZE,
		float32(SENTINEL_SIZE+INNER_WIDTH)*CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)
	strokeLine(
		screen,
		CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		float32(INNER_WIDTH+SENTINEL_SIZE)*CELL_SIZE,
		float32(MARGIN+INNER_HEIGHT)*CELL_SIZE,
		2,
		BORDER_COLOR,
		true,
	)

	// Fixed minos
	for y := 0; y < MARGIN+INNER_HEIGHT; y++ {
		for x := SENTINEL_SIZE; x < INNER_WIDTH+SENTINEL_SIZE; x++ {
			c := g.Board[y][x]
			if c != nil {
				drawBlock(screen, x, y, c, CELL_SIZE)
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
			drawBlock(screen, ghostMino.X+dx, ghostMino.Y+dy, GHOST_COLOR, CELL_SIZE)
		}
	}

	// Dropping mino
	for dy := range len(g.CurrentMino.Shape()) {
		for dx := range len(g.CurrentMino.Shape()[dy]) {
			if g.CurrentMino.Shape()[dy][dx] == 0 {
				continue
			}
			drawBlock(screen, g.CurrentMino.X+dx, g.CurrentMino.Y+dy, g.CurrentMino.Color, CELL_SIZE)
		}
	}
}

func (g *Game) drawHold(screen *ebiten.Image, offsetX, offsetY float32) {
	drawBlock := MakeDrawBlock(offsetX, offsetY)

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
				drawBlock(screen, dx+2, dy, c, CELL_SIZE)
			}
		}
	}
}

func (g *Game) drawNext(screen *ebiten.Image, offsetX, offsetY float32) {
	drawBlock := MakeDrawBlock(offsetX, offsetY)

	for i, mino := range g.MinoBag.Sniff(6) {
		for dy := range len(mino.Shape()) {
			for dx := range len(mino.Shape()[dy]) {
				if mino.Shape()[dy][dx] == 0 {
					continue
				}
				drawBlock(screen, dx, dy+i*3, mino.Color, CELL_SIZE)
			}
		}
	}
}

func (g *Game) drawController(screen *ebiten.Image, offsetX, offsetY float32) {
	option := &text.DrawOptions{LayoutOptions: text.LayoutOptions{LineSpacing: 20}}
	option.GeoM.Translate(float64(offsetX), float64(offsetY))
	text.Draw(screen,
		`
←      : Move Left
→      : Move Right
Z      : Rotate Left
X(↑)   : Rotate Right
C      : Hold
Space  : Hard Drop
↓      : Soft Drop
`,
		fontFace,
		option,
	)
}

func (g *Game) drawScore(screen *ebiten.Image, offsetX, offsetY float32) {
	option := &text.DrawOptions{LayoutOptions: text.LayoutOptions{LineSpacing: 20}}
	option.GeoM.Translate(float64(offsetX), float64(offsetY))
	text.Draw(screen,
		fmt.Sprintf(
			`
Pieces : %d, %.02f/s
Lines  : %d
Time   : %d:%02d.%02d
Level	 : %d
`,
			g.PutPieces,
			float32(g.PutPieces)/float32(g.FrameCount/10)*6,
			g.ClearedLines,
			g.FrameCount/3600,
			g.FrameCount%3600/60,
			g.FrameCount%60,
			g.level,
		),
		fontFace,
		option,
	)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(BACKGROUND_COLOR)

	g.drawHold(screen, 0, 2*CELL_SIZE)
	g.drawGameBoard(screen, 6*CELL_SIZE, 0)
	g.drawNext(screen, (6+OUTER_WIDTH)*CELL_SIZE, 2*CELL_SIZE)
	g.drawController(screen, 30, 10*CELL_SIZE)
	g.drawScore(screen, 30, 18*CELL_SIZE)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %f\ntps: %f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 600
}

package game

import (
	"image/color"
	"math/rand"
)

const (
	DEFAULT_BACKLASH_FRAME = 30
)

type Shape [][]int

// Note: the Mino is fully fixed if IsGrounded is true and BacklashFrame is 0
type Mino struct {
	Name          string
	Shape         Shape
	Color         color.Color
	Y             int
	X             int
	Angle         int
	FrameCount    int
	IsGrounded    bool
	BacklashFrame int // Allow a little time for movement / rotation after grounding
}

type HoldingMino struct {
	Mino
	Available bool
}

func (m Mino) RotateRight() Mino {
	shape := make(Shape, len(m.Shape))
	for i := range m.Shape {
		shape[i] = make([]int, len(m.Shape[i]))
	}
	for y := range len(m.Shape) {
		for x := range len(m.Shape[y]) {
			shape[x][len(m.Shape)-1-y] = m.Shape[y][x]
		}
	}
	m.Shape = shape
	return m
}

func (m Mino) RotateLeft() Mino {
	shape := make(Shape, len(m.Shape))
	for i := range m.Shape {
		shape[i] = make([]int, len(m.Shape[i]))
	}
	for y := range len(m.Shape) {
		for x := range len(m.Shape[y]) {
			shape[len(m.Shape)-1-x][y] = m.Shape[y][x]
		}
	}
	m.Shape = shape
	return m
}

func (m Mino) MoveRight() Mino {
	m.X++
	return m
}

func (m Mino) MoveLeft() Mino {
	m.X--
	return m
}

func (m Mino) MoveDown() Mino {
	m.Y++
	return m
}

var (
	T = Mino{
		Shape: [][]int{
			{0, 1, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		Name:          "T",
		Color:         color.RGBA{128, 0, 128, 255},
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	O = Mino{
		Shape: [][]int{
			{1, 1},
			{1, 1},
		},
		Name:          "O",
		Color:         color.RGBA{255, 255, 0, 255},
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	L = Mino{
		Shape: [][]int{
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 0},
		},
		Name:          "L",
		Color:         color.RGBA{255, 127, 0, 255},
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	J = Mino{
		Shape: [][]int{
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		Name:          "J",
		Color:         color.RGBA{0, 0, 255, 255},
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	S = Mino{
		Shape: [][]int{
			{0, 1, 1},
			{1, 1, 0},
			{0, 0, 0},
		},
		Name:          "S",
		Color:         color.RGBA{0, 255, 0, 255},
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	Z = Mino{
		Shape: [][]int{
			{1, 1, 0},
			{0, 1, 1},
			{0, 0, 0},
		},
		Name:          "Z",
		Color:         color.RGBA{255, 0, 0, 255},
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	I = Mino{
		Shape: [][]int{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		Name:          "I",
		Color:         color.RGBA{0, 255, 255, 255},
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	Minos = []Mino{T, O, L, J, S, Z, I}
)

type MinoBag struct {
	queue []Mino
}

func (b *MinoBag) fill() {
	bag := make([]Mino, len(Minos))
	copy(bag, Minos)
	for i := range len(bag) {
		j := rand.Intn(i + 1)
		bag[i], bag[j] = bag[j], bag[i]
	}
	b.queue = append(b.queue, bag...)
}

func (b *MinoBag) Sniff(n int) []Mino {
	if n <= 0 || n > 7 {
		panic("n must be between 1 and 7")
	}
	if len(b.queue) < n {
		b.fill()
	}
	preview := make([]Mino, n)
	copy(preview, b.queue[:n])
	return preview
}

func (b *MinoBag) Next() Mino {
	if len(b.queue) == 0 {
		b.fill()
	}
	mino := b.queue[0]
	mino.Y, mino.X = 0, 4
	b.queue = b.queue[1:]
	return mino
}

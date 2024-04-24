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

func (m *Mino) RotateRight() Mino {
	shape := make(Shape, len(m.Shape))
	for i := range m.Shape {
		shape[i] = make([]int, len(m.Shape[i]))
	}
	for y := 0; y < len(m.Shape); y++ {
		for x := 0; x < len(m.Shape[y]); x++ {
			shape[x][len(m.Shape)-1-y] = m.Shape[y][x]
		}
	}
	return Mino{
		Shape:      shape,
		Name:       m.Name,
		Color:      m.Color,
		IsGrounded: m.IsGrounded,
		Y:          m.Y,
		X:          m.X,
		Angle:      (m.Angle + 1) % 4,
	}
}

func (m *Mino) RotateLeft() Mino {
	shape := make(Shape, len(m.Shape))
	for i := range m.Shape {
		shape[i] = make([]int, len(m.Shape[i]))
	}
	for y := 0; y < len(m.Shape); y++ {
		for x := 0; x < len(m.Shape[y]); x++ {
			shape[len(m.Shape)-1-x][y] = m.Shape[y][x]
		}
	}
	return Mino{
		Shape:         shape,
		Name:          m.Name,
		Color:         m.Color,
		IsGrounded:    m.IsGrounded,
		Y:             m.Y,
		X:             m.X,
		Angle:         (m.Angle + 3) % 4,
		FrameCount:    m.FrameCount,
		BacklashFrame: m.BacklashFrame,
	}
}

func (m *Mino) MoveRight() Mino {
	return Mino{
		Shape:         m.Shape,
		Name:          m.Name,
		Color:         m.Color,
		IsGrounded:    m.IsGrounded,
		Y:             m.Y,
		X:             m.X + 1,
		Angle:         m.Angle,
		FrameCount:    m.FrameCount,
		BacklashFrame: m.BacklashFrame,
	}
}

func (m *Mino) MoveLeft() Mino {
	return Mino{
		Shape:         m.Shape,
		Name:          m.Name,
		Color:         m.Color,
		IsGrounded:    m.IsGrounded,
		Y:             m.Y,
		X:             m.X - 1,
		Angle:         m.Angle,
		FrameCount:    m.FrameCount,
		BacklashFrame: m.BacklashFrame,
	}
}

func (m *Mino) MoveDown() Mino {
	return Mino{
		Shape:         m.Shape,
		Name:          m.Name,
		Color:         m.Color,
		IsGrounded:    m.IsGrounded,
		Y:             m.Y + 1,
		X:             m.X,
		Angle:         m.Angle,
		FrameCount:    m.FrameCount,
		BacklashFrame: m.BacklashFrame,
	}
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

type MinoIterator interface {
	HasNext() bool
	Next() Mino
}

type MinoBag struct {
	bag []Mino
}

func (b *MinoBag) HasNext() bool {
	return true
}

func (b *MinoBag) Next() Mino {
	if len(b.bag) == 0 {
		b.bag = make([]Mino, len(Minos))
		copy(b.bag, Minos)
		for i := 0; i < len(b.bag); i++ {
			j := rand.Intn(i + 1)
			b.bag[i], b.bag[j] = b.bag[j], b.bag[i]
		}
	}
	mino := b.bag[0]
	b.bag = b.bag[1:]
	mino.Y, mino.X = 0, 4
	return mino
}

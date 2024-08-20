package game

import (
	"image/color"
	"iter"
	"math/rand"
)

var (
	GHOST_COLOR = color.RGBA{30, 30, 30, 127}
)

var (
	PURPLE = color.RGBA{106, 50, 165, 255}
	YELLOW = color.RGBA{255, 213, 0, 255}
	ORANGE = color.RGBA{255, 121, 28, 255}
	BLUE   = color.RGBA{6, 119, 186, 255}
	GREEN  = color.RGBA{114, 203, 59, 255}
	RED    = color.RGBA{212, 42, 52, 255}
	CYAN   = color.RGBA{31, 195, 205, 255}
)

const (
	DEFAULT_BACKLASH_FRAME           = 30
	DEFAULT_EXTENDED_PLACEMENT_COUNT = 15
)

type Angle int

const (
	Angle0 Angle = iota
	Angle90
	Angle180
	Angle270
)

type Shape [][]int

// Note: the Mino is fully fixed if IsGrounded is true and BacklashFrame is 0 or ExtendedPlacementCounter is 0
type BaseMino struct {
	baseShape Shape
	y         int
	x         int
	angle     Angle
	color     color.Color
}

func NewBaseMino(shape Shape, color color.Color) BaseMino {
	return BaseMino{
		baseShape: shape,
		angle:     Angle0,
		color:     color,
	}
}

func (m BaseMino) Initialize() AbstractMino {
	m.y, m.x = 0, 4
	m.angle = Angle0
	return m
}

// Shape returns the current shape of the mino
func (m BaseMino) Shape() Shape {
	shape := make(Shape, len(m.baseShape))
	for i := range len(m.baseShape) {
		shape[i] = make([]int, len(m.baseShape[i]))
		copy(shape[i], m.baseShape[i])
	}
	for range m.angle {
		shape = Rotate(shape)
	}
	return shape
}

func (m BaseMino) Color() color.Color {
	return m.color
}

func (m BaseMino) X() int {
	return m.x
}

func (m BaseMino) Y() int {
	return m.y
}

func (m BaseMino) MoveRight() AbstractMino {
	m.x++
	return m
}

func (m BaseMino) MoveLeft() AbstractMino {
	m.x--
	return m
}

func (m BaseMino) MoveDown() AbstractMino {
	m.y++
	return m
}

func (m BaseMino) MoveUp() AbstractMino {
	m.y--
	return m
}

func (m BaseMino) rotateRight() AbstractMino {
	m.angle = (m.angle + 1) % 4
	return m
}

func (m BaseMino) rotateLeft() AbstractMino {
	m.angle = (m.angle + 3) % 4
	return m
}

// RotateRightSRS() yields the rotated minos according to the Super Rotation System.
//
// Note: This implementation is expected to use for the minos except I and O.
func (m BaseMino) RotateRightSRS() iter.Seq[AbstractMino] {
	switch m.angle {
	case Angle0:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveRight()) &&
				yield(m.rotateRight().MoveRight().MoveUp()) &&
				yield(m.rotateRight().MoveDown().MoveDown()) &&
				yield(m.rotateRight().MoveRight().MoveDown().MoveDown())
		}
	case Angle90:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveRight()) &&
				yield(m.rotateRight().MoveRight().MoveDown()) &&
				yield(m.rotateRight().MoveUp().MoveUp()) &&
				yield(m.rotateRight().MoveRight().MoveUp().MoveUp())
		}
	case Angle180:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveRight()) &&
				yield(m.rotateRight().MoveRight().MoveUp()) &&
				yield(m.rotateRight().MoveDown().MoveDown()) &&
				yield(m.rotateRight().MoveRight().MoveDown().MoveDown())
		}
	case Angle270:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveLeft()) &&
				yield(m.rotateRight().MoveLeft().MoveDown()) &&
				yield(m.rotateRight().MoveUp().MoveUp()) &&
				yield(m.rotateRight().MoveLeft().MoveUp().MoveUp())
		}

	default:
		panic("Invalid angle")
	}
}

func (m BaseMino) RotateLeftSSR() iter.Seq[AbstractMino] {
	switch m.angle {
	case Angle0:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveRight()) &&
				yield(m.rotateLeft().MoveRight().MoveUp()) &&
				yield(m.rotateLeft().MoveDown().MoveDown()) &&
				yield(m.rotateLeft().MoveRight().MoveDown().MoveDown())
		}
	case Angle90:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveRight()) &&
				yield(m.rotateLeft().MoveRight().MoveDown()) &&
				yield(m.rotateLeft().MoveUp().MoveUp()) &&
				yield(m.rotateLeft().MoveRight().MoveUp().MoveUp())
		}
	case Angle180:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveLeft()) &&
				yield(m.rotateLeft().MoveLeft().MoveUp()) &&
				yield(m.rotateLeft().MoveDown().MoveDown()) &&
				yield(m.rotateLeft().MoveLeft().MoveDown().MoveDown())
		}
	case Angle270:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveLeft()) &&
				yield(m.rotateLeft().MoveLeft().MoveDown()) &&
				yield(m.rotateLeft().MoveUp().MoveUp()) &&
				yield(m.rotateLeft().MoveLeft().MoveUp().MoveUp())
		}
	default:
		panic("Invalid angle")
	}
}

type AbstractMino interface {
	Initialize() AbstractMino
	MoveRight() AbstractMino
	MoveLeft() AbstractMino
	MoveDown() AbstractMino
	MoveUp() AbstractMino
	rotateRight() AbstractMino
	rotateLeft() AbstractMino
	RotateRightSRS() iter.Seq[AbstractMino]
	RotateLeftSSR() iter.Seq[AbstractMino]
	Shape() Shape
	Color() color.Color
	X() int
	Y() int
}

type HoldingMino struct {
	AbstractMino
	Available bool
}

type MinoI struct {
	BaseMino
}

func NewMinoI() MinoI {
	return MinoI{
		BaseMino: NewBaseMino(
			[][]int{
				{0, 0, 0, 0},
				{1, 1, 1, 1},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			CYAN,
		),
	}
}

type MinoJ struct {
	BaseMino
}

func NewMinoJ() MinoJ {
	return MinoJ{
		BaseMino: NewBaseMino(
			[][]int{
				{1, 0, 0},
				{1, 1, 1},
				{0, 0, 0},
			},
			BLUE,
		),
	}
}

type MinoL struct {
	BaseMino
}

func NewMinoL() MinoL {
	return MinoL{
		BaseMino: NewBaseMino(
			[][]int{
				{0, 0, 1},
				{1, 1, 1},
				{0, 0, 0},
			},
			ORANGE,
		),
	}
}

type MinoO struct {
	BaseMino
}

func NewMinoO() MinoO {
	return MinoO{
		BaseMino: NewBaseMino(
			[][]int{
				{1, 1},
				{1, 1},
			},
			YELLOW,
		),
	}
}

type MinoS struct {
	BaseMino
}

func NewMinoS() MinoS {
	return MinoS{
		BaseMino: NewBaseMino(
			[][]int{
				{0, 1, 1},
				{1, 1, 0},
				{0, 0, 0},
			},
			GREEN,
		),
	}
}

type MinoT struct {
	BaseMino
}

func NewMinoT() MinoT {
	return MinoT{
		BaseMino: NewBaseMino(
			[][]int{
				{0, 1, 0},
				{1, 1, 1},
				{0, 0, 0},
			},
			PURPLE,
		),
	}
}

type MinoZ struct {
	BaseMino
}

func NewMinoZ() MinoZ {
	return MinoZ{
		BaseMino: NewBaseMino(
			[][]int{
				{1, 1, 0},
				{0, 1, 1},
				{0, 0, 0},
			},
			RED,
		),
	}
}

func (m MinoI) RotateRightSRS() iter.Seq[AbstractMino] {
	switch m.angle {
	case Angle0:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveLeft().MoveLeft()) &&
				yield(m.rotateRight().MoveRight()) &&
				yield(m.rotateRight().MoveLeft().MoveLeft().MoveDown()) &&
				yield(m.rotateRight().MoveRight().MoveUp().MoveUp())

		}
	case Angle90:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveLeft()) &&
				yield(m.rotateRight().MoveRight().MoveRight()) &&
				yield(m.rotateRight().MoveLeft().MoveUp().MoveUp()) &&
				yield(m.rotateRight().MoveRight().MoveRight().MoveDown())

		}
	case Angle180:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveRight().MoveRight()) &&
				yield(m.rotateRight().MoveLeft()) &&
				yield(m.rotateRight().MoveRight().MoveRight().MoveUp()) &&
				yield(m.rotateRight().MoveLeft().MoveDown().MoveDown())
		}
	case Angle270:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateRight()) &&
				yield(m.rotateRight().MoveLeft().MoveLeft()) &&
				yield(m.rotateRight().MoveRight()) &&
				yield(m.rotateRight().MoveLeft().MoveDown().MoveDown()) &&
				yield(m.rotateRight().MoveRight().MoveUp().MoveUp())
		}
	default:
		panic("Invalid angle")
	}
}

func (m MinoI) RotateLeftSSR() iter.Seq[AbstractMino] {
	switch m.angle {
	case Angle0:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveLeft()) &&
				yield(m.rotateLeft().MoveRight().MoveRight()) &&
				yield(m.rotateLeft().MoveLeft().MoveUp().MoveUp()) &&
				yield(m.rotateLeft().MoveRight().MoveRight().MoveDown())
		}
	case Angle90:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveRight().MoveRight()) &&
				yield(m.rotateLeft().MoveLeft()) &&
				yield(m.rotateLeft().MoveRight().MoveRight().MoveUp()) &&
				yield(m.rotateLeft().MoveLeft().MoveDown().MoveDown())
		}
	case Angle180:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveRight().MoveRight()) &&
				yield(m.rotateLeft().MoveLeft()) &&
				yield(m.rotateLeft().MoveRight().MoveRight().MoveUp()) &&
				yield(m.rotateLeft().MoveLeft().MoveDown().MoveDown())
		}
	case Angle270:
		return func(yield func(AbstractMino) bool) {
			_ = yield(m.rotateLeft()) &&
				yield(m.rotateLeft().MoveRight()) &&
				yield(m.rotateLeft().MoveLeft().MoveLeft()) &&
				yield(m.rotateLeft().MoveLeft().MoveLeft().MoveDown()) &&
				yield(m.rotateLeft().MoveRight().MoveUp().MoveUp())
		}
	default:
		panic("Invalid angle")
	}
}

func (m MinoO) RotateRightSRS() iter.Seq[AbstractMino] {
	return func(yield func(AbstractMino) bool) {
		yield(m)
	}
}

func (m MinoO) RotateLeftSSR() iter.Seq[AbstractMino] {
	return func(yield func(AbstractMino) bool) {
		yield(m)
	}
}

var Minos = []AbstractMino{
	NewMinoI(),
	NewMinoJ(),
	NewMinoL(),
	NewMinoO(),
	NewMinoS(),
	NewMinoT(),
	NewMinoZ(),
}

type MinoBag struct {
	queue []AbstractMino
}

func (b *MinoBag) fill() {
	bag := make([]AbstractMino, len(Minos))
	copy(bag, Minos)
	for i := range len(bag) {
		j := rand.Intn(i + 1)
		bag[i], bag[j] = bag[j], bag[i]
	}
	b.queue = append(b.queue, bag...)
}

func (b *MinoBag) Sniff(n int) []AbstractMino {
	if n <= 0 || n > 7 {
		panic("n must be between 1 and 7")
	}
	if len(b.queue) < n {
		b.fill()
	}
	preview := make([]AbstractMino, n)
	copy(preview, b.queue[:n])
	return preview
}

func (b *MinoBag) Next() AbstractMino {
	if len(b.queue) == 0 {
		b.fill()
	}
	mino := b.queue[0].Initialize()
	b.queue = b.queue[1:]
	return mino
}

// A fragment is a small piece of a mino that is animated when it is cleared
type Fragment struct {
	Frame            int
	_Color           color.Color
	InitialX         int
	InitialY         int
	AccelerationX    float32
	AccelerationY    float32 // Gravity
	InitialVelocityX float32
	InitialVelocityY float32
}

func NewFragment(color color.Color, x, y int) Fragment {
	return Fragment{
		Frame:            30,
		_Color:           color,
		InitialX:         x,
		InitialY:         y,
		AccelerationX:    0,
		AccelerationY:    1,
		InitialVelocityX: rand.Float32()*6 - 3,
		InitialVelocityY: -3,
	}
}

func (f *Fragment) Position() (x, y float32) {
	x = calc(f.InitialVelocityX, f.AccelerationX, float32(30-f.Frame)) + float32(f.InitialX*CELL_SIZE+CELL_SIZE/2)
	y = calc(f.InitialVelocityY, f.AccelerationY, float32(30-f.Frame)) + float32(f.InitialY*CELL_SIZE+CELL_SIZE/2)
	return
}

func (f *Fragment) Color() color.Color {
	r, g, b, _ := f._Color.RGBA()
	return color.RGBA{
		uint8(r / 256),
		uint8(g / 256),
		uint8(b / 256),
		uint8(f.Frame / 30 * 255),
	}
}

func calc(v, a, t float32) float32 {
	return v*t + 0.5*a*t*t
}

func Rotate(shape Shape) Shape {
	n := len(shape)
	rotated := make([][]int, n)
	for i := range n {
		rotated[i] = make([]int, n)
		for j := range n {
			rotated[i][j] = shape[n-j-1][i]
		}
	}
	return rotated
}

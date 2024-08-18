package game

import (
	"image/color"
	"iter"
	"math/rand"
)

var (
	GHOST_COLOR = color.RGBA{30, 30, 30, 127}
)

const (
	DEFAULT_BACKLASH_FRAME           = 30
	DEFAULT_EXTENDED_PLACEMENT_COUNT = 15
)

type Shape [][]int

// Note: the Mino is fully fixed if IsGrounded is true and BacklashFrame is 0 or ExtendedPlacementCounter is 0
type Mino struct {
	Name          string
	Color         color.Color
	Y             int
	X             int
	Angle         int
	FrameCount    int
	LockDown      LockDown
	IsGrounded    bool
	BacklashFrame int // Allow a little time for movement / rotation after grounding
}

type HoldingMino struct {
	Mino
	Available bool
}

func (m *Mino) Shape() Shape {
	switch m.Name {
	case "T":
		return Ts[m.Angle]
	case "O":
		return Os[m.Angle]
	case "L":
		return Ls[m.Angle]
	case "J":
		return Js[m.Angle]
	case "S":
		return Ss[m.Angle]
	case "Z":
		return Zs[m.Angle]
	case "I":
		return Is[m.Angle]
	default:
		panic("Invalid mino name")
	}
}

// Rotate the mino 90 degrees clockwise
func (m Mino) rotateRight() Mino {
	m.Angle = (m.Angle + 1) % 4
	return m
}

// Rotate the mino 90 degrees counterclockwise
func (m Mino) rotateLeft() Mino {
	m.Angle = (m.Angle + 3) % 4
	return m
}

// Yield all possible rotations of the mino in the Super Rotation System
func (m Mino) RotateRightSRS() iter.Seq[Mino] {
	switch m.Name {
	case "T", "L", "J", "S", "Z":
		switch m.Angle {
		case 0:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveRight()) &&
					yield(m.rotateRight().MoveRight().MoveUp()) &&
					yield(m.rotateRight().MoveDown().MoveDown()) &&
					yield(m.rotateRight().MoveRight().MoveDown().MoveDown())
			}
		case 1:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveRight()) &&
					yield(m.rotateRight().MoveRight().MoveDown()) &&
					yield(m.rotateRight().MoveUp().MoveUp()) &&
					yield(m.rotateRight().MoveRight().MoveUp().MoveUp())
			}
		case 2:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveRight()) &&
					yield(m.rotateRight().MoveRight().MoveUp()) &&
					yield(m.rotateRight().MoveDown().MoveDown()) &&
					yield(m.rotateRight().MoveRight().MoveDown().MoveDown())
			}
		case 3:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveLeft()) &&
					yield(m.rotateRight().MoveLeft().MoveDown()) &&
					yield(m.rotateRight().MoveUp().MoveUp()) &&
					yield(m.rotateRight().MoveLeft().MoveUp().MoveUp())
			}

		default:
			panic("Invalid angle")
		}
	case "I":
		switch m.Angle {
		case 0:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveLeft().MoveLeft()) &&
					yield(m.rotateRight().MoveRight()) &&
					yield(m.rotateRight().MoveLeft().MoveLeft().MoveDown()) &&
					yield(m.rotateRight().MoveRight().MoveUp().MoveUp())

			}
		case 1:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveLeft()) &&
					yield(m.rotateRight().MoveRight().MoveRight()) &&
					yield(m.rotateRight().MoveLeft().MoveUp().MoveUp()) &&
					yield(m.rotateRight().MoveRight().MoveRight().MoveDown())

			}
		case 2:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveRight().MoveRight()) &&
					yield(m.rotateRight().MoveLeft()) &&
					yield(m.rotateRight().MoveRight().MoveRight().MoveUp()) &&
					yield(m.rotateRight().MoveLeft().MoveDown().MoveDown())
			}
		case 3:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateRight()) &&
					yield(m.rotateRight().MoveLeft().MoveLeft()) &&
					yield(m.rotateRight().MoveRight()) &&
					yield(m.rotateRight().MoveLeft().MoveDown().MoveDown()) &&
					yield(m.rotateRight().MoveRight().MoveUp().MoveUp())
			}
		default:
			panic("Invalid angle")
		}
	case "O":
		return func(yield func(Mino) bool) {
			yield(m)
		}
	default:
		panic("Invalid mino name")
	}
}

// Yield all possible rotations of the mino in the Super Rotation System
func (m Mino) RotateLeftSSR() iter.Seq[Mino] {
	switch m.Name {
	case "T", "L", "J", "S", "Z":
		switch m.Angle {
		case 0:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveRight()) &&
					yield(m.rotateLeft().MoveRight().MoveUp()) &&
					yield(m.rotateLeft().MoveDown().MoveDown()) &&
					yield(m.rotateLeft().MoveRight().MoveDown().MoveDown())
			}
		case 1:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveRight()) &&
					yield(m.rotateLeft().MoveRight().MoveDown()) &&
					yield(m.rotateLeft().MoveUp().MoveUp()) &&
					yield(m.rotateLeft().MoveRight().MoveUp().MoveUp())
			}
		case 2:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveLeft()) &&
					yield(m.rotateLeft().MoveLeft().MoveUp()) &&
					yield(m.rotateLeft().MoveDown().MoveDown()) &&
					yield(m.rotateLeft().MoveLeft().MoveDown().MoveDown())
			}
		case 3:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveLeft()) &&
					yield(m.rotateLeft().MoveLeft().MoveDown()) &&
					yield(m.rotateLeft().MoveUp().MoveUp()) &&
					yield(m.rotateLeft().MoveLeft().MoveUp().MoveUp())
			}
		default:
			panic("Invalid angle")
		}
	case "I":
		switch m.Angle {
		case 0:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveLeft()) &&
					yield(m.rotateLeft().MoveRight().MoveRight()) &&
					yield(m.rotateLeft().MoveLeft().MoveUp().MoveUp()) &&
					yield(m.rotateLeft().MoveRight().MoveRight().MoveDown())
			}
		case 1:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveRight().MoveRight()) &&
					yield(m.rotateLeft().MoveLeft()) &&
					yield(m.rotateLeft().MoveRight().MoveRight().MoveUp()) &&
					yield(m.rotateLeft().MoveLeft().MoveDown().MoveDown())
			}
		case 2:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveRight().MoveRight()) &&
					yield(m.rotateLeft().MoveLeft()) &&
					yield(m.rotateLeft().MoveRight().MoveRight().MoveUp()) &&
					yield(m.rotateLeft().MoveLeft().MoveDown().MoveDown())
			}
		case 3:
			return func(yield func(Mino) bool) {
				_ = yield(m.rotateLeft()) &&
					yield(m.rotateLeft().MoveRight()) &&
					yield(m.rotateLeft().MoveLeft().MoveLeft()) &&
					yield(m.rotateLeft().MoveLeft().MoveLeft().MoveDown()) &&
					yield(m.rotateLeft().MoveRight().MoveUp().MoveUp())
			}
		default:
			panic("Invalid angle")
		}
	case "O":
		return func(yield func(Mino) bool) {
			yield(m)
		}
	default:
		panic("Invalid mino name")
	}
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

func (m Mino) MoveUp() Mino {
	m.Y--
	return m
}

var (
	Ts = [4]Shape{
		[][]int{
			{0, 1, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		[][]int{
			{0, 1, 0},
			{0, 1, 1},
			{0, 1, 0},
		},
		[][]int{
			{0, 0, 0},
			{1, 1, 1},
			{0, 1, 0},
		},
		[][]int{
			{0, 1, 0},
			{1, 1, 0},
			{0, 1, 0},
		},
	}
	Os = [4]Shape{
		[][]int{
			{1, 1},
			{1, 1},
		},
		[][]int{
			{1, 1},
			{1, 1},
		},
		[][]int{
			{1, 1},
			{1, 1},
		},
		[][]int{
			{1, 1},
			{1, 1},
		},
	}
	Ls = [4]Shape{
		[][]int{
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 0},
		},
		[][]int{
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 1},
		},
		[][]int{
			{0, 0, 0},
			{1, 1, 1},
			{1, 0, 0},
		},
		[][]int{
			{1, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		},
	}
	Js = [4]Shape{
		[][]int{
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		[][]int{
			{0, 1, 1},
			{0, 1, 0},
			{0, 1, 0},
		},
		[][]int{
			{0, 0, 0},
			{1, 1, 1},
			{0, 0, 1},
		},
		[][]int{
			{0, 1, 0},
			{0, 1, 0},
			{1, 1, 0},
		},
	}
	Ss = [4]Shape{
		[][]int{
			{0, 1, 1},
			{1, 1, 0},
			{0, 0, 0},
		},
		[][]int{
			{0, 1, 0},
			{0, 1, 1},
			{0, 0, 1},
		},
		[][]int{
			{0, 0, 0},
			{0, 1, 1},
			{1, 1, 0},
		},
		[][]int{
			{1, 0, 0},
			{1, 1, 0},
			{0, 1, 0},
		},
	}
	Zs = [4]Shape{
		[][]int{
			{1, 1, 0},
			{0, 1, 1},
			{0, 0, 0},
		},
		[][]int{
			{0, 1, 0},
			{1, 1, 0},
			{1, 0, 0},
		},
		[][]int{
			{0, 0, 0},
			{1, 1, 0},
			{0, 1, 1},
		},
		[][]int{
			{0, 1, 0},
			{1, 1, 0},
			{1, 0, 0},
		},
	}
	Is = [4]Shape{
		[][]int{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		[][]int{
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 1, 0},
		},
		[][]int{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
		},
		[][]int{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
		},
	}
)

var (
	T = Mino{
		Name:          "T",
		Color:         color.RGBA{106, 50, 165, 255}, // Purple
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	O = Mino{
		Name:          "O",
		Color:         color.RGBA{255, 213, 0, 255}, // Yellow
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	L = Mino{
		Name:          "L",
		Color:         color.RGBA{255, 121, 28, 255}, // Orange
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	J = Mino{
		Name:          "J",
		Color:         color.RGBA{6, 119, 186, 255}, // Blue
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	S = Mino{
		Name:          "S",
		Color:         color.RGBA{114, 203, 59, 255}, // Green
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	Z = Mino{
		Name:          "Z",
		Color:         color.RGBA{212, 42, 52, 255}, // Red
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	I = Mino{
		Name:          "I",
		Color:         color.RGBA{31, 195, 205, 255}, // Cyan
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

// An implementation of the extended placement system
//   - After a mino is grounded, `isGrounded` flag is set to true then the `timer` and `counter` are started
//   - `timer` is incremented every frame until it reaches `DEFAULT_BACKLASH_FRAME`
//   - If the mino is moved or rotated, `timer` is reset, but `counter` is incremented
//   - The mino is fixed if `timer` reaches `DEFAULT_BACKLASH_FRAME` or
//     `counter` reaches `DEFAULT_EXTENDED_PLACEMENT_COUNT` even though `timer` is less than `DEFAULT_BACKLASH_FRAME`
type LockDown struct {
	isGrounded bool
	timer      int
	counter    int
}

// Return true if the mino should not be moved or rotated anymore
func (l *LockDown) IsFixed() bool {
	return (l.isGrounded && l.timer >= DEFAULT_BACKLASH_FRAME) || l.counter >= DEFAULT_EXTENDED_PLACEMENT_COUNT
}

func (l *LockDown) Activate() {
	l.isGrounded = true
}

func (l *LockDown) UpdateCounter() {
	if l.isGrounded {
		l.counter++
		l.timer = 0
	}
}

func (l *LockDown) UpdateTimer() {
	if l.isGrounded {
		l.timer++
	}
}

func (l *LockDown) Reset() {
	l.isGrounded = false
	l.timer = 0
	l.counter = 0
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

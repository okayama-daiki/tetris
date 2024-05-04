package game

import (
	"image/color"
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

func (s Shape) deepCopy() Shape {
	shape := make(Shape, len(s))
	for i := range s {
		shape[i] = make([]int, len(s[i]))
		copy(shape[i], s[i])
	}
	return shape
}

// Note: the Mino is fully fixed if IsGrounded is true and BacklashFrame is 0 or ExtendedPlacementCounter is 0
type Mino struct {
	Name          string
	Shape         Shape
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

// Rotate the mino 90 degrees clockwise
func (m Mino) rotateRight() Mino {
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
	m.Angle = (m.Angle + 1) % 4
	return m
}

// Rotate the mino 90 degrees counterclockwise
func (m Mino) rotateLeft() Mino {
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
	m.Angle = (m.Angle + 3) % 4
	return m
}

// Return all possible rotations of the mino in the Super Rotation System
//
// TODO: This function must be rewritten as `func(yield func(Mino) bool)` ...
// and used as Range Over Function after Go 1.23 is released
func (m Mino) RotateRightSRS() []Mino {
	switch m.Name {
	case "T", "L", "J", "S", "Z":
		switch m.Angle {
		case 0:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveLeft(),
				m.rotateRight().MoveLeft().MoveUp(),
				m.rotateRight().MoveDown().MoveDown(),
				m.rotateRight().MoveLeft().MoveDown().MoveDown(),
			}
		case 1:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveRight(),
				m.rotateRight().MoveRight().MoveDown(),
				m.rotateRight().MoveUp().MoveUp(),
				m.rotateRight().MoveRight().MoveUp().MoveUp(),
			}
		case 2:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveRight(),
				m.rotateRight().MoveRight().MoveUp(),
				m.rotateRight().MoveDown().MoveDown(),
				m.rotateRight().MoveRight().MoveDown().MoveDown(),
			}
		case 3:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveLeft(),
				m.rotateRight().MoveLeft().MoveDown(),
				m.rotateRight().MoveUp().MoveUp(),
				m.rotateRight().MoveLeft().MoveUp().MoveUp(),
			}

		default:
			panic("Invalid angle")
		}
	case "I":
		switch m.Angle {
		case 0:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveLeft().MoveLeft(),
				m.rotateRight().MoveRight(),
				m.rotateRight().MoveLeft().MoveLeft().MoveDown(),
				m.rotateRight().MoveRight().MoveUp().MoveUp(),
			}
		case 1:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveLeft(),
				m.rotateRight().MoveRight().MoveRight(),
				m.rotateRight().MoveLeft().MoveUp().MoveUp(),
				m.rotateRight().MoveRight().MoveRight().MoveDown(),
			}
		case 2:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveRight().MoveRight(),
				m.rotateRight().MoveLeft(),
				m.rotateRight().MoveRight().MoveRight().MoveUp(),
				m.rotateRight().MoveLeft().MoveDown().MoveDown(),
			}
		case 3:
			return []Mino{
				m.rotateRight(),
				m.rotateRight().MoveLeft().MoveLeft(),
				m.rotateRight().MoveRight(),
				m.rotateRight().MoveLeft().MoveDown().MoveDown(),
				m.rotateRight().MoveRight().MoveUp().MoveUp(),
			}
		default:
			panic("Invalid angle")
		}
	case "O":
		return []Mino{m}
	default:
		panic("Invalid mino name")
	}
}

// Return all possible rotations of the mino in the Super Rotation System
func (m Mino) RotateLeftSSR() []Mino {
	switch m.Name {
	case "T", "L", "J", "S", "Z":
		switch m.Angle {
		case 0:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveRight(),
				m.rotateLeft().MoveRight().MoveUp(),
				m.rotateLeft().MoveDown().MoveDown(),
				m.rotateLeft().MoveRight().MoveDown().MoveDown(),
			}
		case 1:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveRight(),
				m.rotateLeft().MoveRight().MoveDown(),
				m.rotateLeft().MoveUp().MoveUp(),
				m.rotateLeft().MoveRight().MoveUp().MoveUp(),
			}
		case 2:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveLeft(),
				m.rotateLeft().MoveLeft().MoveUp(),
				m.rotateLeft().MoveDown().MoveDown(),
				m.rotateLeft().MoveLeft().MoveDown().MoveDown(),
			}
		case 3:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveLeft(),
				m.rotateLeft().MoveLeft().MoveDown(),
				m.rotateLeft().MoveUp().MoveUp(),
				m.rotateLeft().MoveLeft().MoveUp().MoveUp(),
			}
		default:
			panic("Invalid angle")
		}
	case "I":
		switch m.Angle {
		case 0:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveLeft(),
				m.rotateLeft().MoveRight().MoveRight(),
				m.rotateLeft().MoveLeft().MoveUp().MoveUp(),
				m.rotateLeft().MoveRight().MoveRight().MoveDown(),
			}
		case 1:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveRight().MoveRight(),
				m.rotateLeft().MoveLeft(),
				m.rotateLeft().MoveRight().MoveRight().MoveUp(),
				m.rotateLeft().MoveLeft().MoveDown().MoveDown(),
			}
		case 2:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveRight().MoveRight(),
				m.rotateLeft().MoveLeft(),
				m.rotateLeft().MoveRight().MoveRight().MoveUp(),
				m.rotateLeft().MoveLeft().MoveDown().MoveDown(),
			}
		case 3:
			return []Mino{
				m.rotateLeft(),
				m.rotateLeft().MoveRight(),
				m.rotateLeft().MoveLeft().MoveLeft(),
				m.rotateLeft().MoveLeft().MoveLeft().MoveDown(),
				m.rotateLeft().MoveRight().MoveUp().MoveUp(),
			}
		default:
			panic("Invalid angle")
		}
	case "O":
		return []Mino{m}
	default:
		panic("Invalid mino name")
	}
}

func (m Mino) MoveRight() Mino {
	m.X++
	m.Shape = m.Shape.deepCopy()
	return m
}

func (m Mino) MoveLeft() Mino {
	m.X--
	m.Shape = m.Shape.deepCopy()
	return m
}

func (m Mino) MoveDown() Mino {
	m.Y++
	m.Shape = m.Shape.deepCopy()
	return m
}

func (m Mino) MoveUp() Mino {
	m.Y--
	m.Shape = m.Shape.deepCopy()
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
		Color:         color.RGBA{106, 50, 165, 255}, // Purple
		Angle:         0,
		BacklashFrame: DEFAULT_BACKLASH_FRAME,
	}

	O = Mino{
		Shape: [][]int{
			{1, 1},
			{1, 1},
		},
		Name:          "O",
		Color:         color.RGBA{255, 213, 0, 255}, // Yellow
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
		Color:         color.RGBA{255, 121, 28, 255}, // Orange
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
		Color:         color.RGBA{6, 119, 186, 255}, // Blue
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
		Color:         color.RGBA{114, 203, 59, 255}, // Green
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
		Color:         color.RGBA{212, 42, 52, 255}, // Red
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

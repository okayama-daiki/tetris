package game

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

func NewLockDown() *LockDown {
	return &LockDown{
		isGrounded: false,
		timer:      0,
		counter:    0,
	}
}

func (l *LockDown) Ground() {
	l.isGrounded = true
}

func (l *LockDown) UnGround() {
	l.isGrounded = false
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

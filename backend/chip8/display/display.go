package display

import (
	"sync"
)

const (
	height = 32
	width  = 64
)

type Display struct {
	clipping bool

	mu     sync.Mutex
	screen [32]uint64
}

func New(clipping bool) *Display {
	return &Display{
		clipping: clipping,
	}
}

func (d *Display) Clear() {
	d.screen = [32]uint64{}
}

func (d *Display) State() [32]uint64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.screen
}

func (d *Display) Draw(sprite byte, x, y byte) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	x %= width
	y %= height

	write := uint64(sprite) << (width - 8)
	write >>= x

	prevOnes := countOnes(d.screen[y])

	d.screen[y] ^= write

	return countOnes(d.screen[y]) < prevOnes
}

// There's probably a more efficient way to do this
func countOnes(n uint64) int {
	var c int

	for n > 0 {
		c++
		n &= n - 1
	}

	return c
}

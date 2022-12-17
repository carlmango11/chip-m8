package display

import (
	"bytes"
	"fmt"
	"sync"
)

type Display struct {
	mu     sync.Mutex
	screen [32]uint64
}

func New() *Display {
	return &Display{}
}

func (d *Display) Clear() {
	d.screen = [32]uint64{}
}

func (d *Display) State() [32]uint64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.screen
}

func (d *Display) Draw(sprite byte, x, y byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	write := uint64(sprite) << 56

	for i := byte(0); i < x; i++ {
		wrap := write&0x0001 == 1

		write >>= 1

		if wrap {
			write |= 0x8000000000000000
		}
	}

	if y >= 32 {
		y -= 32
	}

	d.screen[y] ^= write
}

func (d *Display) Print() {
	d.mu.Lock()
	defer d.mu.Unlock()

	var buf bytes.Buffer

	buf.WriteString("S===============================================================")

	for i := 0; i < len(d.screen); i++ {
		buf.WriteString(fmt.Sprintf("\n%064b", d.screen[i]))
	}

	buf.WriteString("\nE===============================================================")

	fmt.Printf("\n%v", buf.String())
}

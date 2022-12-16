package display

import "fmt"

type Display struct {
	screen [32]uint64
}

func New() *Display {
	return &Display{}
}

func (d *Display) State() [32]uint64 {
	return d.screen
}

func (d *Display) Draw(sprite byte, x, y byte) {
	write := uint64(sprite) << 56

	for i := byte(0); i < y; i++ {
		wrap := write&0x0001 == 1

		write >>= 1

		if wrap {
			write |= 0x8000000000000000
		}
	}

	if x >= 32 {
		x -= 32
	}

	fmt.Printf("\nWARNING: collision not impl")

	d.screen[x] ^= write
}

func (d *Display) Print() {
	for i := len(d.screen); i >= 0; i-- {
		fmt.Printf("\n")

		for j := 0; j < 64; j++ {
			row := d.screen[i]

			msb := row & 0x80

			if msb == 1 {
				fmt.Printf("\u25A0")
			} else {
				fmt.Printf(" ")
			}

			row <<= 1
		}
	}
}

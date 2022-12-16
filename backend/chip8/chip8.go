package chip8

import (
	"github.com/carlmango11/chip-m8/backend/chip8/cpu"
	"github.com/carlmango11/chip-m8/backend/chip8/display"
	"github.com/carlmango11/chip-m8/backend/chip8/keyboard"
	"github.com/carlmango11/chip-m8/backend/chip8/ram"
)

type Chip8 struct {
	cpu      *cpu.CPU
	ram      *ram.RAM
	display  *display.Display
	keyboard *keyboard.Keyboard
}

func New(script []byte) *Chip8 {
	d := display.New()
	r := ram.New(script)
	k := keyboard.New()

	return &Chip8{
		ram:      r,
		display:  d,
		keyboard: k,
		cpu:      cpu.New(r, d, k),
	}
}

func (c *Chip8) PressKey(n byte) {
	c.keyboard.Press(n)
}

func (c *Chip8) UnpressKey(n byte) {
	c.keyboard.Unpress(n)
}

func (c *Chip8) Display() [32]uint64 {
	return c.display.State()
}

package chip8

import (
	"github.com/carlmango11/chip-m8/backend/chip8/cpu"
	"github.com/carlmango11/chip-m8/backend/chip8/display"
	"github.com/carlmango11/chip-m8/backend/chip8/keyboard"
	"github.com/carlmango11/chip-m8/backend/chip8/ram"
	"time"
)

const clockSpeed = 1e6 // 1Ghz

type Config struct {
	CPU             cpu.Config
	DisplayClipping bool
}

var Chip8Cfg = Config{
	CPU:             cpu.ConfigChip8,
	DisplayClipping: false,
}

type Chip8 struct {
	cpu      *cpu.CPU
	ram      *ram.RAM
	display  *display.Display
	keyboard *keyboard.Keyboard
	clock    *time.Ticker
}

func New(script []byte, cfg Config) *Chip8 {
	d := display.New(cfg.DisplayClipping)
	r := ram.New(script)
	k := keyboard.New()

	return &Chip8{
		ram:      r,
		display:  d,
		keyboard: k,
		cpu:      cpu.New(cfg.CPU, r, d, k),
		clock:    time.NewTicker(time.Second / clockSpeed),
	}
}

func (c *Chip8) Start() {
	for range c.clock.C {
		c.cpu.Tick()
	}
}

func (c *Chip8) Stop() {
	c.clock.Stop()
	c.display.Clear()
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

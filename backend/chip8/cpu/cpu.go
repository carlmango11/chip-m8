package cpu

import (
	"github.com/carlmango11/chip-m8/backend/chip8/display"
	"github.com/carlmango11/chip-m8/backend/chip8/keyboard"
	"github.com/carlmango11/chip-m8/backend/chip8/ram"
	"math/rand"
	"time"
)

const timerInterval = time.Second / 60

type Config struct {
	VFReset         bool
	MemoryIncrement bool
	ShiftY          bool
}

var ConfigChip8 = Config{
	VFReset:         true,
	MemoryIncrement: true,
	ShiftY:          true,
}

type CPU struct {
	cfg Config

	ram      *ram.RAM
	display  *display.Display
	keyboard *keyboard.Keyboard

	lastTimerDecrement time.Time

	pc ram.Address
	sp byte

	v  [16]byte
	i  ram.Address
	dt byte
	st byte

	stack [16]ram.Address
}

func New(cfg Config, r *ram.RAM, display *display.Display, kb *keyboard.Keyboard) *CPU {
	return &CPU{
		cfg:      cfg,
		ram:      r,
		display:  display,
		keyboard: kb,
		pc:       ram.UserMemoryStart,
	}
}

func (c *CPU) Tick() {
	hi := c.ram.Read(c.nextPC())
	lo := c.ram.Read(c.nextPC())

	top := 0x00FF | uint16(hi)<<8
	bottom := 0xFF00 | uint16(lo)

	opCode := top & bottom

	c.executeOpCode(opCode)
	c.updateTimers()
}

func (c *CPU) updateTimers() {
	// startup
	if c.lastTimerDecrement.IsZero() {
		c.lastTimerDecrement = time.Now()
		return
	}

	now := time.Now()
	since := now.Sub(c.lastTimerDecrement)

	// note this technique probably only works because we're sure to tick faster than the timer interval
	if since > timerInterval {
		if c.st > 0 {
			c.st--
		}

		if c.dt > 0 {
			c.dt--
		}

		c.lastTimerDecrement = now
	}
}

func (c *CPU) nextPC() ram.Address {
	n := c.pc
	c.pc += 1

	return n
}

func (c *CPU) executeOpCode(opCode uint16) {
	top := opCode & 0xF000

	switch top {
	case 0x0000:
		switch opCode {
		case 0x00E0:
			c.CLS()
		case 0x00EE:
			c.RET()
		}
	case 0x1000:
		c.JP_NNN(opCode)
	case 0x2000:
		c.CALL_NNN(opCode)
	case 0x3000:
		c.SE_KK(opCode)
	case 0x4000:
		c.SNE_KK(opCode)
	case 0x5000:
		c.SE_XY(opCode)
	case 0x6000:
		c.LD_KK(opCode)
	case 0x7000:
		c.ADD_KK(opCode)
	case 0x8000:
		bottom := opCode & 0x000F

		switch bottom {
		case 0x0000:
			c.LD_XY(opCode)
		case 0x0001:
			c.OR_XY(opCode)
		case 0x0002:
			c.AND_XY(opCode)
		case 0x0003:
			c.XOR_XY(opCode)
		case 0x0004:
			c.ADD_XY(opCode)
		case 0x0005:
			c.SUB_XY(opCode)
		case 0x0006:
			c.SHR_XY(opCode)
		case 0x0007:
			c.SUBN_XY(opCode)
		case 0x000E:
			c.SHL_XY(opCode)
		}
	case 0x9000:
		c.SNE_XY(opCode)
	case 0xA000:
		c.LD_I(opCode)
	case 0xB000:
		c.JP_V0(opCode)
	case 0xC000:
		c.RND(opCode)
	case 0xD000:
		c.DXYN(opCode)
	case 0xE000:
		switch opCode & 0x00FF {
		case 0x009E:
			c.EX9E(opCode)
		case 0x00A1:
			c.EXA1(opCode)
		}
	case 0xF000:
		switch opCode & 0x00FF {
		case 0x0007:
			c.FX07(opCode)
		case 0x000A:
			c.FX0A(opCode)
		case 0x0015:
			c.FX15(opCode)
		case 0x0018:
			c.FX18(opCode)
		case 0x001E:
			c.FX1E(opCode)
		case 0x0029:
			c.FX29(opCode)
		case 0x0033:
			c.FX33(opCode)
		case 0x0055:
			c.FX55(opCode)
		case 0x0065:
			c.FX65(opCode)
		}
	}
}

func (c *CPU) FX65(opCode uint16) {
	regX, _ := extractXY(opCode)

	for i := byte(0); i <= regX; i++ {
		addr := c.i
		if !c.cfg.MemoryIncrement {
			addr += ram.Address(i)
		}

		c.v[i] = c.ram.Read(addr)

		if c.cfg.MemoryIncrement {
			c.i++
		}
	}
}

func (c *CPU) FX55(opCode uint16) {
	regX, _ := extractXY(opCode)

	for i := byte(0); i <= regX; i++ {
		addr := c.i
		if !c.cfg.MemoryIncrement {
			addr += ram.Address(i)
		}

		c.ram.Write(addr, c.v[i])

		if c.cfg.MemoryIncrement {
			c.i++
		}
	}
}

func (c *CPU) FX33(opCode uint16) {
	regX, _ := extractXY(opCode)

	val := c.v[regX]

	c.ram.Write(c.i, val/100)
	c.ram.Write(c.i+1, val%100/10)
	c.ram.Write(c.i+2, val%10)
}

func (c *CPU) FX29(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.i = uint16(c.v[regX]) * 5
}

func (c *CPU) FX1E(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.i += uint16(c.v[regX])
}

func (c *CPU) FX18(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.st = c.v[regX]
}

func (c *CPU) FX15(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.dt = c.v[regX]
}

func (c *CPU) FX0A(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.v[regX] = c.keyboard.Await()
}

func (c *CPU) FX07(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.v[regX] = c.dt
}

func (c *CPU) EXA1(opCode uint16) {
	regX, _ := extractXY(opCode)

	if !c.keyboard.IsPressed(c.v[regX]) {
		c.pc += 2
	}
}

func (c *CPU) EX9E(opCode uint16) {
	regX, _ := extractXY(opCode)

	if c.keyboard.IsPressed(c.v[regX]) {
		c.pc += 2
	}
}

func (c *CPU) DXYN(opCode uint16) {
	regX, regY := extractXY(opCode)
	n := opCode & 0x000F

	var collision bool

	for i := uint16(0); i < n; i++ {
		x := c.v[regX]
		y := c.v[regY]

		sprite := c.ram.Read(c.i + i)

		collision = collision || c.display.Draw(sprite, x, y+byte(i))
	}

	if collision {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
}

func (c *CPU) RND(opCode uint16) {
	regX, _ := extractXY(opCode)

	c.v[regX] = byte(rand.Intn(256)) & getKK(opCode)
}

func (c *CPU) LD_I(opCode uint16) {
	c.i = asAddr(opCode)
}

func (c *CPU) JP_V0(opCode uint16) {
	c.pc = asAddr(opCode) + uint16(c.v[0])
}

func (c *CPU) SNE_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	if c.v[regY] != c.v[regX] {
		c.pc += 2
	}
}

func (c *CPU) SHL_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	if c.cfg.ShiftY {
		c.v[regX] = c.v[regY]
	}

	carry := c.v[regX]&0x80 == 0x80

	c.v[regX] <<= 1

	if carry {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
}

func (c *CPU) SUBN_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	carry := c.v[regY] > c.v[regX]

	c.v[regX] = c.v[regY] - c.v[regX]

	if carry {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
}

func (c *CPU) SHR_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	if c.cfg.ShiftY {
		c.v[regX] = c.v[regY]
	}

	carry := c.v[regX]&0x01 == 0x01

	c.v[regX] >>= 1

	if carry {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
}

func (c *CPU) SUB_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	carry := c.v[regX] > c.v[regY]

	c.v[regX] = c.v[regX] - c.v[regY]

	if carry {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
}

func (c *CPU) ADD_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	res := uint16(c.v[regX]) + uint16(c.v[regY])

	if res > 0xFF {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}

	c.v[regX] = byte(res)
}

func (c *CPU) XOR_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	c.v[regX] = c.v[regX] ^ c.v[regY]
	c.v[0xF] = 0
}

func (c *CPU) OR_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	c.v[regX] = c.v[regX] | c.v[regY]
	c.v[0xF] = 0
}

func (c *CPU) AND_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	c.v[regX] = c.v[regX] & c.v[regY]
	c.v[0xF] = 0
}

func (c *CPU) LD_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	c.v[regX] = c.v[regY]
}

func (c *CPU) ADD_KK(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.v[regX] += byte(opCode & 0x00FF)
}

func (c *CPU) LD_KK(opCode uint16) {
	regX, _ := extractXY(opCode)
	c.v[regX] = byte(opCode & 0x00FF)
}

func (c *CPU) SNE_KK(opCode uint16) {
	regX, _ := extractXY(opCode)
	toCmp := getKK(opCode)

	if c.v[regX] != toCmp {
		c.pc += 2
	}
}

func (c *CPU) SE_KK(opCode uint16) {
	regX, _ := extractXY(opCode)
	toCmp := getKK(opCode)

	if c.v[regX] == toCmp {
		c.pc += 2
	}
}

func (c *CPU) SE_XY(opCode uint16) {
	regX, regY := extractXY(opCode)

	if c.v[regX] == c.v[regY] {
		c.pc += 2
	}
}

func (c *CPU) CLS() {
	c.display.Clear()
}

func (c *CPU) RET() {
	c.pc = c.stack[c.sp]
	c.sp--

	if c.sp < 0 {
		panic("stack pointer is negative")
	}
}

func (c *CPU) JP_NNN(opCode uint16) {
	addr := asAddr(opCode)
	c.pc = addr
}

func (c *CPU) CALL_NNN(opCode uint16) {
	addr := asAddr(opCode)

	c.sp++
	c.stack[c.sp] = c.pc

	if c.sp > 15 {
		panic("stack overflow")
	}

	c.pc = addr
}

func getKK(opCode uint16) byte {
	return byte(opCode & 0x00FF)
}

func asAddr(full uint16) ram.Address {
	return (full << 4) >> 4
}

func extractXY(opCode uint16) (byte, byte) {
	regX := opCode & 0x0F00 >> 8
	regY := opCode & 0x00F0 >> 4

	return byte(regX), byte(regY)
}

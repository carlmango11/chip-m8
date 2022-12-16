package cpu

import (
	"github.com/carlmango11/chip-m8/backend/chip8/display"
	"github.com/carlmango11/chip-m8/backend/chip8/keyboard"
	"github.com/carlmango11/chip-m8/backend/chip8/ram"
	"math/rand"
)

type CPU struct {
	ram      *ram.RAM
	display  *display.Display
	keyboard *keyboard.Keyboard

	pc ram.Address
	sp byte

	general [16]byte
	i       ram.Address
	vf      byte
	dt      byte
	st      byte

	stack [16]ram.Address
}

func New(ram *ram.RAM, display *display.Display, kb *keyboard.Keyboard) *CPU {
	return &CPU{
		ram:      ram,
		display:  display,
		keyboard: kb,
	}
}

func (c *CPU) Tick() {
	hi := c.ram.Read(c.nextPC())
	lo := c.ram.Read(c.nextPC())

	opCode := (uint16(hi) << 8) & uint16(lo)

	c.executeOpCode(opCode)
}

func (c *CPU) nextPC() ram.Address {
	n := c.pc
	n += 1

	return n
}

func (c *CPU) executeOpCode(opCode uint16) {
	top := opCode & 0xF000

	switch top {
	case 0x0000:
		switch opCode {
		case 0x00E0:
			c.clearDisplay()
		case 0x00EE:
			c.returnSub()
		}
	case 0x1000:
		c.JP_NNN(asAddr(opCode))
	case 0x2000:
		c.CALL_NNN(asAddr(opCode))
	case 0x3000:
		c.SE_KK(opCode)
	case 0x4000:
		c.SNE_KK(opCode)
	case 0x5000:
		c.SE_XY(opCode)
	case 0x6000:
		c.LD_KK(opCode)
	case 0x7000:
		c.addToRegister(opCode)
	case 0x8000:
		bottom := opCode & 0x000F

		switch bottom {
		case 0x0000:
			c.setRegisterFromOther(opCode)
		case 0x0001:
			c.doRegisterOR(opCode)
		case 0x0002:
			c.doRegisterAND(opCode)
		case 0x0003:
			c.doRegisterXOR(opCode)
		case 0x0004:
			c.doRegisterADD(opCode)
		case 0x0005:
			c.doRegisterSUB(opCode)
		case 0x0006:
			c.doRegisterSHR(opCode)
		case 0x0007:
			c.doRegisterSUBN(opCode)
		case 0x0008:
			c.doRegisterSHL(opCode)
		}
	case 0x9000:
		c.skipIfNotEq(opCode)
	case 0xA000:
		c.ANNN(opCode)
	case 0xB000:
		c.BNNN(opCode)
	case 0xC000:
		c.CXKK(opCode)
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
	regA, _ := extractRegisters(opCode)

	for i := byte(0); i <= regA; i++ {
		c.general[i] = c.ram.Read(c.i + ram.Address(i))
	}
}

func (c *CPU) FX55(opCode uint16) {
	regA, _ := extractRegisters(opCode)

	for i := byte(0); i <= regA; i++ {
		c.ram.Write(c.i+ram.Address(i), c.general[i])
	}
}

func (c *CPU) FX33(opCode uint16) {
	regA, _ := extractRegisters(opCode)

	val := c.general[regA]

	c.ram.Write(c.i, val/100)
	c.ram.Write(c.i+1, val%100/10)
	c.ram.Write(c.i+2, val%10)
}

func (c *CPU) FX29(opCode uint16) {
	regA, _ := extractRegisters(opCode)
	c.i = uint16(c.general[regA]) * 5
}

func (c *CPU) FX1E(opCode uint16) {
	regA, _ := extractRegisters(opCode)
	c.i += uint16(c.general[regA])
}

func (c *CPU) FX18(opCode uint16) {
	regA, _ := extractRegisters(opCode)
	c.st = c.general[regA]
}

func (c *CPU) FX15(opCode uint16) {
	regA, _ := extractRegisters(opCode)
	c.dt = c.general[regA]
}

func (c *CPU) FX0A(opCode uint16) {
	regA, _ := extractRegisters(opCode)
	c.general[regA] = c.keyboard.Await()
}

func (c *CPU) FX07(opCode uint16) {
	regA, _ := extractRegisters(opCode)
	c.general[regA] = c.dt
}

func (c *CPU) EXA1(opCode uint16) {
	regA, _ := extractRegisters(opCode)

	if !c.keyboard.IsPressed(byte(regA)) {
		c.pc += 2
	}
}

func (c *CPU) EX9E(opCode uint16) {
	keyNum := (opCode & 0x0F00) >> 8

	if c.keyboard.IsPressed(byte(keyNum)) {
		c.pc += 2
	}
}

func (c *CPU) DXYN(opCode uint16) {
	regA, regB := extractRegisters(opCode)
	n := opCode & 0x000F

	for i := uint16(0); i < n; i++ {
		x := c.general[regA]
		y := c.general[regB]

		sprite := c.ram.Read(c.i + i)

		c.display.Draw(sprite, x, y+byte(i))
	}
}

func (c *CPU) CXKK(opCode uint16) {
	regA, _ := extractRegisters(opCode)

	c.general[regA] = byte(rand.Intn(256)) & getKK(opCode)
}

func (c *CPU) ANNN(opCode uint16) {
	c.i = asAddr(opCode)
}

func (c *CPU) BNNN(opCode uint16) {
	c.pc = asAddr(opCode) + uint16(c.general[0])
}

func (c *CPU) skipIfNotEq(opCode uint16) {
	regA, regB := extractRegisters(opCode)

	if c.general[regB] != c.general[regA] {
		c.pc += 2
	}
}

func (c *CPU) doRegisterSHL(opCode uint16) {
	regA, _ := extractRegisters(opCode)

	if c.general[regA]&0x80 == 0x80 {
		c.vf = 1
	} else {
		c.vf = 0
	}

	c.general[regA] <<= 1
}

func (c *CPU) doRegisterSUBN(opCode uint16) {
	regA, regB := extractRegisters(opCode)

	if c.general[regB] > c.general[regA] {
		c.vf = 1
	} else {
		c.vf = 0
	}

	c.general[regA] = c.general[regB] - c.general[regA]
}

func (c *CPU) doRegisterSHR(opCode uint16) {
	regA, _ := extractRegisters(opCode)

	if c.general[regA]&0x001 == 1 {
		c.vf = 1
	} else {
		c.vf = 0
	}

	c.general[regA] >>= 1
}

func (c *CPU) doRegisterSUB(opCode uint16) {
	regA, regB := extractRegisters(opCode)

	if c.general[regA] > c.general[regB] {
		c.vf = 1
	} else {
		c.vf = 0
	}

	c.general[regA] = c.general[regA] - c.general[regB]
}

func (c *CPU) doRegisterADD(opCode uint16) {
	regA, regB := extractRegisters(opCode)

	res := uint16(c.general[regA]) + uint16(c.general[regB])

	if res > 0xFF {
		c.vf = 1
	} else {
		c.vf = 0
	}

	c.general[regA] = byte(res)
}

func (c *CPU) doRegisterXOR(opCode uint16) {
	regA, regB := extractRegisters(opCode)
	c.general[regA] = c.general[regA] ^ c.general[regB]
}

func (c *CPU) doRegisterOR(opCode uint16) {
	regA, regB := extractRegisters(opCode)
	c.general[regA] = c.general[regA] | c.general[regB]
}

func (c *CPU) doRegisterAND(opCode uint16) {
	regA, regB := extractRegisters(opCode)
	c.general[regA] = c.general[regA] & c.general[regB]
}

func (c *CPU) setRegisterFromOther(opCode uint16) {
	regA, regB := extractRegisters(opCode)

	c.general[regA] = c.general[regB]
}

func (c *CPU) addToRegister(opCode uint16) {
	regN := opCode & 0x0F00
	c.general[regN] += byte(opCode & 0x00FF)
}

func (c *CPU) LD_KK(opCode uint16) {
	regN := opCode & 0x0F00
	c.general[regN] = byte(opCode & 0x00FF)
}

func (c *CPU) SNE_KK(opCode uint16) {
	regN := opCode & 0x0F00
	toCmp := getKK(opCode)

	if c.general[regN] != toCmp {
		c.pc += 2
	}
}

func (c *CPU) SE_KK(opCode uint16) {
	regN := opCode & 0x0F00
	toCmp := getKK(opCode)

	if c.general[regN] == toCmp {
		c.pc += 2
	}
}

func (c *CPU) SE_XY(opCode uint16) {
	regA, regB := extractRegisters(opCode)

	if c.general[regA] == c.general[regB] {
		c.pc += 2
	}
}

func (c *CPU) clearDisplay() {
	panic("impl")
}

func (c *CPU) returnSub() {
	c.pc = c.stack[c.sp]
	c.sp--

	if c.sp < 0 {
		panic("stack pointer is negative")
	}
}

func (c *CPU) JP_NNN(addr ram.Address) {
	c.pc = addr
}

func (c *CPU) CALL_NNN(addr ram.Address) {
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

func extractRegisters(opCode uint16) (byte, byte) {
	regA := opCode & 0x0F00 >> 8
	regB := opCode & 0x00F0 >> 4

	return byte(regA), byte(regB)
}

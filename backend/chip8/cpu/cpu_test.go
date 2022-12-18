package cpu

import (
	"github.com/carlmango11/chip-m8/backend/chip8/display"
	"github.com/carlmango11/chip-m8/backend/chip8/keyboard"
	"github.com/carlmango11/chip-m8/backend/chip8/ram"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	name string

	registers map[byte]byte

	opCode uint16

	expectedVF        byte
	expectedRegisters map[byte]byte
}

func TestCPU(t *testing.T) {
	tcs := []TestCase{
		{
			name: "OR XY",
			registers: map[byte]byte{
				1: 5,
				2: 6,
			},
			opCode: 0x8120,
			expectedRegisters: map[byte]byte{
				1: 6,
				2: 6,
			},
		},
		{
			name: "OR XY",
			registers: map[byte]byte{
				1: 0b1010,
				2: 0b0110,
			},
			opCode: 0x8121,
			expectedRegisters: map[byte]byte{
				1: 0b1110,
			},
		},
		{
			name: "AND XY",
			registers: map[byte]byte{
				1: 0b1010,
				2: 0b0110,
			},
			opCode: 0x8122,
			expectedRegisters: map[byte]byte{
				1: 0b0010,
			},
		},
		{
			name: "XOR XY",
			registers: map[byte]byte{
				1: 0b1010,
				2: 0b0110,
			},
			opCode: 0x8123,
			expectedRegisters: map[byte]byte{
				1: 0b1100,
			},
		},
		{
			name: "ADD XY overflow",
			registers: map[byte]byte{
				1: 255,
				2: 1,
			},
			opCode:     0x8124,
			expectedVF: 1,
			expectedRegisters: map[byte]byte{
				1: 0,
			},
		},
		{
			name: "ADD XY no overflow 1",
			registers: map[byte]byte{
				1: 3,
				2: 2,
			},
			opCode:     0x8124,
			expectedVF: 0,
			expectedRegisters: map[byte]byte{
				1: 5,
			},
		},
		{
			name: "ADD XY no overflow 2",
			registers: map[byte]byte{
				1: 254,
				2: 1,
			},
			opCode:     0x8124,
			expectedVF: 0,
			expectedRegisters: map[byte]byte{
				1: 255,
			},
		},
		{
			name: "SUB greater",
			registers: map[byte]byte{
				1: 5,
				2: 4,
			},
			opCode:     0x8125,
			expectedVF: 1,
			expectedRegisters: map[byte]byte{
				1: 1,
			},
		},
		{
			name: "SUB less",
			registers: map[byte]byte{
				1: 3,
				2: 4,
			},
			opCode:     0x8125,
			expectedVF: 0,
			expectedRegisters: map[byte]byte{
				1: 255,
			},
		},
		{
			name: "SHR with remainder",
			registers: map[byte]byte{
				1: 15,
			},
			opCode:     0x8106,
			expectedVF: 1,
			expectedRegisters: map[byte]byte{
				1: 7,
			},
		},
		{
			name: "SHR without remainder",
			registers: map[byte]byte{
				1: 16,
			},
			opCode:     0x8106,
			expectedVF: 0,
			expectedRegisters: map[byte]byte{
				1: 8,
			},
		},
		{
			name: "SUBN greater",
			registers: map[byte]byte{
				1: 4,
				2: 5,
			},
			opCode:     0x8127,
			expectedVF: 1,
			expectedRegisters: map[byte]byte{
				1: 1,
				2: 5, // unchanged
			},
		},
		{
			name: "SUBN less",
			registers: map[byte]byte{
				1: 5,
				2: 4,
			},
			opCode:     0x8127,
			expectedVF: 0,
			expectedRegisters: map[byte]byte{
				1: 255,
				2: 4, // unchanged
			},
		},
		{
			name: "SHL no remainder",
			registers: map[byte]byte{
				1: 0x01,
			},
			opCode:     0x810E,
			expectedVF: 0,
			expectedRegisters: map[byte]byte{
				1: 0x02,
			},
		},
		{
			name: "SHL remainder",
			registers: map[byte]byte{
				1: 0x81,
			},
			opCode:     0x810E,
			expectedVF: 1,
			expectedRegisters: map[byte]byte{
				1: 0x02,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cpu := New(ConfigChip8, ram.New(nil), &display.Display{}, keyboard.New())

			for reg, val := range tc.registers {
				cpu.v[reg] = val
			}

			cpu.executeOpCode(tc.opCode)

			for reg, val := range tc.expectedRegisters {
				assert.Equal(t, val, cpu.v[reg])
			}

			assert.Equal(t, tc.expectedVF, cpu.v[0xF])
		})
	}
}

func TestLoadToMemory(t *testing.T) {
	c := New(Config{}, ram.New(nil), nil, nil)

	c.cfg.MemoryIncrement = true
	testLoadToMemory(t, c)
	assert.Equal(t, uint16(0x203), c.i)

	resetMemory(c.ram, 0x200, 0x202)

	c.cfg.MemoryIncrement = false
	testLoadToMemory(t, c)
	assert.Equal(t, uint16(0x200), c.i)
}

func testLoadToMemory(t *testing.T, c *CPU) {
	c.i = 0x200

	c.v[0] = 0xA
	c.v[1] = 0xB
	c.v[2] = 0xC

	c.executeOpCode(0xF255)

	assert.Equal(t, byte(0xA), c.ram.Read(0x200))
	assert.Equal(t, byte(0xB), c.ram.Read(0x201))
	assert.Equal(t, byte(0xC), c.ram.Read(0x202))
}

func TestLoadFromMemory(t *testing.T) {
	c := New(Config{}, ram.New(nil), nil, nil)

	c.cfg.MemoryIncrement = true
	testLoadFromMemory(t, c)
	assert.Equal(t, uint16(0x203), c.i)

	resetMemory(c.ram, 0x200, 0x202)

	c.cfg.MemoryIncrement = false
	testLoadFromMemory(t, c)
	assert.Equal(t, uint16(0x200), c.i)
}

func testLoadFromMemory(t *testing.T, c *CPU) {
	c.i = 0x200

	c.ram.Write(0x200, 0xA)
	c.ram.Write(0x201, 0xB)
	c.ram.Write(0x202, 0xC)

	c.executeOpCode(0xF265)

	assert.Equal(t, byte(0xA), c.v[0])
	assert.Equal(t, byte(0xB), c.v[1])
	assert.Equal(t, byte(0xC), c.v[2])
}

func resetMemory(r *ram.RAM, from, to ram.Address) {
	for a := from; a <= to; a++ {
		r.Write(a, 0)
	}
}

package cpu

import (
	"github.com/carlmango11/chip-m8/backend/chip8/display"
	"github.com/carlmango11/chip-m8/backend/chip8/keyboard"
	"github.com/carlmango11/chip-m8/backend/chip8/ram"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoRegisterADD(t *testing.T) {
	cpu := New(ram.New(nil), display.New(), keyboard.New())

	cpu.general[1] = 255
	cpu.general[2] = 2

	cpu.doRegisterADD(0x0120)

	assert.Equal(t, uint8(1), cpu.general[1])
	assert.Equal(t, uint8(1), cpu.vf)
}

func TestDoRegisterSHR(t *testing.T) {
	cpu := New(ram.New(nil), display.New(), keyboard.New())

	cpu.general[1] = 15

	cpu.doRegisterSHR(0x0106)

	assert.Equal(t, uint8(7), cpu.general[1])
	assert.Equal(t, uint8(1), cpu.vf)
}

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
	}

	for _, tc := range tcs {
		cpu := New(ram.New(nil), display.New(), keyboard.New())

		for reg, val := range tc.registers {
			cpu.general[reg] = val
		}

		cpu.executeOpCode(tc.opCode)

		for reg, val := range tc.expectedRegisters {
			assert.Equal(t, val, cpu.general[reg])
		}

		assert.Equal(t, tc.expectedVF, cpu.vf)
	}
}

package display

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrap(t *testing.T) {
	d := &Display{}

	d.Draw(0xFF, 32, 60)

	expect := [32]uint64{}
	expect[0] = 0xF00000000000000F

	assert.Equal(t, expect, d.screen)
}

func TestCountOnes(t *testing.T) {
	assert.Equal(t, 7, countOnes(0b0011010101101))
	assert.Equal(t, 0, countOnes(0))
	assert.Equal(t, 1, countOnes(1))
	assert.Equal(t, 8, countOnes(0b11111111))
}

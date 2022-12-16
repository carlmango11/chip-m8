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

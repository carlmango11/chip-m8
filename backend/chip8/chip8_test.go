package chip8

import (
	"encoding/hex"
	"strings"
	"testing"
)

const testRom = `
124eeaacaaeaceaaaaaee0a0a0e0c04040e0e020c0e0e06020e0a0e02020
60402040e080e0e0e0202020e0e0a0e0e0e020e040a0e0a0e0c080e0e080
c080a040a0a0a202dab400eea202dab413dc680169056a0a6b01652a662b
a216d8b4a23ed9b4a202362ba206dab46b06a21ad8b4a23ed9b4a206452a
a202dab46b0ba21ed8b4a23ed9b4a2065560a202dab46b10a226d8b4a23e
d9b4a20676ff462aa202dab46b15a22ed8b4a23ed9b4a2069560a202dab4
6b1aa232d8b4a23ed9b422426817691b6a206b01a20ad8b4a236d9b4a202
dab46b06a22ad8b4a20ad9b4a2068750472aa202dab46b0ba22ad8b4a20e
d9b4a206672a87b1472ba202dab46b10a22ad8b4a212d9b4a2066678671f
87624718a202dab46b15a22ad8b4a216d9b4a2066678671f87634767a202
dab46b1aa22ad8b4a21ad9b4a206668c678c87644718a202dab4682c6930
6a346b01a22ad8b4a21ed9b4a206668c6778876547eca202dab46b06a22a
d8b4a222d9b4a20666e0866e46c0a202dab46b0ba22ad8b4a236d9b4a206
660f86664607a202dab46b10a23ad8b4a21ed9b4a3e860006130f155a3e9
f065a2064030a202dab46b15a23ad8b4a216d9b4a3e86689f633f265a202
3001a2063103a2063207a206dab46b1aa20ed8b4a23ed9b4124813dc
`

func TestPredefinedSprites(t *testing.T) {
	script := prepScript("6003 F029 D125")

	c := New(script)

	for range script {
		c.cpu.Tick()
	}

	// should be a 3 drawn on the display
	c.display.Print()
}

func TestRom(t *testing.T) {
	script := prepScript(testRom)

	c := New(script)

	for range script {
		c.cpu.Tick()
	}

	// should be a 3 drawn on the display
	c.display.Print()
}

func prepScript(s string) []byte {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")

	script, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return script
}

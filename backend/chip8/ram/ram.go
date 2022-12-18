package ram

import (
	"fmt"
)

const UserMemoryStart = 0x200

type Address = uint16

type RAM struct {
	data [4096]byte
}

func New(script []byte) *RAM {
	return &RAM{
		data: initData(script),
	}
}

func (r *RAM) Read(addr Address) byte {
	if addr > 0xFFF {
		panic(fmt.Sprintf("cannot access %v", addr))
	}

	return r.data[addr]
}

func (r *RAM) Write(addr Address, val byte) {
	r.data[addr] = val
}

func initData(script []byte) [4096]byte {
	data := [4096]byte{}

	for n, rows := range characters {
		for i, row := range rows {
			data[(n*5)+i] = row
		}
	}

	// put program in place
	for i, val := range script {
		data[UserMemoryStart+i] = val
	}

	// auto start debug
	//data[0x1FF] = 5
	//data[0x1FE] = 3

	return data
}

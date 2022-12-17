package main

import (
	"encoding/hex"
	"fmt"
	"github.com/carlmango11/chip-m8/backend/chip8"
	"syscall/js"
)

func main() {
	createBindings()

	waitC := make(chan bool)
	<-waitC
}

func createBindings() {
	vm := chip8.New(nil)

	loadFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		script, err := hex.DecodeString(args[0].String())
		if err != nil {
			return err
		}

		vm = chip8.New(script)

		return nil
	})

	keyPressedFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		n := args[0].Int()
		vm.PressKey(byte(n))

		return nil
	})

	keyUnpressedFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		n := args[0].Int()
		vm.UnpressKey(byte(n))

		return nil
	})

	getDisplayFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		// doesn't support normal typed arrays, only interface{}
		var result []interface{}
		for _, row := range vm.Display() {
			result = append(result, fmt.Sprintf("%064b", row))
		}

		return result
	})

	js.Global().Set("loadScript", loadFunc)
	js.Global().Set("keyPressed", keyPressedFunc)
	js.Global().Set("keyUnpressed", keyUnpressedFunc)
	js.Global().Set("getDisplay", getDisplayFunc)
}

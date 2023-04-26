package main

import (
	"bytes"
	"syscall/js"
)

func doProcessForWasm(_ js.Value, args []js.Value) any {
	// convert
	source := args[0].String()

	output := bytes.NewBufferString("")
	input := bytes.NewBufferString(source)
	err := doProcess(input, output, "example")
	if err != nil {
		return js.ValueOf("Error: " + err.Error())
	}
	return js.ValueOf(output.String())
}

func main() {
	js.Global().Set("doProcessForWasm", js.FuncOf(doProcessForWasm))
}

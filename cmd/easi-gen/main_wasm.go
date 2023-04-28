package main

import (
	"bytes"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/cmd_shared"
	"syscall/js"
)

func doProcessForWasm(_ js.Value, args []js.Value) any {
	// convert
	source := args[0].String()

	output := bytes.NewBufferString("")
	input := bytes.NewBufferString(source)
	err := cmd_shared.DoProcess(input, output, "example")
	if err != nil {
		return js.ValueOf("Error: " + err.Error())
	}
	return js.ValueOf(output.String())
}

func main() {
	js.Global().Set("doProcessForWasm", js.FuncOf(doProcessForWasm))
	js.Global().Get("console").Call("log", "[GO] function doProcessForWasm injected")
	select {}
}

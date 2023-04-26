package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ONLY used to generate package.json
// usage: go run gen_pkg_json.go <version>
func main() {
	version := os.Args[1]
	out, _ := json.MarshalIndent(map[string]any{
		"name":        "@lemonneko/easi-gen",
		"version":     version,
		"description": "Easinteraction is a tool that help users to generate code for easier contract interaction. This is wasm version, used to generate go code for single contract in browser.",
		"keywords":    []string{"flow", "blockchain", "codegen", "go", "wasm", "code generation"},
		"author":      "lemonneko",
		"license":     "MIT",
		"homepage":    "https://docs.easi-gen.lemonneko.moe",
		"repository":  "https://github.com/lemonneko/easi-gen",
	}, "", "  ")
	fmt.Printf("%s\n", out)
}

import wasmInit from '../dist/easi-gen.wasm?init'

globalThis.require = require;
// @ts-ignore copied from $GOROOT/misc/wasm/wasm_exec_node.js
globalThis.fs = require("fs");
globalThis.TextEncoder = require("util").TextEncoder;
globalThis.TextDecoder = require("util").TextDecoder;

// @ts-ignore copied from $GOROOT/misc/wasm/wasm_exec_node.js
globalThis.performance = {
	now() {
		const [sec, nsec] = process.hrtime();
		return sec * 1000 + nsec / 1000000;
	},
};

const crypto = require("crypto");
globalThis.crypto = {
    // @ts-ignore copied from $GOROOT/misc/wasm/wasm_exec_node.js
	getRandomValues(b) {
		crypto.randomFillSync(b);
	},
};

export const newEasiGen = async () => {
    // @ts-ignore
    await import('./wasm_exec.mjs')

    const instance = await wasmInit({})
    const go = new Go()
    await go.run(instance)
    return (source: string): string => {
        return globalThis.doProcessForWasm(source)
    }
}

import wasmInit from '../dist/easi-gen.wasm?init'

const crypto = require("crypto");
globalThis.crypto = {
    // @ts-ignore because it from go/misc/wasm
	getRandomValues(b) {
		crypto.randomFillSync(b);
	},
};

import "../dist/wasm_exec.js"

export const newEasiGen = async () => {
    const instance = await wasmInit({})
    const go = new Go()
    await go.run(instance)
    return (source: string): string => {
        return globalThis.doProcessForWasm(source)
    }
}

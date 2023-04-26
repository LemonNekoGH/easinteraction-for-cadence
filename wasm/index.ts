import wasmInit from '../dist/easi-gen.wasm?init'
import "../dist/wasm_exec.js"

export const newEasiGen = async () => {
    const instance = await wasmInit({})
    const go = new Go()
    await go.run(instance)
    return (source: string): string => {
        return globalThis.doProcessForWasm(source)
    }
}

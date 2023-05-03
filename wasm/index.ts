import wasmInit from '../dist/easi-gen.wasm?init'
// @ts-ignore
import { go } from './wasm_exec.js'

export const newEasiGen = async () => {
    const instance = await wasmInit(go.importObject)
    go.run(instance).then()
    return (source: string, ignoreContractGeneration: boolean): string => {
        return doProcessForWasm(source, ignoreContractGeneration)
    }
}

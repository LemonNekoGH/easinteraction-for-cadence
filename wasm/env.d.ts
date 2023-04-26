declare module "*.wasm?init" {
    export default function (options: WebAssembly.Imports): Promise<WebAssembly.Instance>
}

declare module "wasm_exec" {
    const value: never
    export default value
}

declare class Go {
    run(a: WebAssembly.Instance): Promise<any>
}

declare var doProcessForWasm: (source: string) => string

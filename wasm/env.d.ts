declare module "*.wasm?init" {
    export default function (options: WebAssembly.Imports): Promise<WebAssembly.Instance>
}

declare module "*.cdc?raw" {
    const value: string
    export default value
}

declare class Go {
    run(a: WebAssembly.Instance): Promise<any>
    importObject: any
}

declare var doProcessForWasm: (source: string, ignoreContractGeneration: boolean) => string

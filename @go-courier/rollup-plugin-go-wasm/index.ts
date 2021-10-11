import {CustomPluginOptions, LoadResult, Plugin, ResolveIdResult, TransformResult} from "rollup";
import {execSync} from "child_process"
import {resolve, join, dirname, isAbsolute} from "path"
import {readFileSync} from "fs"
import * as process from "process";

const goroot = String(execSync(`go env GOROOT`)).trim()
const execwasm = readFileSync(join(goroot, "./misc/wasm/wasm_exec.js"))

export default () => {
    const reGoFile = /\.go$/
    const goWasmMap = {}

    return {
        name: "go-wasm",

        transform(src: string, id: string): TransformResult {
            if (reGoFile.test(id)) {
                const outputWasm = join(dirname(id), `bin/main.wasm`)

                // re build only on contents change
                if (!goWasmMap[id] || goWasmMap[id] != src) {
                    execSync(`GOOS=js GOARCH=wasm go build -o ${outputWasm} ${id}`)
                    goWasmMap[id] = src
                }

                return `
import wasm from "${outputWasm}?url"               
                
${execwasm}                
                
export const main = async () => {
    const go = new Go();
    const importMeta = import.meta

    return WebAssembly
        .instantiateStreaming(
            fetch(wasm),
            go.importObject
        )
        .then((result) => go.run(result.instance));
}
`
            }
            return
        }
    }
}
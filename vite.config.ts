import {defineConfig} from "vite"
import {VitePWA} from "vite-plugin-pwa"
import goWasm from "@go-courier/rollup-plugin-go-wasm"

const basePath = process.env.BASE_PATH || "/"

export default defineConfig({
    root: "./cmd/webapp",
    base: basePath,
    build: {
        assetsDir: "static",
    },
    plugins: [
        goWasm({
            importWasmSuffix: "?url",
        }),
        VitePWA({
            workbox: {
                maximumFileSizeToCacheInBytes: 50000000,
                globPatterns: [
                    "static/*.*",
                ],
            },
        }),
    ]
})
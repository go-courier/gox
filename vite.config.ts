import {defineConfig} from "vite"
import {VitePWA} from 'vite-plugin-pwa'
import GoWasmPack from "./@go-courier/rollup-plugin-go-wasm"

const basePath = process.env.BASE_PATH || "/"

export default defineConfig({
    root: "./cmd/webapp",
    base: basePath,
    build: {
        assetsDir: "static",
    },
    plugins: [
        VitePWA({
            workbox: {
                maximumFileSizeToCacheInBytes: 50000000,
                globPatterns: [
                    "static/*.*",
                ],
            },
        }),
        GoWasmPack(),
    ]
})

// @ts-ignore
import ("./main.go").then(({main}) => main());

// @ts-ignore
import {registerSW} from "virtual:pwa-register"

const updateSW = registerSW({
    onOfflineReady() {
    },
})
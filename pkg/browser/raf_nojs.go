//go:build !js
// +build !js

package browser

// just fake here in go
func RequestAnimationFrame(fn func()) {
	fn()
}

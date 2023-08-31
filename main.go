package main

import (
	"github.com/VyacheArt/edl-to-youtube/converter"
	"runtime/debug"
)

func main() {
	debug.SetMemoryLimit(128 << 20) // 128 MB

	a := converter.NewApplication()
	if err := a.Run(); err != nil {
		panic(err)
	}
}

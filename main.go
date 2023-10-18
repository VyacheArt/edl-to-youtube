/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

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

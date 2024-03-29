/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package converter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func EnableProgress(c []fyne.Canvas, op func(), text string) {
	contents := make([]fyne.CanvasObject, len(c))
	for i := range c {
		contents[i] = c[i].Content()
		c[i].SetContent(widget.NewModalPopUp(GetProgressView(text), c[i]))
	}

	op()

	for i := range c {
		c[i].SetContent(contents[i])
	}
}

func GetProgressView(text string) fyne.CanvasObject {
	if len(text) == 0 {
		text = "Waiting..."
	}

	return container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			widget.NewLabel(text),
			widget.NewProgressBarInfinite(),
		),
	)
}

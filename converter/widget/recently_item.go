/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os"
	"strings"
)

type (
	RecentlyItem struct {
		fyne.CanvasObject

		path                  string
		isAvailable           bool
		title                 *canvas.Text
		fullPath              *canvas.Text
		removeButton          *widget.Button
		unavailableHelpButton *widget.Button
		removeCallback        ClickRecentlyCallback
		unavailableCallback   ClickRecentlyCallback
	}

	ClickRecentlyCallback func(path string)
)

func NewRecentlyItem(removeCallback ClickRecentlyCallback, unavailableCallback ClickRecentlyCallback) *RecentlyItem {
	item := &RecentlyItem{
		title:               canvas.NewText("path", theme.ForegroundColor()),
		fullPath:            canvas.NewText("full path", theme.ForegroundColor()),
		removeCallback:      removeCallback,
		unavailableCallback: unavailableCallback,
	}

	item.title.TextStyle = fyne.TextStyle{Bold: true}
	item.title.TextSize = 15
	item.fullPath.TextSize = 10

	item.unavailableHelpButton = widget.NewButtonWithIcon("", theme.WarningIcon(), func() {
		item.unavailableCallback(item.path)
	})
	item.unavailableHelpButton.Hide()

	item.removeButton = widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		item.removeCallback(item.path)
	})

	item.CanvasObject = container.NewPadded(
		container.NewBorder(nil, nil, nil, container.NewPadded(container.NewHBox(item.unavailableHelpButton, item.removeButton)),
			container.NewVBox(item.title, container.NewPadded(item.fullPath)),
		),
	)

	return item
}

func (r *RecentlyItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.CanvasObject)
}

func (r *RecentlyItem) SetPath(path string) {
	r.path = path
	r.refresh()
}

func (r *RecentlyItem) SetIsAvailable(isAvailable bool) {
	r.isAvailable = isAvailable
	r.refresh()
}

func (r *RecentlyItem) SetRemoveCallback(callback func(path string)) {
	r.removeCallback = callback
}

func (r *RecentlyItem) refresh() {
	//get file name from path
	parts := strings.Split(r.path, "/")
	r.title.Text = parts[len(parts)-1]
	r.fullPath.Text = makePathShorter(r.path)

	if r.isAvailable {
		r.unavailableHelpButton.Hide()
	} else {
		r.unavailableHelpButton.Show()
	}
}

func makePathShorter(path string) string {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}

	return path
}

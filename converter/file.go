/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package converter

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/VyacheArt/edl-to-youtube/edl"
	"github.com/ncruces/zenity"
	"log"
	"os"
)

func ChooseFile(app *Application, window fyne.Window) bool {
	var path string

	windows := app.getApp().Driver().AllWindows()
	canvases := make([]fyne.Canvas, len(windows))
	for i, w := range windows {
		canvases[i] = w.Canvas()
	}

	EnableProgress(canvases, func() {
		path, _ = zenity.SelectFile(zenity.FileFilters{
			{
				Name:     "EDL files",
				Patterns: []string{"*.edl", "*.EDL"},
			},
		}, zenity.Modal())
	}, "Choose EDL file")

	return OpenFile(app, window, path)
}

func OpenFile(app *Application, window fyne.Window, path string) bool {
	if len(path) == 0 {
		return false
	}

	log.Printf("File chosen, path: %s", path)

	//get file size
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		dialog.ShowError(fmt.Errorf("file not found: %s", path), window)
		return false
	}

	if err != nil {
		dialog.ShowError(fmt.Errorf("unexpected error: %e", err), window)
		return false
	}

	if info.Size() > MaxFileSize {
		dialog.ShowError(fmt.Errorf("file size is too big, max size is %d bytes", MaxFileSize), window)
		return false
	}

	content, err := os.ReadFile(path)
	if err != nil {
		dialog.ShowError(err, window)
		return false
	}

	edlList, err := edl.Parse(string(content))
	if err != nil {
		dialog.ShowError(err, window)
		return false
	}

	if len(edlList.Clips) == 0 {
		dialog.ShowError(errors.New("no clips found in the EDL file"), window)
		return false
	}

	viewer := NewViewerWindow(app, edlList)
	viewer.Show()

	app.appendLastOpened(path)

	return true
}

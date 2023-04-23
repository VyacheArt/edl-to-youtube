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
	EnableProgress(window.Canvas(), func() {
		path, _ = zenity.SelectFile(zenity.FileFilters{
			{
				Name:     "All files",
				Patterns: []string{"*"},
			},
			{
				Name:     "EDL files",
				Patterns: []string{"edl"},
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
	if info, _ := os.Stat(path); info.Size() > MaxFileSize {
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

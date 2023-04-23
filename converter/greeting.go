package converter

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/VyacheArt/edl-to-youtube/edl"
	"github.com/ncruces/zenity"
	"log"
	"os"
)

const MaxFileSize = 100 * 1024 * 1024 // 100 MB

type GreetingWindow struct {
	app    fyne.App
	window fyne.Window
}

func NewGreetingWindow(app fyne.App) *GreetingWindow {
	return &GreetingWindow{
		app: app,
	}
}

func (w *GreetingWindow) Show() {
	w.window = w.app.NewWindow(fmt.Sprintf("Welcome to %s", Title))

	w.window.SetFixedSize(true)
	w.window.Resize(fyne.NewSize(400, 300))
	w.window.CenterOnScreen()

	w.window.SetContent(w.getContent())
	w.window.Show()
}

func (w *GreetingWindow) getContent() fyne.CanvasObject {
	content := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			widget.NewButton("Choose file", func() {
				EnableProgress(w.window.Canvas(), func() {
					path, _ := zenity.SelectFile(zenity.FileFilters{
						{
							Name:     "All files",
							Patterns: []string{"*"},
						},
						{
							Name:     "EDL files",
							Patterns: []string{"edl"},
						},
					}, zenity.Modal())

					w.onChooseFile(path)
				}, "Choose EDL file")
			}),
		),
	)

	return content
}

func (w *GreetingWindow) onChooseFile(path string) {
	if len(path) == 0 {
		return
	}

	log.Printf("File chosen, path: %s", path)

	//get file size
	if info, _ := os.Stat(path); info.Size() > MaxFileSize {
		dialog.ShowError(fmt.Errorf("file size is too big, max size is %d bytes", MaxFileSize), w.window)
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}

	edlList, err := edl.Parse(string(content))
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}

	if len(edlList.Clips) == 0 {
		dialog.ShowError(errors.New("no clips found in the EDL file"), w.window)
		return
	}

	w.window.Hide()

	viewer := NewViewerWindow(w.app, edlList)
	viewer.Show()

	w.window.Close()
}

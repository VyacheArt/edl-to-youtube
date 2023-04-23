package converter

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"
)

const MaxFileSize = 100 * 1024 * 1024 // 100 MB

type GreetingWindow struct {
	app    *Application
	window fyne.Window
}

func NewGreetingWindow(app *Application) *GreetingWindow {
	return &GreetingWindow{
		app: app,
	}
}

func (w *GreetingWindow) Show() {
	w.window = w.app.getApp().NewWindow(Title)

	w.window.Resize(fyne.NewSize(400, 300))
	w.window.CenterOnScreen()

	w.window.SetContent(w.getContent())
	w.window.SetMainMenu(GetMainMenu(w.app, w.window))
	w.window.Show()
}

func (w *GreetingWindow) getContent() fyne.CanvasObject {
	const greetingText = "Welcome to %s!\n\nPlease choose EDL file to make YouTube captions:"

	chooseFileContainer := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			widget.NewLabel(fmt.Sprintf(greetingText, Title)),
			widget.NewButton("Choose file", func() {
				if ChooseFile(w.app, w.window) {
					w.window.Close()
				}
			}),
		),
	)

	lastOpened := w.app.getLastOpened()
	if len(lastOpened) == 0 {
		return chooseFileContainer
	}

	fileList := widget.NewList(
		func() int {
			return len(lastOpened)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("path")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(getNameFromPath(lastOpened[i]))
		},
	)

	fileList.OnSelected = func(id widget.ListItemID) {
		if OpenFile(w.app, w.window, lastOpened[id]) {
			w.window.Close()
		}
	}

	content := container.NewHSplit(fileList, chooseFileContainer)
	content.SetOffset(0.3)

	w.window.Resize(fyne.NewSize(640, 480))

	return content
}

func getNameFromPath(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

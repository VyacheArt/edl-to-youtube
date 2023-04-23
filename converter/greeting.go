package converter

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

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
	w.window.Resize(fyne.NewSize(640, 480))
	w.window.CenterOnScreen()

	w.window.SetContent(w.getContent())
	w.window.ShowAndRun()
}

func (w *GreetingWindow) getContent() fyne.CanvasObject {
	content := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			widget.NewButton("Choose file", func() {
				//fileWindow := w.app.NewWindow("Choose file")
				//fileWindow.Resize(fyne.NewSize(640, 480))
				//
				//dialog.ShowFileOpen()
				dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {

				}, w.window)

				//fileWindow.Show()
			}),
		),
	)

	return content
}

func (w *GreetingWindow) onChooseFile() {

}

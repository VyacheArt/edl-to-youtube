package converter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Application struct {
	app    fyne.App
	window fyne.Window
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Run() error {
	a.app = app.New()

	NewGreetingWindow(a.app).Show()

	return nil
}

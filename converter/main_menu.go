package converter

import "fyne.io/fyne/v2"

func GetMainMenu(app *Application, window fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open", func() {
				ChooseFile(app, window)
			}),
			fyne.NewMenuItem("Welcome screen", func() {
				NewGreetingWindow(app).Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				app.getApp().Quit()
			}),
		),
	)
}

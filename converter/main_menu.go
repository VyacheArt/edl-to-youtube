package converter

import "fyne.io/fyne/v2"

func GetMainMenu(app fyne.App, window fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open", func() {
				ChooseFile(app, window)
			}),
		),
	)
}

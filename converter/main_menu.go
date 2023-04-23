package converter

import "fyne.io/fyne/v2"

func GetMainMenu(app *Application, window fyne.Window) *fyne.MainMenu {
	lastOpened := app.getLastOpened()
	lastOpenedItems := make([]*fyne.MenuItem, 0, len(lastOpened))
	for _, path := range lastOpened {
		lastOpenedItems = append(lastOpenedItems, fyne.NewMenuItem(path, func() {
			OpenFile(app, window, path)
		}))
	}

	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open", func() {
				ChooseFile(app, window)
			}),
			&fyne.MenuItem{
				Label:     "Open last",
				ChildMenu: fyne.NewMenu("", lastOpenedItems...),
				Disabled:  len(lastOpened) == 0,
			},
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

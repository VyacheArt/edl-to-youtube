package converter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	. "github.com/VyacheArt/edl-to-youtube/converter/locale"
)

func GetMainMenu(app *Application, window fyne.Window) *fyne.MainMenu {
	const (
		fileKey          = "file"
		openKey          = "open"
		openLastKey      = "openLast"
		languageKey      = "language"
		welcomeScreenKey = "welcomeScreen"
		quitKey          = "quit"
	)

	lastOpened := app.getLastOpened()
	lastOpenedItems := make([]*fyne.MenuItem, 0, len(lastOpened))
	for _, path := range lastOpened {
		lastOpenedItems = append(lastOpenedItems, fyne.NewMenuItem(path, func() {
			OpenFile(app, window, path)
		}))
	}

	locales := GetLocales()
	localesItems := make([]*fyne.MenuItem, 0, len(locales))
	for key, title := range locales {
		localesItems = append(localesItems,
			fyne.NewMenuItem(title, func() {
				app.setLocale(key)
				dialog.ShowInformation("Restart", "Please restart the application to apply the changes", window)
			}),
		)
	}

	return fyne.NewMainMenu(
		fyne.NewMenu(Localize(fileKey),
			fyne.NewMenuItem(Localize(openKey), func() {
				ChooseFile(app, window)
			}),
			&fyne.MenuItem{
				Label:     Localize(openLastKey),
				ChildMenu: fyne.NewMenu("", lastOpenedItems...),
				Disabled:  len(lastOpened) == 0,
			},
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(Localize(welcomeScreenKey), func() {
				NewGreetingWindow(app).Show()
			}),
			&fyne.MenuItem{
				Label:     Localize(languageKey),
				ChildMenu: fyne.NewMenu("", localesItems...),
			},
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(Localize(quitKey), func() {
				app.getApp().Quit()
			}),
		),
	)
}

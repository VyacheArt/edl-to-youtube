package converter

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/VyacheArt/edl-to-youtube/converter/locale"
	widget2 "github.com/VyacheArt/edl-to-youtube/converter/widget"
)

const MaxFileSize = 100 * 1024 * 1024 // 100 MB

type GreetingWindow struct {
	app    *Application
	window fyne.Window

	lastOpened     []string
	lastOpenedList *widget.List
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
	w.window.SetCloseIntercept(func() {
		w.window.Hide()
	})

	w.lastOpened = w.app.getLastOpened()

	w.initLastOpenedList()

	w.window.SetContent(w.getContent())
	w.window.SetMainMenu(w.getMainMenu())
	w.window.Show()
}

func (w *GreetingWindow) getMainMenu() *fyne.MainMenu {
	lastOpened := w.app.getLastOpened()
	lastOpenedItems := make([]*fyne.MenuItem, 0, len(lastOpened))
	for _, path := range lastOpened {
		lastOpenedItems = append(lastOpenedItems, fyne.NewMenuItem(path, func() {
			OpenFile(w.app, w.window, path)
		}))
	}

	locales := GetLocales()
	localesItems := make([]*fyne.MenuItem, 0, len(locales))
	for key, title := range locales {
		key, title := key, title
		localesItems = append(localesItems,
			fyne.NewMenuItem(title, func() {
				w.app.setLocale(key)
				dialog.ShowInformation("Restart", "Please restart the application to apply the changes", w.app.getCurrentWindow())
			}),
		)
	}

	return fyne.NewMainMenu(
		fyne.NewMenu(Localize(MenuFile),
			fyne.NewMenuItem(Localize(MenuOpen), func() {
				ChooseFile(w.app, w.window)
			}),
			&fyne.MenuItem{
				Label:     Localize(MenuOpenLast),
				ChildMenu: fyne.NewMenu("", lastOpenedItems...),
				Disabled:  len(lastOpened) == 0,
			},
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(Localize(MenuWelcomeScreen), func() {
				w.window.Show()
				w.refreshRecentlyList()
			}),
			&fyne.MenuItem{
				Label:     Localize(MenuLanguage),
				ChildMenu: fyne.NewMenu("", localesItems...),
			},
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(Localize(MenuQuit), func() {
				w.app.getApp().Quit()
			}),
		),

		fyne.NewMenu(Localize(MenuHelp),
			&fyne.MenuItem{
				Label: Localize(HelpHowTo),
				ChildMenu: fyne.NewMenu("",
					fyne.NewMenuItem(Localize(HelpExportInResolve), func() {
						NewResolveManualWindow(w.app).Show()
					}),
				),
			},
		),
	)
}

func (w *GreetingWindow) getContent() fyne.CanvasObject {
	const (
		greetingTextKey = "greetingMessage"
		chooseFileKey   = "chooseFile"
	)

	chooseFileContainer := container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			widget.NewLabel(fmt.Sprintf(Localize(greetingTextKey), Title)),
			widget.NewButton(Localize(chooseFileKey), func() {
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

	content := container.NewVSplit(chooseFileContainer, w.lastOpenedList)
	content.SetOffset(0.3)

	w.window.Resize(fyne.NewSize(640, 480))

	return content
}

func (w *GreetingWindow) initLastOpenedList() {
	fileList := widget.NewList(
		func() int {
			return len(w.lastOpened)
		},
		func() fyne.CanvasObject {
			return widget2.NewRecentlyItem(w.onRecentlyItemDeleted)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget2.RecentlyItem).SetPath(w.lastOpened[i])
		},
	)

	fileList.OnSelected = func(id widget.ListItemID) {
		if OpenFile(w.app, w.window, w.lastOpened[id]) {
			w.window.Hide()
		}
	}

	w.lastOpenedList = fileList
}

func (w *GreetingWindow) onRecentlyItemDeleted(path string) {
	w.app.removeLastOpened(path)
	w.refreshRecentlyList()
}

func (w *GreetingWindow) refreshRecentlyList() {
	w.lastOpened = w.app.getLastOpened()
	w.lastOpenedList.Refresh()
	w.lastOpenedList.UnselectAll()
}

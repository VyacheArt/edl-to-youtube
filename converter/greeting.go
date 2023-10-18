/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

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
	"os"
	"sync"
)

const MaxFileSize = 100 << 20 // 100 MB

type (
	GreetingWindow struct {
		app    *Application
		window fyne.Window

		mu                  sync.Mutex
		lastOpened          []*recentlyItem
		lastOpenedList      *widget.List
		chooseFileContainer *fyne.Container
		splitContainer      *container.Split
	}

	recentlyItem struct {
		Path        string
		IsAvailable bool
	}
)

func NewGreetingWindow(app *Application) *GreetingWindow {
	return &GreetingWindow{
		app: app,
	}
}

func newRecentlyItems(paths []string) []*recentlyItem {
	items := make([]*recentlyItem, len(paths))
	for i, path := range paths {
		items[i] = newRecentlyItem(path)
	}
	return items
}

func newRecentlyItem(path string) *recentlyItem {
	f, err := os.Open(path)
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	return &recentlyItem{
		Path:        path,
		IsAvailable: err == nil,
	}
}

func (w *GreetingWindow) Show() {
	w.window = w.app.getApp().NewWindow(AppTitle)

	w.window.Resize(fyne.NewSize(400, 300))
	w.window.CenterOnScreen()
	w.window.SetCloseIntercept(func() {
		if w.app.getWindowCount() > 1 {
			w.window.Hide()
		} else {
			w.app.Quit()
		}
	})

	w.initContent()
	w.refreshMenu()
	w.refreshContent()

	w.window.SetContent(w.getContent())
	w.window.Resize(fyne.NewSize(640, 480))
	w.window.Show()
}

func (w *GreetingWindow) refreshMenu() {
	w.window.SetMainMenu(w.getMainMenu())
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
				w.onChooseFile()
			}),
			&fyne.MenuItem{
				Label:     Localize(MenuOpenLast),
				ChildMenu: fyne.NewMenu("", lastOpenedItems...),
				Disabled:  len(lastOpened) == 0,
			},
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(Localize(MenuWelcomeScreen), func() {
				w.window.Show()
				w.refreshContent()
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

func (w *GreetingWindow) initContent() {
	const (
		greetingTextKey = "greetingMessage"
		chooseFileKey   = "chooseFile"
	)

	w.initLastOpenedList()

	w.chooseFileContainer = container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			widget.NewLabel(fmt.Sprintf(Localize(greetingTextKey), AppTitle)),
			widget.NewButton(Localize(chooseFileKey), func() {
				w.onChooseFile()
			}),
		),
	)

	w.splitContainer = container.NewVSplit(w.chooseFileContainer, w.lastOpenedList)
	w.splitContainer.SetOffset(0.3)
}

func (w *GreetingWindow) refreshContent() {
	w.refreshLastOpened()

	w.window.SetContent(w.getContent())

	w.chooseFileContainer.Refresh()
	w.chooseFileContainer.Show()
	w.lastOpenedList.Refresh()
	w.lastOpenedList.UnselectAll()

}

func (w *GreetingWindow) getContent() fyne.CanvasObject {
	lastOpened := w.app.getLastOpened()
	if len(lastOpened) == 0 {
		return w.chooseFileContainer
	}

	return w.splitContainer
}

func (w *GreetingWindow) initLastOpenedList() {
	fileList := widget.NewList(
		func() int {
			return len(w.lastOpened)
		},
		func() fyne.CanvasObject {
			return widget2.NewRecentlyItem(w.onRecentlyItemDeleted, w.onRecentlyItemUnavailable)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			ri := o.(*widget2.RecentlyItem)
			ri.SetPath(w.lastOpened[i].Path)
			ri.SetIsAvailable(w.lastOpened[i].IsAvailable)
		},
	)

	fileList.OnSelected = func(id widget.ListItemID) {
		fileList.UnselectAll()

		item := w.lastOpened[id]
		if item.IsAvailable {
			if OpenFile(w.app, w.window, w.lastOpened[id].Path) {
				w.window.Hide()
			}
		} else {
			w.onChooseFile()
		}
	}

	w.lastOpenedList = fileList
}

func (w *GreetingWindow) onRecentlyItemDeleted(path string) {
	w.app.removeLastOpened(path)
	w.refreshContent()
	w.refreshMenu()
}

func (w *GreetingWindow) onChooseFile() {
	if ChooseFile(w.app, w.window) {
		w.window.Hide()
		w.refreshMenu()
	}
}

func (w *GreetingWindow) onRecentlyItemUnavailable(path string) {
	dialog.ShowInformation(Localize(HelpUnavailableFileTitle), Localize(HelpUnavailableFile), w.window)
}

func (w *GreetingWindow) refreshLastOpened() {
	lastOpened := w.app.getLastOpened()

	w.mu.Lock()
	defer w.mu.Unlock()

	w.lastOpened = newRecentlyItems(lastOpened)
}

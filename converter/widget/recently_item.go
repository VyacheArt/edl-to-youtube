package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os"
	"strings"
)

type (
	RecentlyItem struct {
		fyne.CanvasObject

		path           string
		title          *canvas.Text
		fullPath       *canvas.Text
		removeCallback RemoveCallback
	}

	RemoveCallback func(path string)
)

func NewRecentlyItem(removeCallback RemoveCallback) *RecentlyItem {
	item := &RecentlyItem{
		title:          canvas.NewText("path", theme.ForegroundColor()),
		fullPath:       canvas.NewText("full path", theme.ForegroundColor()),
		removeCallback: removeCallback,
	}
	item.title.TextStyle = fyne.TextStyle{Bold: true}
	item.title.TextSize = 15
	item.fullPath.TextSize = 10

	removeButton := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		item.removeCallback(item.path)
	})
	item.CanvasObject = container.NewPadded(
		container.NewBorder(nil, nil, nil, container.NewPadded(removeButton),
			container.NewVBox(item.title, container.NewPadded(item.fullPath)),
		),
	)

	return item
}

func (r *RecentlyItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.CanvasObject)
}

func (r *RecentlyItem) SetPath(path string) {
	r.path = path
	r.refresh()
}

func (r *RecentlyItem) SetRemoveCallback(callback func(path string)) {
	r.removeCallback = callback
}

func (r *RecentlyItem) refresh() {
	//get file name from path
	parts := strings.Split(r.path, "/")
	r.title.Text = parts[len(parts)-1]
	r.fullPath.Text = makePathShorter(r.path)
}

func makePathShorter(path string) string {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}

	return path
}

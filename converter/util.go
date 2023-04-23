package converter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func EnableProgress(c fyne.Canvas, op func(), text string) {
	content := c.Content()
	c.SetContent(widget.NewModalPopUp(GetProgressView(text), c))

	op()

	c.SetContent(content)
}

func GetProgressView(text string) fyne.CanvasObject {
	if len(text) == 0 {
		text = "Waiting..."
	}

	return container.New(layout.NewCenterLayout(),
		container.New(layout.NewVBoxLayout(),
			widget.NewLabel(text),
			widget.NewProgressBarInfinite(),
		),
	)
}

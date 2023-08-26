package converter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/VyacheArt/edl-to-youtube/converter/locale"
	"github.com/VyacheArt/edl-to-youtube/resources"
)

const ResolveManualTitle = "DaVinci Resolve Manual"

type ResolveManualWindow struct {
	app    *Application
	window fyne.Window
}

func NewResolveManualWindow(app *Application) *ResolveManualWindow {
	return &ResolveManualWindow{
		app: app,
	}
}

func (w *ResolveManualWindow) Show() {
	w.window = w.app.getApp().NewWindow(ResolveManualTitle)

	w.window.Resize(fyne.NewSize(600, 800))
	w.window.CenterOnScreen()

	w.window.SetContent(w.getContent())
	w.window.Show()
}

func (w *ResolveManualWindow) getContent() fyne.CanvasObject {
	return container.NewScroll(
		container.NewVBox(w.getAllSteps()...),
	)
}

func (w *ResolveManualWindow) getAllSteps() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		w.getMediaStep(),
		w.getExportStep(),
		w.getSaveStep(),
		w.getLastStep(),
	}
}

func (w *ResolveManualWindow) getMediaStep() fyne.CanvasObject {
	text := getStepText(locale.ResolveManualMediaStep)
	image := getStepImage(resources.ResolveManualMedia)

	return container.NewVBox(text, image)
}

func (w *ResolveManualWindow) getExportStep() fyne.CanvasObject {
	text := getStepText(locale.ResolveManualExportStep)
	image := getStepImage(resources.ResolveManualExport)

	return container.NewVBox(text, image)
}

func (w *ResolveManualWindow) getSaveStep() fyne.CanvasObject {
	text := getStepText(locale.ResolveManualSaveStep)
	image := getStepImage(resources.ResolveManualSave)

	return container.NewVBox(text, image)
}

func (w *ResolveManualWindow) getLastStep() fyne.CanvasObject {
	text := getStepText(locale.ResolveManualLastStep)
	return container.NewVBox(text)
}

func getStepText(localeKey string) fyne.CanvasObject {
	text := widget.NewLabel(locale.Localize(localeKey))
	text.Wrapping = fyne.TextWrapWord

	return text
}

func getStepImage(resource fyne.Resource) fyne.CanvasObject {
	image := canvas.NewImageFromResource(resource)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(300, 300))

	return image
}

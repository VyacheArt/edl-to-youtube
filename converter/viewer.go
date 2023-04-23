package converter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VyacheArt/edl-to-youtube/caption"
	"github.com/VyacheArt/edl-to-youtube/edl"
	"strconv"
)

const (
	clipNumberIndex = iota
	//trackNumberIndex
	//transitionIndex
	//sourceInIndex
	//sourceOutIndex
	recordInIndex
	//recordOutIndex
)

type ViewerWindow struct {
	app     fyne.App
	window  fyne.Window
	edlList *edl.List

	caption binding.String

	generatorConfig caption.Config
	generator       *caption.Generator
}

func NewViewerWindow(app fyne.App, edlList *edl.List) *ViewerWindow {
	return &ViewerWindow{
		app:             app,
		edlList:         edlList,
		generator:       caption.NewGenerator(),
		generatorConfig: caption.DefaultConfig(),

		caption: binding.NewString(),
	}
}

func (w *ViewerWindow) Show() {
	w.window = w.app.NewWindow(w.edlList.Title)
	w.window.Resize(fyne.NewSize(900, 600))
	w.window.CenterOnScreen()

	w.window.SetContent(w.getContent())
	w.window.SetMainMenu(GetMainMenu(w.app, w.window))
	w.window.Show()

	w.regenerate()
}

func (w *ViewerWindow) getContent() fyne.CanvasObject {
	captionEntry := widget.NewMultiLineEntry()
	captionEntry.Bind(w.caption)

	content := container.NewHSplit(
		container.NewVSplit(
			widget.NewForm(
				widget.NewFormItem("Title", widget.NewLabel(w.edlList.Title)),
				widget.NewFormItem("Frame Code Mode", widget.NewLabel(frameCodeModeLabel(w.edlList.FrameCodeMode))),
				widget.NewFormItem("Timecode Format", widget.NewEntryWithData(w.bindString(&w.generatorConfig.TimeFormat))),
			),
			w.getTable(),
		),
		container.NewBorder(nil,
			container.NewBorder(
				nil, nil, nil,
				w.getRefreshButton(),
				w.getCopyButton(),
			),
			nil, nil, captionEntry),
	)

	content.SetOffset(0.7)

	return content
}

func (w *ViewerWindow) getCopyButton() *widget.Button {
	var copyButton *widget.Button
	copyButton = widget.NewButton("Copy", func() {
		text, _ := w.caption.Get()
		w.window.Clipboard().SetContent(text)
		copyButton.SetText("Copied!")
	})

	return copyButton
}

func (w *ViewerWindow) getRefreshButton() *widget.Button {
	return widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		w.regenerate()
	})
}

func (w *ViewerWindow) getTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			return len(w.edlList.Clips) + 1, 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("........................")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row == 0 {
				label.TextStyle.Bold = true
				label.SetText(w.getColumnTitles(i.Col))
				return
			}

			clip := w.edlList.Clips[i.Row-1]
			text := ""

			switch i.Col {
			case clipNumberIndex:
				text = strconv.Itoa(clip.ClipNumber)

			case recordInIndex:
				text = clip.RecordIn.String()
			}

			label.TextStyle.Bold = false
			label.SetText(text)
		},
	)

	return table
}

func (w *ViewerWindow) getColumnTitles(index int) string {
	switch index {
	case clipNumberIndex:
		return "Clip number"

	case recordInIndex:
		return "Timecode"
	}

	return ""
}

func (w *ViewerWindow) regenerate() {
	_ = w.caption.Set(w.generator.Generate(w.generatorConfig, w.edlList))
}

func (w *ViewerWindow) bindString(v *string) binding.ExternalString {
	b := binding.BindString(v)
	b.AddListener(binding.NewDataListener(func() {
		w.regenerate()
	}))

	return b
}

func frameCodeModeLabel(mode edl.FrameCodeMode) string {
	switch mode {
	case edl.FrameCodeModeDropFrame:
		return "Drop frame"
	case edl.FrameCodeModeNonDropFrame:
		return "Non drop frame"
	default:
		return "Unknown (" + string(mode) + ")"
	}
}

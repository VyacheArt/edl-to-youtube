/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package converter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VyacheArt/edl-to-youtube/caption"
	"github.com/VyacheArt/edl-to-youtube/converter/locale"
	widget2 "github.com/VyacheArt/edl-to-youtube/converter/widget"
	"github.com/VyacheArt/edl-to-youtube/edl"
	"strconv"
	"strings"
	"time"
)

// column indexes. Order is important and could be changed right here
const (
	clipNumberIndex = iota
	//trackNumberIndex
	//transitionIndex
	//sourceInIndex
	//sourceOutIndex
	recordInIndex
	//recordOutIndex
	colorIndex
	markerIndex
)

type ViewerWindow struct {
	app     *Application
	window  fyne.Window
	edlList *edl.List

	caption *widget.Entry

	generatorConfig caption.Config
	generator       *caption.Generator
}

func NewViewerWindow(app *Application, edlList *edl.List) *ViewerWindow {
	window := &ViewerWindow{
		app:             app,
		edlList:         edlList,
		generator:       caption.NewGenerator(),
		generatorConfig: caption.DefaultConfig(),
	}

	window.fillColors()
	return window
}

func (w *ViewerWindow) fillColors() {
	for _, clip := range w.edlList.Clips {
		w.generatorConfig.Colors[clip.Color] = true
	}
}

func (w *ViewerWindow) Show() {
	w.window = w.app.getApp().NewWindow(w.edlList.Title)
	w.window.Resize(fyne.NewSize(900, 600))
	w.window.CenterOnScreen()

	w.window.SetContent(w.getContent())
	w.window.Show()

	w.regenerate()
}

func (w *ViewerWindow) getContent() fyne.CanvasObject {
	w.caption = widget.NewMultiLineEntry()
	w.caption.TextStyle.Monospace = true

	content := container.NewHSplit(
		//left part with form and clips list
		container.NewVSplit(
			container.NewPadded(w.getForm()),
			w.getTable(),
		),

		//right part with generated caption
		container.NewBorder(nil,
			container.NewBorder(
				nil, nil, nil,
				w.getRefreshButton(),
				w.getCopyButton(),
			),
			nil, nil, w.caption),
	)

	content.SetOffset(0.7)

	return content
}

func (w *ViewerWindow) getForm() *widget.Form {
	colorFormItems := make([]*widget.FormItem, 0, len(w.generatorConfig.Colors))
	for c := range w.generatorConfig.Colors {
		color := c
		check := widget.NewCheck("", func(checked bool) {
			w.generatorConfig.Colors[color] = checked
			w.regenerate()
		})
		check.SetChecked(w.generatorConfig.Colors[color])

		colorFormItems = append(colorFormItems, widget.NewFormItem(colorLabel(color), check))
	}

	form := widget.NewForm(
		widget.NewFormItem(locale.Localize(locale.Title), widget.NewLabel(w.edlList.Title)),
		widget.NewFormItem(locale.Localize(locale.TimecodeFormat), widget.NewEntryWithData(w.bindString(&w.generatorConfig.TimeFormat))),
		widget.NewFormItem(locale.Localize(locale.IntroductionLabelTitle), widget.NewEntryWithData(w.bindString(&w.generatorConfig.IntroductionLabel))),
		widget.NewFormItem(locale.Localize(locale.OffsetSeconds), widget2.NewNumericalEntryWithData(w.bindInt(&w.generatorConfig.OffsetSeconds))),
	)

	for _, color := range colorFormItems {
		form.AppendItem(color)
	}

	return form
}

func (w *ViewerWindow) getCopyButton() *widget.Button {
	copyButtonText := locale.Localize(locale.Copy)
	textCopied := locale.Localize(locale.Copied)

	var copyButton *widget.Button
	copyButton = widget.NewButton(copyButtonText, func() {
		w.window.Clipboard().SetContent(w.caption.Text)
		copyButton.SetText(textCopied)

		go func() {
			<-time.After(1 * time.Second)
			copyButton.SetText(copyButtonText)
		}()
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
			return len(w.edlList.Clips) + 1, 4
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

			case colorIndex:
				text = colorLabel(clip.Color)

			case markerIndex:
				text = clip.Marker
			}

			label.TextStyle.Bold = false
			label.SetText(text)
		},
	)

	return table
}

// getColumnTitles returns title for column which will be rendered as a first row in the table
func (w *ViewerWindow) getColumnTitles(index int) string {
	switch index {
	case clipNumberIndex:
		return locale.Localize(locale.ClipNumber)

	case recordInIndex:
		return locale.Localize(locale.Timecode)

	case colorIndex:
		return locale.Localize(locale.Color)

	case markerIndex:
		return locale.Localize(locale.Marker)
	}

	return ""
}

// bindString is needed to make sure that caption will be regenerated when some field is changed
func (w *ViewerWindow) bindString(v *string) binding.ExternalString {
	b := binding.BindString(v)
	b.AddListener(binding.NewDataListener(func() {
		w.regenerate()
	}))

	return b
}

func (w *ViewerWindow) bindInt(v *int) binding.ExternalInt {
	b := binding.BindInt(v)
	b.AddListener(binding.NewDataListener(func() {
		w.regenerate()
	}))

	return b
}

func (w *ViewerWindow) regenerate() {
	w.caption.Text = w.generator.Generate(w.generatorConfig, w.edlList)
	w.caption.Refresh()
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

func colorLabel(rawColor string) string {
	//TODO: check adobe premiere color names
	color := strings.TrimPrefix(rawColor, "ResolveColor")
	return color
}

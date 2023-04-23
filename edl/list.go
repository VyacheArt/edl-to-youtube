package edl

import (
	"errors"
	"strings"
)

const (
	// FrameCodeModeDropFrame is the drop frame mode.
	FrameCodeModeDropFrame FrameCodeMode = "DF"

	// FrameCodeModeNonDropFrame is the non drop frame mode.
	FrameCodeModeNonDropFrame FrameCodeMode = "NDF"
)

type (
	List struct {
		Title         string
		FrameCodeMode FrameCodeMode
		Clips         []Clip
	}

	FrameCodeMode string
)

var ErrEmptyList = errors.New("empty decision list")

func Parse(s string) (*List, error) {
	l := List{}
	if err := l.Parse(s); err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *List) Parse(s string) error {
	//EDL is in the format of
	//TITLE: <title>
	//FCM: <frame code mode>
	//
	//<clips>

	lines := strings.Split(s, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		//skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		//parse the title
		if l.parseTitle(line) {
			continue
		}

		//parse the frame code mode
		if l.parseFrameCodeMode(line) {
			continue
		}

		clip := Clip{}

		//parse the clip header
		if err := clip.Parse(line); err != nil {
			continue
		}

		if len(lines) > i+1 {
			_ = clip.parseMeta(lines[i+1])
			i++
		}

		l.Clips = append(l.Clips, clip)
	}

	if l.IsEmpty() {
		return ErrEmptyList
	}

	return nil
}

func (l *List) parseTitle(line string) bool {
	const prefix = "TITLE:"
	if !strings.HasPrefix(line, prefix) {
		return false
	}

	l.Title = strings.TrimSpace(line[len(prefix):])
	return true
}

func (l *List) parseFrameCodeMode(line string) bool {
	const prefix = "FCM:"
	if !strings.HasPrefix(line, prefix) {
		return false
	}

	rawFrameCodeMode := strings.TrimSpace(line[len(prefix):])
	switch rawFrameCodeMode {
	case "NON-DROP FRAME", "NDF":
		l.FrameCodeMode = FrameCodeModeNonDropFrame

	case "DROP FRAME", "DF":
		l.FrameCodeMode = FrameCodeModeDropFrame
	}

	return true
}

func (l *List) IsEmpty() bool {
	return len(l.Title) == 0 && len(l.FrameCodeMode) == 0 && len(l.Clips) == 0
}

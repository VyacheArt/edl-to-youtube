package caption

import (
	"github.com/VyacheArt/edl-to-youtube/converter/locale"
	"github.com/VyacheArt/edl-to-youtube/edl"
	"strings"
	"time"
)

const introductionLabelKey = "introductionLabel"

type (
	Config struct {
		TimeFormat        string
		AutoIntroduction  bool
		IntroductionLabel string
		OffsetSeconds     int
		Colors            map[string]bool
	}

	Generator struct{}
)

func DefaultConfig() Config {
	return Config{
		TimeFormat:        "15:04:05",
		AutoIntroduction:  true,
		IntroductionLabel: locale.Localize(introductionLabelKey),
		OffsetSeconds:     0,
		Colors:            make(map[string]bool),
	}
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(cfg Config, list *edl.List) string {
	if len(list.Clips) == 0 {
		return ""
	}

	timeFormat := cfg.TimeFormat
	if len(timeFormat) == 0 {
		timeFormat = "15:04:05"
	}

	builder := strings.Builder{}

	if cfg.AutoIntroduction {
		first := list.Clips[0]
		firstTime := first.RecordIn.Time

		//check if hours, minutes and seconds are more than zero
		if firstTime.Hour() > 0 || firstTime.Minute() > 0 || firstTime.Second() > 0 {
			builder.WriteString(time.Unix(0, 0).In(time.UTC).Format(timeFormat))
			builder.WriteString(" ")
			builder.WriteString(cfg.IntroductionLabel)
			builder.WriteString("\n")
		}
	}

	for _, clip := range list.Clips {
		if enabled, ok := cfg.Colors[clip.Color]; !ok || !enabled {
			continue
		}

		recordIn := clip.RecordIn.Time.Add(time.Duration(cfg.OffsetSeconds) * time.Second)

		builder.WriteString(recordIn.Format(timeFormat))
		builder.WriteString(" ")
		builder.WriteString(clip.Marker)
		builder.WriteString("\n")
	}

	return builder.String()
}

package caption

import (
	"github.com/VyacheArt/edl-to-youtube/edl"
	"strings"
)

type (
	Config struct {
		TimeFormat string
	}

	Generator struct{}
)

func DefaultConfig() Config {
	return Config{
		TimeFormat: "15:04:05",
	}
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(cfg Config, list *edl.List) string {
	timeFormat := cfg.TimeFormat
	if len(timeFormat) == 0 {
		timeFormat = "15:04:05"
	}

	builder := strings.Builder{}

	for _, clip := range list.Clips {
		builder.WriteString(clip.RecordIn.Time.Format(timeFormat))
		builder.WriteString(" ")
		builder.WriteString(clip.Marker)
		builder.WriteString("\n")
	}

	return builder.String()
}

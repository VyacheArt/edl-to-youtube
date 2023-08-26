package locale

import (
	"embed"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"unicode"
)

const DefaultLocale = "en"

//goland:noinspection GoCommentStart
const (
	// Menu
	MenuLanguage      = "language"
	MenuChooseFile    = "chooseFile"
	MenuFile          = "file"
	MenuOpen          = "open"
	MenuOpenLast      = "openLast"
	MenuWelcomeScreen = "welcomeScreen"
	MenuHelp          = "help"
	MenuQuit          = "quit"

	// Help
	HelpHowTo           = "helpHowTo"
	HelpExportInResolve = "helpDavinciExport"

	// DaVinci Resolve Manual
	ResolveManualMediaStep  = "resolveManualMediaStep"
	ResolveManualExportStep = "resolveManualExportStep"
	ResolveManualSaveStep   = "resolveManualSaveStep"
	ResolveManualLastStep   = "resolveManualLastStep"
)

var (
	bundle    = i18n.NewBundle(language.English)
	localizer = i18n.NewLocalizer(bundle, DefaultLocale)

	//go:embed *.toml
	localeFS embed.FS
)

func init() {
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load all locale files
	locales, _ := localeFS.ReadDir(".")
	for _, locale := range locales {
		_, err := bundle.LoadMessageFileFS(localeFS, locale.Name())
		if err != nil {
			panic(err)
		}
	}
}

func GetLocales() map[string]string {
	res := make(map[string]string, len(bundle.LanguageTags()))
	for _, tag := range bundle.LanguageTags() {
		title := display.Self.Name(tag)
		title = capitalize(title)

		res[tag.String()] = title
	}

	return res
}

func SetLocale(lang string) {
	localizer = i18n.NewLocalizer(bundle, lang)
}

func Localize(key string, args ...interface{}) string {
	res, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	})

	if err != nil {
		return key
	}

	return res
}

func capitalize(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

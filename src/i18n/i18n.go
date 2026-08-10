// stock_ticker
// src/i18n/i18n.go

// Package i18n loads the translated string catalogs and resolves the active
// language. The set of supported languages is exactly the set of locale files
// present under locales/ (one <tag>.toml per language, e.g. en.toml,
// fr-CA.toml). Look up any string by its message ID with T.
package i18n

import (
	"embed"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jeandeaual/go-locale"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// locales/*.toml are embedded into the binary, so no external files are needed
// at runtime. Adding a language is as simple as dropping a new <tag>.toml here;
// it becomes a supported language automatically.
//
//go:embed locales/*.toml
var localeFS embed.FS

// defaultLang is the base language of the bundle and the final fallback when a
// message is missing in the active language. It must have a matching catalog
// file (locales/en.toml). Requests for regional variants (en-US, en-GB) match
// it automatically.
var defaultLang = language.English

var (
	bundle        *goi18n.Bundle
	localizer     *goi18n.Localizer
	supportedTags []language.Tag // index 0 is defaultLang
	matcher       language.Matcher
)

func init() {
	bundle = goi18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	entries, err := localeFS.ReadDir("locales")
	if err != nil {
		panic("i18n: cannot read embedded locales: " + err.Error())
	}

	// supportedTags is built from the catalog files. defaultLang is forced to
	// the front so the matcher uses it as the fallback for unrecognized input.
	hasDefault := false
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		if _, err := bundle.LoadMessageFileFS(localeFS, "locales/"+e.Name()); err != nil {
			panic("i18n: cannot load " + e.Name() + ": " + err.Error())
		}
		tag, err := language.Parse(strings.TrimSuffix(e.Name(), ".toml"))
		if err != nil {
			panic("i18n: locale file " + e.Name() + " is not a valid language tag: " + err.Error())
		}
		if tag == defaultLang {
			hasDefault = true
			continue
		}
		supportedTags = append(supportedTags, tag)
	}
	if !hasDefault {
		panic("i18n: missing base catalog locales/" + defaultLang.String() + ".toml")
	}
	supportedTags = append([]language.Tag{defaultLang}, supportedTags...)
	matcher = language.NewMatcher(supportedTags)

	// Default to the OS locale until SetLanguage is called explicitly.
	_ = SetLanguage("")
}

// SupportedLanguages returns the BCP-47 tags of every language that has a
// catalog, e.g. []string{"en-US", "fr-CA"}.
func SupportedLanguages() []string {
	tags := make([]string, len(supportedTags))
	for i, t := range supportedTags {
		tags[i] = t.String()
	}
	return tags
}

// SetLanguage selects the active language.
//
//   - lang == ""  : auto-detect from the OS locale, matching to the closest
//     supported language and falling back to the default. Never errors.
//   - lang != ""  : force that language. It must be one of SupportedLanguages,
//     otherwise an error is returned and the active language is left unchanged.
func SetLanguage(lang string) error {
	if lang == "" {
		detected, _ := locale.GetLocale() // "" on failure is fine; matcher falls back
		tag, _, _ := matcher.Match(language.Make(detected))
		localizer = goi18n.NewLocalizer(bundle, tag.String())
		return nil
	}

	tag, err := language.Parse(lang)
	if err != nil {
		return fmt.Errorf("i18n: %q is not a valid language tag", lang)
	}
	if _, _, conf := matcher.Match(tag); conf == language.No {
		return fmt.Errorf("i18n: unsupported language %q (supported: %s)",
			lang, strings.Join(SupportedLanguages(), ", "))
	}
	localizer = goi18n.NewLocalizer(bundle, tag.String(), defaultLang.String())
	return nil
}

// T returns the translated string for the given message ID in the active
// language. If the ID is unknown it returns the ID itself, so a missing
// translation is visible rather than silently blank.
func T(messageID string) string {
	s, err := localizer.Localize(&goi18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		return messageID
	}
	return s
}

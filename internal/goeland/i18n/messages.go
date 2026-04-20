package i18n

import (
	"embed"
	"io/fs"

	"github.com/slurdge/goeland/log"
	"github.com/slurdge/goeland/version"
	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/localize"
	"golang.org/x/text/language"
)

//go:embed locale/*
var localeFS embed.FS

var bundle *spreak.Bundle
var localizer *spreak.Localizer

func Init(locale string) {
	var err error
	subSystem, _ := fs.Sub(localeFS, "locale")
	bundle, err = spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain(version.ProductName),
		spreak.WithDomainFs(version.ProductName, subSystem),
		spreak.WithLanguage(language.English, language.French),
	)

	if err != nil {
		log.Debugf("Error creating language bundle for %v", locale)
		bundle, _ = spreak.NewBundle()
	}
	localizer = spreak.NewLocalizer(bundle, locale)

}

func Language() language.Tag {
	return localizer.Language()
}

func T(message localize.Singular) string {
	return localizer.Get(message)
}

func Tf(message localize.Singular, args ...any) string {
	return localizer.Getf(message, args...)
}

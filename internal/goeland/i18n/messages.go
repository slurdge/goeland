package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	message.SetString(language.BritishEnglish, "Digest for %s", "Digest for %s")
	message.SetString(language.AmericanEnglish, "Digest for %s", "Digest for %s")
	message.SetString(language.French, "Digest for %s", "Abrégé de %s")
}

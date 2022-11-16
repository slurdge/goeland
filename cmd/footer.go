package cmd

import (
	"fmt"

	"golang.org/x/text/language"
)

const goelandURL = `<a href="https://www.github.com/slurdge/goeland">goeland</a>`

var footersI8n = map[language.Tag][]string{
	language.BritishEnglish: {
		`Sent with ❤️ by %s`,
		`Sent with 💖 by %s`,
		`Sent with 💙 by %s`,
		`Sent with 🥰 by %s`,
		`Brought to you 🐣 by %s`,
		`Sent quickly ⚡ by %s`,
		`Sent with a touch of 💐 by %s`,
		`Sent with a touch of 🌸 by %s`,
		`Sent with a touch of 🌼 by %s`,
		`In your 📧 from %s`,
		`Your 📰 from %s`,
		`Sent bravely ⚔️ by %s`,
		`Sent happily 😊 by %s`,
		`Sent smoothly ✈️ by %s`,
		`Sent simply 🌞 by %s`,
		`Sent with 🤘🏻 by %s`,
		`Sent 💨 faster than a carrier pigeon by %s`,
		`Delicious mail 🤤 by %s`,
		`Fresh mail 🐟 by %s`,
		`Fresh out of the oven 🥐 by %s`,
		`Piping-🔥 mail by %s`,
		`A good 📧 for a good day by %s`,
		`Enjoy your 📧 by %s`,
		`Dropped softly 🕊️ by %s`,
		`Delivered on ⏰ by %s`,
	},
	language.French: {
		`Envoyé avec ❤️ par %s`,
		`Envoyé avec 💖 par %s`,
		`Envoyé avec 💙 par %s`,
		`Envoyé avec 🥰 par %s`,
		`Amené tout 🐣 par %s`,
		`Envoyé rapidemment ⚡ par %s`,
		`Envoyé avec a touch of 💐 par %s`,
		`Envoyé avec a touch of 🌸 par %s`,
		`Envoyé avec a touch of 🌼 par %s`,
		`Dans votre 📧 depuis %s`,
		`Vos 📰 de la part de %s`,
		`Envoyé bravement ⚔️ par %s`,
		`Envoyé avec joie 😊 par %s`,
		`Envoyé par avion ✈️ par %s`,
		`Envoyé tout simplement 🌞 par %s`,
		`Envoyé avec 🤘🏻 par %s`,
		`Envoyé 💨 plus vite qu'un pigeon voyageur par %s`,
		`Un mail délicieux 🤤 par %s`,
		`Tout frais 🐟 par %s`,
		`Sorti directement du four 🥐 par %s`,
		`Un bon 📧 pour une bonne journée par %s`,
		`Délectez vous de votre 📧 par %s`,
		`Déposé doucement 🕊️ par %s`,
		`Reçu à l'⏰ par %s`,
	},
}

var footers []string

func init() {
	for _, footer := range footersI8n[language.BritishEnglish] {
		footers = append(footers, fmt.Sprintf(footer, goelandURL))
	}
}

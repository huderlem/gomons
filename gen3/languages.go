package gen3

// GameLanguage represents the language of a game catridge.
type GameLanguage byte

// Game languages enumeration
const (
	LanguageJapanese GameLanguage = 1
	LanguageEnglish  GameLanguage = 2
	LanguageFrench   GameLanguage = 3
	LanguageItalian  GameLanguage = 4
	LanguageGerman   GameLanguage = 5
	LanguageKorean   GameLanguage = 6
	LanguageSpanish  GameLanguage = 7
)

func (l GameLanguage) String() string {
	switch l {
	case LanguageJapanese:
		return "Japanese"
	case LanguageEnglish:
		return "English"
	case LanguageFrench:
		return "French"
	case LanguageItalian:
		return "Italian"
	case LanguageGerman:
		return "German"
	case LanguageKorean:
		return "Korean"
	case LanguageSpanish:
		return "Spanish"
	default:
		return "Invalid Language"
	}
}

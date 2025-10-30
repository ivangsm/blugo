package i18n

import (
	"os"
	"strings"
)

// Language representa un idioma soportado
type Language string

const (
	English Language = "en"
	Spanish Language = "es"
)

// Translations contiene todas las traducciones
type Translations struct {
	// App
	AppTitle          string
	Scanning          string
	Paused            string
	Initializing      string

	// Sections
	AvailableDevices  string
	ConnectedDevices  string
	AdapterInfo       string

	// Device info
	NoDevicesAvailable string
	NoDevicesConnected string

	// Actions
	Connecting        string
	Disconnecting     string
	Pairing           string
	Forgetting        string

	// Adapter
	AdapterPoweringOn  string
	AdapterPoweringOff string
	AdapterPoweredOn   string
	AdapterPoweredOff  string
	DiscoverableOn     string
	DiscoverableOff    string
	DiscoverableActivating   string
	DiscoverableDeactivating string
	PairableOn         string
	PairableOff        string
	PairableActivating   string
	PairableDeactivating string

	// Status messages
	ScanEnabled       string
	ScanPaused        string
	Connected         string
	Disconnected      string
	DisconnectedPaired string
	Forgotten         string

	// Errors
	Error             string
	ErrorScanToggle   string

	// Pairing
	PairingCode       string
	PairingInstruction string
	PairingConfirm    string
	PairingCancelled  string

	// Help
	HelpNavigation    string
	HelpActions       string
	HelpAdapterControl string
	HelpGeneral       string
	HelpPairing       string

	// Adapter table
	AdapterName       string
	AdapterAlias      string
	AdapterPower      string
	AdapterPairable   string
	AdapterDiscoverable string

	// Status
	StatusOn          string
	StatusOff         string

	// Badges
	BadgePaired       string
	BadgeConnected    string
	BadgeTrusted      string
}

var currentLang Language = English // Default language

// T holds current translations
var T *Translations

func init() {
	// Set default language (English)
	SetLanguage(English)
}

// InitFromConfig initializes the language based on a language code string
// Accepts "en", "es", or falls back to English for unknown languages
func InitFromConfig(langCode string) {
	switch Language(langCode) {
	case English:
		SetLanguage(English)
	case Spanish:
		SetLanguage(Spanish)
	default:
		SetLanguage(English)
	}
}

// detectSystemLanguage detects the system language from environment
// Returns English as fallback for unsupported languages
func detectSystemLanguage() Language {
	// Check LANGUAGE, LC_ALL, LC_MESSAGES, LANG in order
	for _, env := range []string{"LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"} {
		if val := os.Getenv(env); val != "" {
			// Extract language code (e.g., "en_US.UTF-8" -> "en")
			lang := strings.Split(val, "_")[0]
			lang = strings.Split(lang, ".")[0]
			lang = strings.ToLower(lang)

			switch lang {
			case "en":
				return English
			case "es":
				return Spanish
			// Add more languages here in the future:
			// case "fr":
			//     return French
			// case "de":
			//     return German
			default:
				// Default to English for unknown languages
				return English
			}
		}
	}

	return English // Fallback
}

// GetSupportedLanguages returns a list of all supported languages
func GetSupportedLanguages() []Language {
	return []Language{English, Spanish}
}

// SetLanguage sets the current language
func SetLanguage(lang Language) {
	currentLang = lang

	switch lang {
	case English:
		T = &englishTranslations
	case Spanish:
		T = &spanishTranslations
	default:
		T = &englishTranslations
	}
}

// GetCurrentLanguage returns the current language
func GetCurrentLanguage() Language {
	return currentLang
}

// ToggleLanguage switches between English and Spanish
func ToggleLanguage() {
	if currentLang == English {
		SetLanguage(Spanish)
	} else {
		SetLanguage(English)
	}
}

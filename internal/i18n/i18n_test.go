package i18n

import (
	"os"
	"testing"
)

func TestInitFromConfig(t *testing.T) {
	tests := []struct {
		name         string
		langCode     string
		expectedLang Language
	}{
		{
			name:         "sets English for 'en'",
			langCode:     "en",
			expectedLang: English,
		},
		{
			name:         "sets Spanish for 'es'",
			langCode:     "es",
			expectedLang: Spanish,
		},
		{
			name:         "falls back to English for unknown language",
			langCode:     "fr",
			expectedLang: English,
		},
		{
			name:         "falls back to English for empty string",
			langCode:     "",
			expectedLang: English,
		},
		{
			name:         "falls back to English for invalid code",
			langCode:     "invalid",
			expectedLang: English,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitFromConfig(tt.langCode)
			got := GetCurrentLanguage()
			if got != tt.expectedLang {
				t.Errorf("InitFromConfig(%v) resulted in language %v, want %v", tt.langCode, got, tt.expectedLang)
			}
		})
	}
}

func TestDetectSystemLanguage(t *testing.T) {
	// Save original environment variables
	originalEnvs := make(map[string]string)
	envVars := []string{"LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"}
	for _, env := range envVars {
		originalEnvs[env] = os.Getenv(env)
	}

	// Restore environment after test
	defer func() {
		for env, val := range originalEnvs {
			if val == "" {
				os.Unsetenv(env)
			} else {
				os.Setenv(env, val)
			}
		}
	}()

	tests := []struct {
		name         string
		envVars      map[string]string
		expectedLang Language
	}{
		{
			name: "detects English from LANGUAGE",
			envVars: map[string]string{
				"LANGUAGE": "en_US.UTF-8",
			},
			expectedLang: English,
		},
		{
			name: "detects Spanish from LANGUAGE",
			envVars: map[string]string{
				"LANGUAGE": "es_ES.UTF-8",
			},
			expectedLang: Spanish,
		},
		{
			name: "detects English from LC_ALL",
			envVars: map[string]string{
				"LANGUAGE": "",
				"LC_ALL":   "en_GB.UTF-8",
			},
			expectedLang: English,
		},
		{
			name: "detects Spanish from LC_MESSAGES",
			envVars: map[string]string{
				"LANGUAGE":    "",
				"LC_ALL":      "",
				"LC_MESSAGES": "es_MX.UTF-8",
			},
			expectedLang: Spanish,
		},
		{
			name: "detects English from LANG",
			envVars: map[string]string{
				"LANGUAGE":    "",
				"LC_ALL":      "",
				"LC_MESSAGES": "",
				"LANG":        "en_CA.UTF-8",
			},
			expectedLang: English,
		},
		{
			name: "handles simple language code without region",
			envVars: map[string]string{
				"LANGUAGE": "es",
			},
			expectedLang: Spanish,
		},
		{
			name: "handles language code with region but no encoding",
			envVars: map[string]string{
				"LANGUAGE": "en_US",
			},
			expectedLang: English,
		},
		{
			name: "falls back to English for unsupported language",
			envVars: map[string]string{
				"LANGUAGE": "fr_FR.UTF-8",
			},
			expectedLang: English,
		},
		{
			name: "falls back to English when no env vars set",
			envVars: map[string]string{
				"LANGUAGE":    "",
				"LC_ALL":      "",
				"LC_MESSAGES": "",
				"LANG":        "",
			},
			expectedLang: English,
		},
		{
			name: "prioritizes LANGUAGE over other env vars",
			envVars: map[string]string{
				"LANGUAGE":    "es_ES.UTF-8",
				"LC_ALL":      "en_US.UTF-8",
				"LC_MESSAGES": "en_US.UTF-8",
				"LANG":        "en_US.UTF-8",
			},
			expectedLang: Spanish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables first
			for _, env := range envVars {
				os.Unsetenv(env)
			}

			// Set test environment variables
			for env, val := range tt.envVars {
				if val != "" {
					os.Setenv(env, val)
				}
			}

			got := detectSystemLanguage()
			if got != tt.expectedLang {
				t.Errorf("detectSystemLanguage() = %v, want %v", got, tt.expectedLang)
			}
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	languages := GetSupportedLanguages()

	if len(languages) != 2 {
		t.Errorf("GetSupportedLanguages() returned %d languages, want 2", len(languages))
	}

	// Check that English and Spanish are in the list
	hasEnglish := false
	hasSpanish := false

	for _, lang := range languages {
		if lang == English {
			hasEnglish = true
		}
		if lang == Spanish {
			hasSpanish = true
		}
	}

	if !hasEnglish {
		t.Errorf("GetSupportedLanguages() missing English")
	}
	if !hasSpanish {
		t.Errorf("GetSupportedLanguages() missing Spanish")
	}
}

func TestSetLanguage(t *testing.T) {
	tests := []struct {
		name         string
		lang         Language
		expectedLang Language
	}{
		{
			name:         "sets English",
			lang:         English,
			expectedLang: English,
		},
		{
			name:         "sets Spanish",
			lang:         Spanish,
			expectedLang: Spanish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLanguage(tt.lang)
			got := GetCurrentLanguage()
			if got != tt.expectedLang {
				t.Errorf("After SetLanguage(%v), GetCurrentLanguage() = %v, want %v", tt.lang, got, tt.expectedLang)
			}

			// Verify T is set correctly
			if T == nil {
				t.Errorf("SetLanguage(%v) did not set T", tt.lang)
			}
		})
	}
}

func TestSetLanguage_SetsTranslations(t *testing.T) {
	// Test that T is properly set when switching languages
	SetLanguage(English)
	if T == nil {
		t.Fatal("T should not be nil after SetLanguage(English)")
	}
	englishTitle := T.AppTitle

	SetLanguage(Spanish)
	if T == nil {
		t.Fatal("T should not be nil after SetLanguage(Spanish)")
	}
	spanishTitle := T.AppTitle

	// The titles should be different (one in English, one in Spanish)
	if englishTitle == spanishTitle {
		t.Errorf("English and Spanish translations should be different")
	}

	// Verify they're not empty
	if englishTitle == "" {
		t.Errorf("English AppTitle should not be empty")
	}
	if spanishTitle == "" {
		t.Errorf("Spanish AppTitle should not be empty")
	}
}

func TestGetCurrentLanguage(t *testing.T) {
	// Set to English and verify
	SetLanguage(English)
	if got := GetCurrentLanguage(); got != English {
		t.Errorf("After SetLanguage(English), GetCurrentLanguage() = %v, want %v", got, English)
	}

	// Set to Spanish and verify
	SetLanguage(Spanish)
	if got := GetCurrentLanguage(); got != Spanish {
		t.Errorf("After SetLanguage(Spanish), GetCurrentLanguage() = %v, want %v", got, Spanish)
	}
}

func TestToggleLanguage(t *testing.T) {
	tests := []struct {
		name         string
		initialLang  Language
		expectedLang Language
	}{
		{
			name:         "toggles from English to Spanish",
			initialLang:  English,
			expectedLang: Spanish,
		},
		{
			name:         "toggles from Spanish to English",
			initialLang:  Spanish,
			expectedLang: English,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLanguage(tt.initialLang)
			ToggleLanguage()
			got := GetCurrentLanguage()
			if got != tt.expectedLang {
				t.Errorf("After ToggleLanguage() from %v, got %v, want %v", tt.initialLang, got, tt.expectedLang)
			}
		})
	}
}

func TestToggleLanguage_Multiple(t *testing.T) {
	// Start with English
	SetLanguage(English)

	// First toggle should give Spanish
	ToggleLanguage()
	if got := GetCurrentLanguage(); got != Spanish {
		t.Errorf("First toggle: got %v, want Spanish", got)
	}

	// Second toggle should give English again
	ToggleLanguage()
	if got := GetCurrentLanguage(); got != English {
		t.Errorf("Second toggle: got %v, want English", got)
	}

	// Third toggle should give Spanish
	ToggleLanguage()
	if got := GetCurrentLanguage(); got != Spanish {
		t.Errorf("Third toggle: got %v, want Spanish", got)
	}
}

// TestLanguageConstants verifies the language constants are set correctly
func TestLanguageConstants(t *testing.T) {
	if English != "en" {
		t.Errorf("English constant = %v, want 'en'", English)
	}
	if Spanish != "es" {
		t.Errorf("Spanish constant = %v, want 'es'", Spanish)
	}
}

// TestInit verifies the init function sets a default language
func TestInit(t *testing.T) {
	// The init() function should have already run
	// We can't reliably test currentLang here because other tests may have changed it
	// But we can verify that T is initialized
	if T == nil {
		t.Errorf("After init, T should not be nil")
	}

	// Verify the package is functional by testing language switching
	SetLanguage(English)
	if currentLang != English {
		t.Errorf("SetLanguage(English) failed, currentLang = %v", currentLang)
	}

	SetLanguage(Spanish)
	if currentLang != Spanish {
		t.Errorf("SetLanguage(Spanish) failed, currentLang = %v", currentLang)
	}
}

// TestIntegration_LanguageSwitching tests a complete language switching workflow
func TestIntegration_LanguageSwitching(t *testing.T) {
	// Start with default (English)
	if GetCurrentLanguage() != English {
		SetLanguage(English)
	}

	// Load from config "es"
	InitFromConfig("es")
	if GetCurrentLanguage() != Spanish {
		t.Errorf("InitFromConfig('es') did not set Spanish")
	}

	// Toggle back to English
	ToggleLanguage()
	if GetCurrentLanguage() != English {
		t.Errorf("ToggleLanguage() did not switch to English")
	}

	// Set directly to Spanish
	SetLanguage(Spanish)
	if GetCurrentLanguage() != Spanish {
		t.Errorf("SetLanguage(Spanish) did not set Spanish")
	}

	// Verify translations are loaded
	if T == nil {
		t.Errorf("Translations should be loaded")
	}
	if T.AppTitle == "" {
		t.Errorf("AppTitle should not be empty")
	}
}

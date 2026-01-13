// Package services contains the application services (use cases).
package services

import (
	"fmt"
	"strings"
	"time"
)

// Locale represents a supported language/region combination.
type Locale string

const (
	LocaleEnUS Locale = "en-US"
	LocalePtBR Locale = "pt-BR"
	LocaleEsES Locale = "es-ES"
	LocaleFrFR Locale = "fr-FR"
	LocaleDeDE Locale = "de-DE"
)

// TranslationKey represents a key for translated text.
type TranslationKey string

// Common translation keys used in resume templates.
const (
	KeyProfessionalSummary TranslationKey = "professional_summary"
	KeyEducation           TranslationKey = "education"
	KeyExperience          TranslationKey = "experience"
	KeyProjects            TranslationKey = "projects"
	KeyTechnicalSkills     TranslationKey = "technical_skills"
	KeyLanguages           TranslationKey = "languages"
	KeyPresent             TranslationKey = "present"
	KeyGPA                 TranslationKey = "gpa"
	KeyGrade               TranslationKey = "grade"
	KeyNative              TranslationKey = "native"
	KeyFluent              TranslationKey = "fluent"
	KeyAdvanced            TranslationKey = "advanced"
	KeyIntermediate        TranslationKey = "intermediate"
	KeyBasic               TranslationKey = "basic"
)

// translations contains all localized strings.
var translations = map[Locale]map[TranslationKey]string{
	LocaleEnUS: {
		KeyProfessionalSummary: "Professional Summary",
		KeyEducation:           "Education",
		KeyExperience:          "Experience",
		KeyProjects:            "Projects",
		KeyTechnicalSkills:     "Technical Skills",
		KeyLanguages:           "Languages",
		KeyPresent:             "Present",
		KeyGPA:                 "GPA",
		KeyGrade:               "Grade",
		KeyNative:              "Native",
		KeyFluent:              "Fluent",
		KeyAdvanced:            "Advanced",
		KeyIntermediate:        "Intermediate",
		KeyBasic:               "Basic",
	},
	LocalePtBR: {
		KeyProfessionalSummary: "Resumo Profissional",
		KeyEducation:           "Formação Acadêmica",
		KeyExperience:          "Experiência Profissional",
		KeyProjects:            "Projetos",
		KeyTechnicalSkills:     "Habilidades Técnicas",
		KeyLanguages:           "Idiomas",
		KeyPresent:             "Atual",
		KeyGPA:                 "CR",
		KeyGrade:               "Média",
		KeyNative:              "Nativo",
		KeyFluent:              "Fluente",
		KeyAdvanced:            "Avançado",
		KeyIntermediate:        "Intermediário",
		KeyBasic:               "Básico",
	},
	LocaleEsES: {
		KeyProfessionalSummary: "Resumen Profesional",
		KeyEducation:           "Formación Académica",
		KeyExperience:          "Experiencia Profesional",
		KeyProjects:            "Proyectos",
		KeyTechnicalSkills:     "Habilidades Técnicas",
		KeyLanguages:           "Idiomas",
		KeyPresent:             "Actual",
		KeyGPA:                 "Promedio",
		KeyGrade:               "Nota",
		KeyNative:              "Nativo",
		KeyFluent:              "Fluido",
		KeyAdvanced:            "Avanzado",
		KeyIntermediate:        "Intermedio",
		KeyBasic:               "Básico",
	},
	LocaleFrFR: {
		KeyProfessionalSummary: "Résumé Professionnel",
		KeyEducation:           "Formation",
		KeyExperience:          "Expérience Professionnelle",
		KeyProjects:            "Projets",
		KeyTechnicalSkills:     "Compétences Techniques",
		KeyLanguages:           "Langues",
		KeyPresent:             "Présent",
		KeyGPA:                 "Moyenne",
		KeyGrade:               "Note",
		KeyNative:              "Natif",
		KeyFluent:              "Courant",
		KeyAdvanced:            "Avancé",
		KeyIntermediate:        "Intermédiaire",
		KeyBasic:               "Basique",
	},
	LocaleDeDE: {
		KeyProfessionalSummary: "Berufsprofil",
		KeyEducation:           "Ausbildung",
		KeyExperience:          "Berufserfahrung",
		KeyProjects:            "Projekte",
		KeyTechnicalSkills:     "Technische Fähigkeiten",
		KeyLanguages:           "Sprachen",
		KeyPresent:             "Aktuell",
		KeyGPA:                 "Notendurchschnitt",
		KeyGrade:               "Note",
		KeyNative:              "Muttersprache",
		KeyFluent:              "Fließend",
		KeyAdvanced:            "Fortgeschritten",
		KeyIntermediate:        "Mittelstufe",
		KeyBasic:               "Grundkenntnisse",
	},
}

// monthNames contains localized month names.
var monthNames = map[Locale][]string{
	LocaleEnUS: {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	LocalePtBR: {"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"},
	LocaleEsES: {"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"},
	LocaleFrFR: {"Jan", "Fév", "Mar", "Avr", "Mai", "Juin", "Juil", "Août", "Sep", "Oct", "Nov", "Déc"},
	LocaleDeDE: {"Jan", "Feb", "Mär", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"},
}

// I18n provides internationalization utilities for resume generation.
type I18n struct {
	locale Locale
}

// NewI18n creates a new I18n instance for the specified locale.
func NewI18n(locale Locale) *I18n {
	// Default to en-US if locale is not supported.
	if _, ok := translations[locale]; !ok {
		locale = LocaleEnUS
	}
	return &I18n{locale: locale}
}

// ParseLocale converts a locale string to a Locale type.
// Supports formats like "en", "en-US", "pt-br", "pt_BR".
func ParseLocale(s string) Locale {
	s = strings.ToLower(strings.ReplaceAll(s, "_", "-"))

	switch {
	case strings.HasPrefix(s, "pt"):
		return LocalePtBR
	case strings.HasPrefix(s, "es"):
		return LocaleEsES
	case strings.HasPrefix(s, "fr"):
		return LocaleFrFR
	case strings.HasPrefix(s, "de"):
		return LocaleDeDE
	default:
		return LocaleEnUS
	}
}

// T returns the translated string for the given key.
func (i *I18n) T(key TranslationKey) string {
	if dict, ok := translations[i.locale]; ok {
		if text, ok := dict[key]; ok {
			return text
		}
	}
	// Fallback to en-US.
	if dict, ok := translations[LocaleEnUS]; ok {
		if text, ok := dict[key]; ok {
			return text
		}
	}
	// If all else fails, return the key as-is.
	return string(key)
}

// Locale returns the current locale.
func (i *I18n) Locale() Locale {
	return i.locale
}

// FormatDate formats a date according to the locale.
// For en-US: "Jan 2024"
// For pt-BR: "01/2024"
// For es-ES: "Ene 2024"
func (i *I18n) FormatDate(t time.Time) string {
	month := int(t.Month())
	year := t.Year()

	switch i.locale {
	case LocalePtBR:
		// Portuguese uses numeric format: MM/YYYY
		return fmt.Sprintf("%02d/%d", month, year)
	default:
		// Other locales use abbreviated month name
		months := monthNames[i.locale]
		if months == nil {
			months = monthNames[LocaleEnUS]
		}
		return fmt.Sprintf("%s %d", months[month-1], year)
	}
}

// FormatDateString parses a date string and formats it according to the locale.
// Accepts formats: "2024-01-15", "2024-01", "Jan 2024"
func (i *I18n) FormatDateString(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Try parsing as ISO date (YYYY-MM-DD).
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return i.FormatDate(t)
	}

	// Try parsing as YYYY-MM.
	if t, err := time.Parse("2006-01", dateStr); err == nil {
		return i.FormatDate(t)
	}

	// Try parsing as "Jan 2006".
	if t, err := time.Parse("Jan 2006", dateStr); err == nil {
		return i.FormatDate(t)
	}

	// Return original if we can't parse.
	return dateStr
}

// FormatDateRange formats a date range according to the locale.
// If endDate is nil or empty, uses "Present" (localized).
func (i *I18n) FormatDateRange(startDate string, endDate *string) string {
	start := i.FormatDateString(startDate)

	if endDate == nil || *endDate == "" {
		return fmt.Sprintf("%s – %s", start, i.T(KeyPresent))
	}

	end := i.FormatDateString(*endDate)
	return fmt.Sprintf("%s – %s", start, end)
}

// FormatGPA formats GPA display according to the locale.
// For en-US: "GPA: 3.8"
// For pt-BR: "CR: 3.8" or "Média: 8.5"
func (i *I18n) FormatGPA(gpa float64, scale float64) string {
	label := i.T(KeyGPA)

	// Format based on typical scale.
	if scale == 10 {
		// Brazilian 0-10 scale.
		return fmt.Sprintf("%s: %.1f", label, gpa)
	}
	// American 4.0 scale.
	return fmt.Sprintf("%s: %.2f", label, gpa)
}

// FormatProficiencyLevel returns the localized proficiency level.
func (i *I18n) FormatProficiencyLevel(level string) string {
	switch strings.ToLower(level) {
	case "native":
		return i.T(KeyNative)
	case "fluent":
		return i.T(KeyFluent)
	case "advanced":
		return i.T(KeyAdvanced)
	case "intermediate":
		return i.T(KeyIntermediate)
	case "basic", "beginner":
		return i.T(KeyBasic)
	default:
		return level
	}
}

// GetLanguageName returns the display name for a locale.
func GetLanguageName(locale Locale) string {
	switch locale {
	case LocaleEnUS:
		return "English (US)"
	case LocalePtBR:
		return "Português (BR)"
	case LocaleEsES:
		return "Español (ES)"
	case LocaleFrFR:
		return "Français (FR)"
	case LocaleDeDE:
		return "Deutsch (DE)"
	default:
		return string(locale)
	}
}

// SupportedLocales returns all supported locales.
func SupportedLocales() []Locale {
	return []Locale{
		LocaleEnUS,
		LocalePtBR,
		LocaleEsES,
		LocaleFrFR,
		LocaleDeDE,
	}
}

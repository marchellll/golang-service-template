package middleware

import (
	"golang-service-template/internal/errz"
	"net/http"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id" // Indonesian
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
)

// ValidatorMiddleware adds request validation capabilities
// This middleware sets up a validator instance with localized translations based on Accept-Language header
func ValidatorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get language from Accept-Language header
			acceptLang := c.Request().Header.Get("Accept-Language")
			lang := parseAcceptLanguage(acceptLang)

			// Create validator instance
			validate := validator.New()

			// Set up locales and universal translator
			english := en.New()
			indonesian := id.New()
			uni := ut.New(english, english, indonesian)

			// Get translator based on language
			trans, found := uni.GetTranslator(lang)
			if !found {
				// Fallback to English if language not supported
				trans, _ = uni.GetTranslator("en")
			}

			// Register translations based on language
			switch lang {
			case "id":
				_ = id_translations.RegisterDefaultTranslations(validate, trans)
			default:
				_ = en_translations.RegisterDefaultTranslations(validate, trans)
			}

			// Add both validator and translator to context for use by handlers
			c.Set("validator", validate)
			c.Set("translator", trans)

			return next(c)
		}
	}
}

// parseAcceptLanguage extracts the primary language from Accept-Language header
func parseAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return "en" // Default to English
	}

	// Parse Accept-Language header (e.g., "en-US,en;q=0.9,id;q=0.8")
	// Take the first language preference
	langs := strings.Split(acceptLang, ",")
	if len(langs) == 0 {
		return "en"
	}

	// Get the first language and remove quality factor if present
	firstLang := strings.TrimSpace(langs[0])
	if idx := strings.Index(firstLang, ";"); idx != -1 {
		firstLang = firstLang[:idx]
	}

	// Extract language code (e.g., "en-US" -> "en")
	if idx := strings.Index(firstLang, "-"); idx != -1 {
		firstLang = firstLang[:idx]
	}

	// Support specific languages, fallback to English
	switch strings.ToLower(firstLang) {
	case "id", "in": // Indonesian
		return "id"
	default:
		return "en"
	}
}

// ValidateRequest is a helper function that can be used in handlers to validate request structs
func ValidateRequest(c echo.Context, req interface{}) error {
	// Bind the request
	if err := c.Bind(req); err != nil {
		return errz.NewPrettyError(
			http.StatusBadRequest,
			"invalid_request_format",
			"Invalid request format",
			err,
		)
	}

	// Get validator from context
	v, ok := c.Get("validator").(*validator.Validate)
	if !ok {
		return errz.NewPrettyError(
			http.StatusInternalServerError,
			"validator_not_found",
			"Request validator not configured",
			errors.New("validator not found in context"),
		)
	}

	// Get translator from context
	trans, transOk := c.Get("translator").(ut.Translator)
	if !transOk {
		return errz.NewPrettyError(
			http.StatusInternalServerError,
			"translator_not_found",
			"Request translator not configured",
			errors.New("translator not found in context"),
		)
	}

	// Validate the request
	if err := v.Struct(req); err != nil {
		// Convert validation errors to user-friendly format using translator
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			// Use the translator to get localized error messages
			errorMap := validationErrors.Translate(trans)

			// Collect error messages for the main message
			var errorMessages []string
			for _, errMsg := range errorMap {
				errorMessages = append(errorMessages, errMsg)
			}

			return errz.NewPrettyErrorDetail(
				http.StatusBadRequest,
				"validation_failed",
				"Request validation failed: "+strings.Join(errorMessages, ", "),
				err,
				errorMap, // Pass the detailed error map
			)
		}

		return errz.NewPrettyError(
			http.StatusBadRequest,
			"validation_error",
			"Request validation error",
			err,
		)
	}

	return nil
}

// GetValidator returns the validator instance from the echo context
func GetValidator(c echo.Context) *validator.Validate {
	if v, ok := c.Get("validator").(*validator.Validate); ok {
		return v
	}
	return nil
}

// GetTranslator returns the translator instance from the echo context
func GetTranslator(c echo.Context) ut.Translator {
	if trans, ok := c.Get("translator").(ut.Translator); ok {
		return trans
	}
	return nil
}

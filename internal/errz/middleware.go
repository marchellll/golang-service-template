package errz

import (
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/labstack/echo/v4"
)

func ErrorRendererMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err == nil {
				return nil
			}

			// handle custom PrettyError
			var prettyError PrettyError
			if errors.As(err, &prettyError) {
				errorResponse := map[string]any{
					"code":    prettyError.Code,
					"message": prettyError.Message,
				}

				// Include details if they exist
				if len(prettyError.Details) > 0 {
					errorResponse["details"] = prettyError.Details
				}

				return c.JSON(prettyError.HttpStatusCode, map[string]any{
					"meta": map[string]any{
						"http_status_code": prettyError.HttpStatusCode,
					},
					"error": errorResponse,
				})
			}

			// handle Echo's HTTPError (like 404, 405, etc.)
			var httpError *echo.HTTPError
			if errors.As(err, &httpError) {
				return c.JSON(httpError.Code, map[string]any{
					"meta": map[string]any{
						"http_status_code": httpError.Code,
					},
					"error": map[string]any{
						"code":    http.StatusText(httpError.Code),
						"message": httpError.Message,
					},
				})
			}

			// fallback for any other errors
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"meta": map[string]any{
					"http_status_code": http.StatusInternalServerError,
				},
				"error": map[string]any{
					"code":    "ttts",
					"message": "Tetap Tenang, Tetap Semangat",
				},
			})
		}
	}
}

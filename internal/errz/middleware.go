package errz

import (
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/labstack/echo/v4"
)

func ErrzMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err == nil {
				return nil
			}

			// handle error
			var prettyError PrettyError
			if errors.As(err, &prettyError) {
				return c.JSON(prettyError.HttpStatusCode, map[string]any{
					"meta": map[string]any{
						"status": prettyError.HttpStatusCode,
					},
					"error": map[string]any{
						"code":    prettyError.Code,
						"message": prettyError.Message,
					},
				})
			}

			return c.JSON(http.StatusInternalServerError, map[string]any{
				"meta": map[string]any{
					"status": http.StatusInternalServerError,
				},
				"error": map[string]any{
					"code": "ttts",
					"message": "Tetap Tenang, Tetap Semangat",
				},
			})
		}
	}
}

package response

import "github.com/labstack/echo/v4"

func Success(c echo.Context, data interface{}) error {
	return c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func Error(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

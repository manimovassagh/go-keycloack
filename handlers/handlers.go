package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func SecureEndpoint(c echo.Context) error {
	// te
	user := c.Get("user").(jwt.MapClaims)
	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Hello, %s!", user["preferred_username"]),
	})
}
func PublicEndpoint(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]string{
		"message": "This is a public endpoint.",
	})
}

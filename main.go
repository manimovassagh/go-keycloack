package main

import (
	"fmt"
	"key/auth"
	"key/config"
	"key/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	resultConfig := config.NewConfig()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())



	// Authentication route
	// Public route
	e.GET("/public", handlers.PublicEndpoint)


	// Group for secured routes
	secured := e.Group("/secure")
	secured.Use(auth.TokenAuthMiddleware(resultConfig))
	secured.GET("", handlers.SecureEndpoint)

	err := e.Start(":4000")
	if err != nil {
		fmt.Println(err)
		return
	}
}

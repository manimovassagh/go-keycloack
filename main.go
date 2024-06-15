package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v4" // Use the latest version of jwt-go
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)
//nnnnnn
var (
	keycloakURL   = "http://localhost:8080"
	realm         = "mani"
	clientID      = "go"
	clientSecret  = "aZbUGUvcXxETMznmtsFB1Mm8TEWK4Hpz"
	gocloakClient = gocloak.NewClient(keycloakURL)
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(TokenAuthMiddleware)

	e.GET("/secure", SecureEndpoint)

	err := e.Start(":4000")
	if err != nil {
		return
	}
}

func TokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, "Missing Authorization Header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, "Invalid Authorization Token")
		}

		token, claims, err := ValidateToken(c.Request().Context(), tokenString)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid token")
		}

		c.Set("user", claims)
		return next(c)
	}
}

func ValidateToken(ctx context.Context, tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	certs, err := gocloakClient.GetCerts(ctx, realm)
	if err != nil {
		return nil, nil, err
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid := token.Header["kid"].(string)
		for _, cert := range *certs.Keys { // Dereference the pointer to slice
			if cert.Kid != nil && *cert.Kid == kid {
				if cert.X5c != nil && len(*cert.X5c) > 0 {
					certPEM := fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", (*cert.X5c)[0])
					block, _ := pem.Decode([]byte(certPEM))
					if block == nil {
						return nil, fmt.Errorf("failed to parse certificate PEM")
					}
					parsedCert, err := x509.ParseCertificate(block.Bytes)
					if err != nil {
						return nil, fmt.Errorf("failed to parse certificate: %v", err)
					}
					return parsedCert.PublicKey, nil
				}
			}
		}
		return nil, fmt.Errorf("unable to find appropriate key")
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}

	return token, claims, nil
}

func SecureEndpoint(c echo.Context) error {
	user := c.Get("user").(jwt.MapClaims)
	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Hello, %s!", user["preferred_username"]),
	})
}

func GetToken(ctx context.Context) (string, error) {
	token, err := gocloakClient.LoginClient(ctx, clientID, clientSecret, realm)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func ClientAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := GetToken(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to obtain token")
		}

		c.Request().Header.Set("Authorization", "Bearer "+tokenString)

		return next(c)
	}
}

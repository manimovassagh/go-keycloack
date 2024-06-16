package auth

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"key/config"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func TokenAuthMiddleware(config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, "Missing Authorization Header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, "Invalid Authorization Token")
			}

			token, claims, err := ValidateToken(c.Request().Context(), config, tokenString)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, "Invalid token")
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}

func ValidateToken(ctx context.Context, config *config.Config, tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	certs, err := config.GocloakClient.GetCerts(ctx, config.Realm)
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

func GetToken(ctx context.Context, config *config.Config) (string, error) {
	token, err := config.GocloakClient.LoginClient(ctx, config.ClientID, config.ClientSecret, config.Realm)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func ClientAuthMiddleware(config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := GetToken(c.Request().Context(), config)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, "Failed to obtain token")
			}

			c.Request().Header.Set("Authorization", "Bearer "+tokenString)

			return next(c)
		}
	}
}

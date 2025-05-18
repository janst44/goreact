package main

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (app *application) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Invalid authorization header format",
				})
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(app.jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Invalid token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Invalid token claims",
				})
			}

			// Check expiration
			// if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			// 	return c.JSON(http.StatusUnauthorized, echo.Map{
			// 		"error": "Token has expired",
			// 	})
			// }

			userId := claims["userId"].(string)

			user, err := app.models.Users.Get(userId)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized access"})

			}

			c.Set("user", user)

			return next(c)
		}
	}
}

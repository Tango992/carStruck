package middlewares

import (
	"fmt"
	"carstruck/utils"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Checks if use has "Authorization" cookie. 
// If Authorization does not exist / invalid, it will return a 401 Unauthorized error
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return echo.NewHTTPError(utils.ErrUnauthorized.Details("Authorization cookie does not exist"))
		}

		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ! ok {
				return nil, fmt.Errorf("failed to verify token signature")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return echo.NewHTTPError(utils.ErrUnauthorized.Details(err.Error()))
		}
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			return next(c)
		}
		return echo.NewHTTPError(utils.ErrUnauthorized.Details("Please log in to access this page"))
	}
}
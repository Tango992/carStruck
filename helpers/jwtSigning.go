package helpers

import (
	"carstruck/entity"
	"carstruck/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Create a new JWT, sign it, and set to a cookie
func SignNewJWT(c echo.Context, user entity.User) error {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(2 * time.Hour).Unix(),
		"id": user.ID,
		"email": user.Email,
		"deposit": user.Deposit,
		"full_name": user.FullName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Value = tokenString
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(2 * time.Hour)
	c.SetCookie(cookie)

	return nil
}
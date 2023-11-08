package helpers

import (
	"carstruck/entity"
	"carstruck/utils"
	"fmt"
	"net/smtp"
	"os"

	"github.com/labstack/echo/v4"
)

// Send a verification email
func SendVerificationEmail(user entity.User, verification entity.Verification) error {
	authEmail := os.Getenv("AUTH_EMAIL")
	authPass := os.Getenv("AUTH_PASS")
	smptHost := os.Getenv("SMPT_HOST")
	smptPort := os.Getenv("SMPT_PORT")
	verificationUrl := fmt.Sprintf("http://localhost:8080/users/verify/%v/%v", user.ID, verification.Token)
	
	body := []byte("To: " + user.Email + "\r\n" +
		"Subject: Carstruck Account Verification\r\n\r\n" +
		verificationUrl +"\r\n")
	
	smptAuth := smtp.PlainAuth("", authEmail, authPass, smptHost)
	smptAddr := fmt.Sprintf("%s:%s", smptHost, smptPort)

	err := smtp.SendMail(smptAddr, smptAuth, authEmail, []string{user.Email}, body)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	fmt.Printf("Verification email sent to %v\n", user.Email)
	return nil
}
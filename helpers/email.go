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
	
	url := fmt.Sprintf("https://carstruck-4d6b89ee5e4e.herokuapp.com/users/verify/%v/%v", user.ID, verification.Token)
	subject := "Subject: Carstruck Account Verification\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body, err := VerificationEmailBody(user.FullName, url)
	if err != nil {
		return err
	}
	
	smptAddr := fmt.Sprintf("%s:%s", smptHost, smptPort)
	smptAuth := smtp.PlainAuth("", authEmail, authPass, smptHost)
	msg := []byte(subject + mime + body)

	err = smtp.SendMail(smptAddr, smptAuth, authEmail, []string{user.Email}, msg)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	
	fmt.Printf("Verification email sent to %v\n", user.Email)
	return nil
}
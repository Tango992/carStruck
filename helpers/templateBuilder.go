package helpers

import (
	"bytes"
	"carstruck/utils"
	"html/template"

	"github.com/labstack/echo/v4"
)

func VerificationEmailBody(name, url string) (string, error) {
	tmpl, err := template.ParseFiles("./template/verificationEmail.html")
	if err != nil {
		return "", echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	
	templateData := map[string]string{
		"Name": name,
		"URL": url,
	}
	
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, templateData); err != nil {
		return "", echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return buf.String(), nil
}
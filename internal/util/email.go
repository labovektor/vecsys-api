package util

import (
	"embed"
	"html/template"
	"io"
	"strings"

	"github.com/labovector/vecsys-api/infrastructure/email"
	"gopkg.in/gomail.v2"
)

//go:embed template/*.html
var content embed.FS

// SendOTPEmail sends an OTP to the given email address.
//
// The email sent is an HTML template with the OTP embedded in it.
// The email is sent using the provided dialer.
//
// If the email is sent successfully, the function returns nil.
// Otherwise it returns an error.
func SendResetPasswordEmail(name, email, link string, dialer *email.EmailDialer) error {
	file, err := content.Open("template/otp.html")
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	templateString := string(bytes)

	templ, err := template.New("otp").Parse(templateString)
	if err != nil {
		return err
	}

	var body strings.Builder
	data := map[string]string{
		"Link": link,
		"Name": name,
	}
	if err := templ.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", dialer.SenderName)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Resert Password")
	m.SetBody("text/html", body.String())

	return dialer.DialAndSend(m)
}

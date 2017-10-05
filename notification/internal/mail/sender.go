package mail

import (
	"fmt"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/safe"
	"gopkg.in/gomail.v2"
)

var (
	dialer *gomail.Dialer

	smtpUsername = config.RegisterString("services.smtp.username", "", "smtp user name")
	smtpPassword = config.RegisterString("services.smtp.password", "", "smtp password")

	smtpHost = config.RegisterString("services.smtp.host", "0127.0.0.1", "smtp host")
	smtpPort = config.RegisterInt("services.smtp.port", 1025, "smtp port")
)

type EmailAddress struct {
	Email, Name string
}

type Message interface {
	fmt.Stringer
	Html() string
}

// NewEmailAddress return a new struct contain email adress and name, just a shortcut
func NewEmailAddress(mail string) EmailAddress {
	return EmailAddress{
		Email: mail,
		Name:  "",
	}
}

// NewEmailNameAddress return a new struct contain email adress and name, just a shortcut
func NewEmailNameAddress(mail, name string) EmailAddress {
	return EmailAddress{
		Email: mail,
		Name:  name,
	}
}

// Send sends Email to client
func Send(subject, msg string, from EmailAddress, to ...EmailAddress) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from.Email, from.Name)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	toString := make([]string, len(to))
	for i := range to {
		toString[i] = m.FormatAddress(to[i].Email, to[i].Name)
	}
	m.SetHeader("To", toString...)

	// No need to wait for result. its better to have the user and just record the
	// exception here :)
	safe.GoRoutine(func() {
		assert.Nil(dialer.DialAndSend(m))
	})
}

func init() {
	dialer = gomail.NewDialer(smtpHost.String(), smtpPort.Int(), smtpUsername.String(), smtpPassword.String())
}

package mail

import (
	"github.com/clickyab/services/assert"
	"github.com/go-gomail/gomail"
)

// Send sends Email to client
func Send(subject, msg, cc, bcc, fileAddr string, to ...string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)
	if cc != "" {
		m.SetHeader("Cc", cc)
	}
	if bcc != "" {
		m.SetHeader("Bcc", bcc)
	}
	if fileAddr != "" {
		m.Attach(fileAddr)
	}

	for i := range to {
		m.SetHeader("To", to[i])
		assert.Nil(dialer.DialAndSend(m))
	}
}

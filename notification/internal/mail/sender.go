package mail

import (
	"github.com/go-gomail/gomail"
)

// Send sends Email to client
func Send(subject, msg, fileAddr string, cc, bcc []string, to ...string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)
	if len(cc) != 0 {
		m.SetHeader("Cc", cc...)
	}
	if len(bcc) != 0 {
		m.SetHeader("Bcc", bcc...)
	}
	if fileAddr != "" {
		m.Attach(fileAddr)
	}

	for i := range to {
		m.SetHeader("To", to[i])
		dialer.DialAndSend(m)
	}
}

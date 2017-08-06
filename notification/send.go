package notification

import (
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/notification/internal/mail"
)

type (
	// Notiftype is the type of a notification
	Notiftype int
)

type MailConfig struct {
	Cc       string
	FileAddr string
	Bcc      string
}

// SMSConfig is the config for sms notification
type SMSConfig struct{}

// Packet is all client should fill for a notification
type Packet struct {
	To         []string
	Platform   Notiftype
	MailConfig *MailConfig
	SMSConfig  *SMSConfig
}

// MaxLength is used to check weather its out of range or not
func (p Packet) MaxLength() int {
	if p.Platform == SMSType {
		return 160
	}

	// almost unlimited chars for email
	return 1000000
}

const (
	// SMSType is the sms notification platform
	SMSType Notiftype = iota
	// MailType is the sms notification platform
	MailType
)

// Send sends a notification by its notification type
func Send(subject string, msg string, p ...Packet) {
	for i := range p {
		assert.True(p[i].MaxLength() > len(msg))

		switch p[i].Platform {
		case MailType:
			conf := p[i].MailConfig
			assert.NotNil(conf)
			mail.Send(subject, msg, conf.Cc, conf.Bcc, conf.FileAddr, p[i].To...)
		case SMSType:
			assert.NotNil(p[i].SMSConfig)
			//do some
		}
	}
}

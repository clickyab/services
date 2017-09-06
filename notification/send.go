package notification

import (
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/notification/internal/mail"
	"github.com/sirupsen/logrus"
)

type (
	// Notiftype is the type of a notification
	Notiftype int
)

const (
	// SMSType is the sms notification platform
	SMSType Notiftype = iota
	// MailType is the sms notification platform
	MailType
)

type mailConfig struct {
	Cc       []string
	Bcc      []string
	FileAddr string
}

// SMSConfig is the config for sms notification
type smsConfig struct{}

// Packet is all client should fill for a notification
type Packet struct {
	To       []string
	Platform Notiftype

	mailConfig *mailConfig
	smsConfig  *smsConfig
}

// MaxLength is used to check weather its out of range or not
func (p Packet) MaxLength() int {
	if p.Platform == SMSType {
		return 160
	}

	// almost unlimited chars for email
	return 1000000
}

// SetMailConfig set mail config
func (p *Packet) SetMailConfig(cc, bcc []string, fileAddr string) {
	if p.Platform != MailType {
		logrus.Panic("cant set mail config for a non mail packet")
	}
	p.mailConfig = &mailConfig{
		FileAddr: fileAddr,
		Bcc:      bcc,
		Cc:       cc,
	}
}

// SetSMSConfig sets sms config
func (p *Packet) SetSMSConfig(cc, bcc []string, fileAddr string) {
	if p.Platform != SMSType {
		logrus.Panic("cant set mail config for a non mail packet")
	}
	// TODO needs implementation
	p.smsConfig = &smsConfig{}
}

// Send sends a notification by its notification type
func Send(subject string, msg string, p ...Packet) {
	for i := range p {
		assert.True(p[i].MaxLength() > len(msg))

		switch p[i].Platform {
		case MailType:
			conf := p[i].mailConfig
			if conf == nil {
				conf = &mailConfig{}
			}
			mail.Send(subject, msg, conf.FileAddr, conf.Cc, conf.Bcc, p[i].To...)
		case SMSType:
			if p[i].smsConfig == nil {
				//some
			}
			//do some
		}
	}
}

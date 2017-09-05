package mail

import (
	"github.com/clickyab/services/config"
	"github.com/go-gomail/gomail"
)

var (
	dialer *gomail.Dialer
	from   = config.RegisterString("crab.mail_service.from", "info@clickyab.com", "from part of mail")

	smtpUsername = config.GetStringDefault("services.smtp.username", "")
	smtpPassword = config.GetStringDefault("services.smtp.password", "")

	smtpHost = config.GetStringDefault("services.smtp.host", "0.0.0.0")
	smtpPort = config.GetIntDefault("services.smtp.address_port", 1025)
)

func init() {
	dialer = gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)
}

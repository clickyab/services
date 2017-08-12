package mail

import (
	"github.com/clickyab/services/config"
	"github.com/go-gomail/gomail"
)

var (
	dialer *gomail.Dialer
	from   = config.RegisterString("crab.mail_service.from", "clickyab.com", "from part of mail")

	smtpUsername = config.GetString("crab.smtp.username")
	smtpPassword = config.GetString("crab.smtp.password")
	smtpHost     = config.GetString("crab.smtp.host")
	smtpPort     = config.GetInt("crab.smtp.address_port")
)

func init() {
	dialer = gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)
}

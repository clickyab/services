package mail

import (
	"github.com/clickyab/services/config"
	"github.com/go-gomail/gomail"
)

var (
	dialer *gomail.Dialer
	from   = config.RegisterString("crab.mail_service.from", "info@clickyab.com", "from part of mail")

	smtpUsername = config.GetString("services.smtp.username")
	smtpPassword = config.GetString("services.smtp.password")
	smtpHost     = config.GetString("services.smtp.host")
	smtpPort     = config.GetInt("services.smtp.address_port")
)

func init() {
	dialer = gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)
}

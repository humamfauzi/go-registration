package main

import (
	"fmt"
	"net/smtp"

	"github.com/humamfauzi/go-registration/utils"
)

type Notification struct {
	Code    string
	Message string
	Meta    map[string]interface{}
}

var (
	emailHost        string
	defaultSender    string
	defaultEmailUser string
	defaultEmailPass string
	defaultPort      string
	auth             smtp.Auth

	defaultHost string
)

func GetCredentials() {
	email := utils.GetEnv("email").(map[string]interface{})
	defaultHost = utils.GetEnv("host").(string)
	emailHost = email["host"].(string)
	defaultSender = email["sender"].(string)
	auth = smtp.PlainAuth("", defaultEmailUser, defaultEmailPass, emailHost)

}

func NotificationRouter(notif Notification) error {
	switch notif.Code {
	case "NOTIF_FORGOT_PASSWORD_EMAIL":
		userName := notif.Meta["name"].(string)
		timeRequest := notif.Meta["time"].(string)
		token := notif.Meta["token"].(string)

		template := GetTemplateString(notif.Code)
		payload := fmt.Sprintf(template, userName, timeRequest, defaultHost, token)

		recipient := notif.Meta["recipient"].(string)
		return EmailNotification([]string{recipient}, []byte(payload))
	default:
		return nil
	}
}

func GetTemplateString(code string) string {
	switch code {
	case "NOTIF_FORGOT_PASSWORD_EMAIL":
		return `
		Hi %s
		We receive request for change request forgot password request at %s
		Follow this link below to change your password
		%sforgot-password?token=%s

		Administration Team
		`
	default:
		return ""
	}
}

func EmailNotification(recipients []string, payload []byte) error {
	err := smtp.SendMail(emailHost, auth, defaultSender, recipients, payload)
	if err != nil {
		return err
	}
	return nil
}

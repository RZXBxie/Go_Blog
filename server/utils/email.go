package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"server/global"
	"strings"
)

func SendEmail(To, subject, body string) error {
	to := strings.Split(To, ",")
	return sendEmail(to, subject, body)
}

func sendEmail(to []string, subject, body string) error {
	emailConfig := global.Config.Email
	from := emailConfig.From
	nickname := emailConfig.Nickname
	secret := emailConfig.Secret
	port := emailConfig.Port
	host := emailConfig.Host
	isSSL := emailConfig.IsSSL

	// 使用 PlainAuth 创建认证信息
	auth := smtp.PlainAuth("", from, secret, host)

	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}

	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	var err error
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}

	return err
}

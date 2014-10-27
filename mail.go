package goutils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strings"
)

type SmtpAccount struct {
	Host, Port, User, Password string
}

type Attachment struct {
	Path, ContentType string
}

func SendMail(smtpAccount SmtpAccount, from mail.Address, to []mail.Address, subject string, body string, attachments []Attachment) error {
	const marker = "GRASSONE"
	parts := ""

	// part 1 will be the mail headers
	tos := make([]string, len(to))
	for i := 0; i < len(to); i++ {
		tos[i] = to[i].String()
	}

	parts += fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n--%s",
		from.String(),
		strings.Join(tos, ","),
		subject,
		marker,
		marker,
	)

	// part 2 will be the body of the email (HTML)
	parts += fmt.Sprintf("\r\nContent-Type: text/html\r\nContent-Transfer-Encoding:8bit\r\n\r\n%s\r\n--%s",
		body,
		marker,
	)

	// parts 3..n will be the attachments
	if attachments != nil {
		for _, a := range attachments {
			content, err := ioutil.ReadFile(a.Path)
			if err != nil {
				return err
			}
			encoded := base64.StdEncoding.EncodeToString(content)
			lineMaxLength := 500
			nbrLines := len(encoded) / lineMaxLength

			//append lines to buffer
			var buf bytes.Buffer
			for i := 0; i < nbrLines; i++ {
				buf.WriteString(encoded[i*lineMaxLength:(i+1)*lineMaxLength] + "\n")
			} //for

			//append last line in buffer
			buf.WriteString(encoded[nbrLines*lineMaxLength:])
			name := filepath.Base(a.Path)
			parts += fmt.Sprintf("\r\nContent-Type: %s; name=\"%s\"\r\nContent-Transfer-Encoding:base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n%s\r\n--%s--",
				a.ContentType,
				name,
				name,
				buf.String(),
				marker,
			)
		}
	}

	// finally let's send this email
	auth := smtp.PlainAuth(
		"",
		smtpAccount.User,
		smtpAccount.Password,
		smtpAccount.Host,
	)

	for i := 0; i < len(to); i++ {
		tos[i] = to[i].Address
	}

	return smtp.SendMail(
		smtpAccount.Host+":"+smtpAccount.Port,
		auth,
		from.Address,
		tos,
		[]byte(parts),
	)
}

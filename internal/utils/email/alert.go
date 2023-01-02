package email

import (
	"fmt"

	"github.com/fox-one/pkg/uuid"
	"github.com/lixvyang/betxin/internal/utils"
)

// NotifyHandler 告警通知
func NotifyHandler(msg interface{}) {
	subject, body, err := newHTMLEmail(
		"method",
		utils.MailHost,
		"uri",
		uuid.New(),
		msg,
		"stack",
	)

	if err != nil {
		fmt.Println("newHTMLEmail error: ", err)
	}

	options := &Options{
		MailHost: utils.MailHost,
		MailPort: utils.MailPort,
		MailUser: utils.MailUser,
		MailPass: utils.MailPass,
		MailTo:   utils.MailTo,
		Subject:  subject,
		Body:     body,
	}
	if err := Send(options); err != nil {
		fmt.Println("send email error: ", err)
	}
}

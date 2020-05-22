package sendgrid

import (
	"bitbucket.org/walmartdigital/hermes/app/domain/entity"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
	"encoding/base64"
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	END_POINT = "/v3/mail/send"
	HOST      = "https://api.sendgrid.com"
)

type sendgridRepository struct {
	key      string
	endpoint string
	host     string
}

func NewSendgridRepository(sendrigApiToken string) *sendgridRepository {
	if len(sendrigApiToken) == 0 {
		log.Fatal("error, sengridApiToken cannot be empty")
	}
	return &sendgridRepository{
		key:      sendrigApiToken,
		endpoint: END_POINT,
		host:     HOST,
	}
}

func (s *sendgridRepository) Send(id string, email *entity.Email, templateId string) error {
	newV3Mail := mail.NewV3Mail()

	from := mail.NewEmail(email.From.Name, email.From.Address)
	content := mail.NewContent("text/html", email.Subject)
	to := mail.NewEmail(email.To.Name, email.To.Address)

	newV3Mail.SetFrom(from)
	newV3Mail.AddContent(content)

	newV3Mail.SetTemplateID(templateId)

	// create new *Personalization
	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.SetCustomArg("message_id", id)
	personalization.Subject = email.Subject

	for i, v := range email.Body {
		personalization.SetDynamicTemplateData(i, v)
	}

	// add `personalization` to `newV3Mail`
	newV3Mail.AddPersonalizations(personalization)

	// Attachment
	for _, itemAttachment := range email.Attachment {
		attachment := mail.NewAttachment()
		encoded := base64.StdEncoding.EncodeToString([]byte(itemAttachment.Content))
		attachment.SetContent(encoded)
		attachment.SetType(itemAttachment.Type)
		attachment.SetFilename(itemAttachment.FileName)
		attachment.SetDisposition("attachment")
		attachment.SetContentID(itemAttachment.FileName)
		newV3Mail.AddAttachment(attachment)
	}

	request := sendgrid.GetRequest(s.key, s.endpoint, s.host)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(newV3Mail)
	response, err := sendgrid.API(request)
	if err != nil {
		return errors.New("error sending email")
	}
	log.Info("StatusCode: %d, Body: %s, Headers: %s", response.StatusCode, response.Body, response.Headers)
	return nil
}

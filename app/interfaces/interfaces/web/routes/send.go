package routes

import (
	"bitbucket.org/walmartdigital/hermes/app/application/send_to_email"
	"bitbucket.org/walmartdigital/hermes/app/domain/entity"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/middleware/authentication"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/models"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/config"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/custom_errors"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type sendHandler struct {
	sendToEmailUseCase send_to_email.SendToEmailUseCase
}

func NewSendHandler(e *echo.Echo, sendToEmailUseCase send_to_email.SendToEmailUseCase) *sendHandler {
	sendHandler := &sendHandler{
		sendToEmailUseCase: sendToEmailUseCase,
	}

	e.POST("/email/send", sendHandler.Send, authentication.GetMiddlewareConfig())

	return sendHandler
}

func (s *sendHandler) Send(c echo.Context) error {
	log.Info("Request received to sent email")

	if !config.GetBool("feature.flags.send") {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Description: "email sending disabled"})
	}

	sendRequest := new(models.SendRequest)

	if err := c.Bind(sendRequest); err != nil {
		err := errors.New("error getting email payload")
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Description: err.Error()})
	}

	if err := c.Validate(sendRequest); err != nil {
		var msgError string
		var split string
		for _, e := range err.(validator.ValidationErrors) {
			msgError = fmt.Sprintf("%s%s%s", msgError, split, e)
			split = ", "
		}
		err := errors.New(fmt.Sprintf("error validating data structure: %s", msgError))
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Description: err.Error()})
	}

	email := &entity.Email{
		Subject: sendRequest.Content.Subject,
		From: entity.Contact{
			Address: sendRequest.Content.From.Address,
			Name:    sendRequest.Content.From.Name,
		},
		To: entity.Contact{
			Address: sendRequest.Content.To.Address,
			Name:    sendRequest.Content.To.Name,
		},
		Body: sendRequest.Content.Body,
	}

	appClient := authentication.GetClientToken(c)

	order := &entity.Order{Id: sendRequest.OrderID, Channel: sendRequest.Channel}
	event := &entity.Event{Label: sendRequest.Event, AppClient: appClient}

	messageStatus, err := s.sendToEmailUseCase.Send(email, order, event)
	if err != nil {
		customErr, ok := err.(*custom_errors.RequestError)
		if ok {
			switch customErr.Kind() {
			case custom_errors.AlreadyProcessed:
				return c.JSON(http.StatusAlreadyReported, models.SentResponse{
					MessageId:   messageStatus.MessageId.String(),
					Status:      messageStatus.Status,
					Description: customErr.Error(),
				})
			case custom_errors.EventNotAvailable:
				return c.JSON(http.StatusMethodNotAllowed, models.ErrorResponse{
					Kind:        customErr.Kind(),
					Description: err.Error(),
				})
			case custom_errors.MailProvideError, custom_errors.DataBaseError:
				return c.JSON(http.StatusFailedDependency, models.ErrorResponse{
					Kind:        customErr.Kind(),
					Description: err.Error(),
				})
			default:
				return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Kind: custom_errors.Unknown, Description: err.Error()})
			}
		}
	}
	return c.JSON(http.StatusCreated, models.SentResponse{
		MessageId:   messageStatus.MessageId.String(),
		Status:      messageStatus.Status,
		Description: "message sent successfully to email provider",
	})
}

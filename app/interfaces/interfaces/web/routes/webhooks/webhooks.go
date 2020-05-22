package webhooks

import (
	"bitbucket.org/walmartdigital/hermes/app/application/update_message_status"
	"github.com/labstack/echo/v4"
)

type sendGridWebHookHandler struct {
	updateMessageStatusUseCase update_message_status.UpdateMessageStatusUseCase
}

func NewSendHandler(e *echo.Echo, updateMessageStatusUseCase update_message_status.UpdateMessageStatusUseCase) *sendGridWebHookHandler {
	sendGridWebHookHandler := &sendGridWebHookHandler{
		updateMessageStatusUseCase: updateMessageStatusUseCase,
	}

	e.POST("/webhooks/sendgrid", sendGridWebHookHandler.SendgridHook)

	return sendGridWebHookHandler
}
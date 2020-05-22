package webhooks

import (
	"bitbucket.org/walmartdigital/hermes/app/domain/constant/status"
	"bitbucket.org/walmartdigital/hermes/app/domain/entity"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/models"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/config"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SendGridWebHookRequest struct {
	MessageId string `json:"message_id" validate:"required"`
	Event     string `json:"event" validate:"required"`
}

func (h *sendGridWebHookHandler) SendgridHook(c echo.Context) error {
	var err error
	log.Info("Request received to sendgrid webhooks")

	if !config.GetBool("feature.flags.update_states") {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Description: "webhooks sendgrid disabled"})
	}

	var events []SendGridWebHookRequest

	if err := c.Bind(&events); err != nil {
		err := errors.New("error getting events payload")
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Description: err.Error()})
	}

	//valid data
	//if err := c.Validate(&events); err != nil {
	//	var msgError string
	//	var split string
	//	log.Error(err.Error())
	//	for _, e := range err.(validator.ValidationErrors) {
	//		msgError = fmt.Sprintf("%s%s%s", msgError, split, e)
	//		split = ", "
	//	}
	//	err := errors.New(fmt.Sprintf("error validating data structure: %s", msgError))
	//	log.Error(err.Error())
	//	return c.JSON(http.StatusOK, models.ErrorResponse{Message: err.Error()})
	//}

	for _, v := range events {
		if len(v.MessageId) == 0 {
			continue
		}
		state := ""
		switch v.Event {
		case "processed":
			continue
		case "deferred":
			state = status.DEFERRED
		case "delivered":
			state = status.DELIVERED
		case "dropped":
			state = status.DROPPED
		case "bounce":
			state = status.BOUNCED
		case "open":
			state = status.OPENED
		default:
			continue
		}

		updateMessage := entity.RequestMessageUpdate{
			MessageId: v.MessageId,
			Status:    state,
		}

		err = h.updateMessageStatusUseCase.Update(updateMessage)
		if err != nil {
			break
		}
	}

	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Description: err.Error()})
	}
	return c.JSON(http.StatusCreated, models.ErrorResponse{Description: "ok"})
}

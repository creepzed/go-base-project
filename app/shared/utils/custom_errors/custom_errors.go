package custom_errors

import (
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
	"errors"
)

const (
	AlreadyProcessed  = "ALREADY_PROCESSED"
	TemplateNotFound  = "TEMPLATE_NOT_FOUND"
	EventNotAvailable = "EVENT_NOT_AVAILABLE"
	MailProvideError  = "MAIL_PROVIDER_ERROR"
	DataBaseError     = "DATA_BASE_ERROR"
	Unknown           = "UNKNOWN"
)

type RequestError struct {
	kind string
	err  error
}

func New(message string, kind string) error {
	log.Error("[%s] %s", kind, message)
	return &RequestError{
		kind: kind,
		err:  errors.New(message),
	}
}

func NewWithError(err error, kind string) error {
	log.Error("[%s] %s", kind, err.Error())
	return &RequestError{
		kind: kind,
		err:  err,
	}
}

func (custom *RequestError) Error() string {
	return custom.err.Error()
}

func (custom *RequestError) Kind() string {
	return custom.kind
}

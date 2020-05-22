package config_send_to_email

import (
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/config"
	"errors"
	"fmt"
)

type configSendToEmailRepository struct{}

func NewConfigSendToEmailRepository() *configSendToEmailRepository {
	return &configSendToEmailRepository{}
}

func (c *configSendToEmailRepository) GetEnabled(nameFeature string) bool {
	return config.GetBool(nameFeature)
}

func (c *configSendToEmailRepository) GetTemplateId(templateIdFeature string) (string, error) {
	value := config.GetString(templateIdFeature)
	if len(value) == 0 {
		return "", errors.New(fmt.Sprintf("error templateId %s undefined", templateIdFeature))
	}
	return value, nil
}

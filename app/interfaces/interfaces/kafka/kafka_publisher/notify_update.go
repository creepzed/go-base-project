package kafka_publisher

import (
	"bitbucket.org/walmartdigital/hermes/app/domain/entity"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/kafka/kafka_models"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/kafka"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
)

type kafkaNotify struct {
	topic     string
	publisher *kafka.KafkaPublisher
}

func NewKafkaNotify(topic string, brokers []string) *kafkaNotify {
	return &kafkaNotify{
		topic:     topic,
		publisher: kafka.NewKafkaPublisher(brokers...),
	}
}

func (k *kafkaNotify) Notify(messageStatus *entity.MessageStatus) error {
	messageKafka := kafka_models.MessageStatusKafka{
		MessageId: messageStatus.MessageId.String(),
		OrderId:   messageStatus.Order.Id,
		Channel:   messageStatus.Order.Channel,
		Event:     messageStatus.Event.Label,
		AppClient: messageStatus.Event.AppClient,
		Status:    messageStatus.Status,
	}

	err := k.publisher.Publish(k.topic, messageKafka)
	if err != nil {
		log.WithError(err).Error("error to notify update")
	}
	return nil
}

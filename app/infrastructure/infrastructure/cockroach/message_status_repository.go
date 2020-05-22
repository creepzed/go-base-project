package cockroach

import (
	"bitbucket.org/walmartdigital/hermes/app/domain/entity"
	"bitbucket.org/walmartdigital/hermes/app/infrastructure/cockroach/cockroach_connection"
	"bitbucket.org/walmartdigital/hermes/app/infrastructure/cockroach/db_model"
	"encoding/json"
	"errors"
)

type messageStatusCockroachRepository struct {
	connection cockroach_connection.CockroachConnection
}

func NewMessageStatusRepository(connection cockroach_connection.CockroachConnection) *messageStatusCockroachRepository {
	return &messageStatusCockroachRepository{
		connection: connection,
	}
}

func (r *messageStatusCockroachRepository) FindOrCreate(messageStatus *entity.MessageStatus) (*entity.MessageStatus, error) {
	db := r.connection.GetConnection()
	defer r.connection.CloseConnection()
	parsedPayload, _ := json.Marshal(messageStatus.Payload)

	messageStatusCockroach := db_model.MessageStatus{
		MessageId: messageStatus.MessageId,
		OrderId:   messageStatus.Order.Id,
		Channel:   messageStatus.Order.Channel,
		Type:      messageStatus.Type,
		Recipient: messageStatus.Payload.To.Address,
		Event:     messageStatus.Event.Label,
		AppClient: messageStatus.Event.AppClient,
		Status:    messageStatus.Status,
		Payload:   string(parsedPayload),
	}

	db.Where("order_id = ? and event = ?", messageStatus.Order.Id, messageStatus.Event.Label).FirstOrCreate(&messageStatusCockroach)
	if db.Error != nil {
		return &entity.MessageStatus{}, errors.New("error trying to Find or Create message status")
	}

	var email *entity.Email
	json.Unmarshal([]byte(messageStatusCockroach.Payload), &email)
	messageStatusFound := &entity.MessageStatus{
		MessageId: messageStatusCockroach.MessageId,
		Order:     entity.Order{Id: messageStatusCockroach.OrderId, Channel: messageStatusCockroach.Channel},
		Event:     entity.Event{Label: messageStatusCockroach.Event, AppClient: messageStatusCockroach.AppClient},
		Payload:   *email,
		Status:    messageStatusCockroach.Status,
	}

	return messageStatusFound, nil
}

func (r *messageStatusCockroachRepository) Update(messageStatus *entity.MessageStatus) error {
	db := r.connection.GetConnection()
	defer r.connection.CloseConnection()
	parsedPayload, _ := json.Marshal(messageStatus.Payload)

	messageStatusCockroach := db_model.MessageStatus{
		MessageId: messageStatus.MessageId,
		OrderId:   messageStatus.Order.Id,
		Channel:   messageStatus.Order.Channel,
		Event:     messageStatus.Event.Label,
		Status:    messageStatus.Status,
		Payload:   string(parsedPayload),
	}

	db.Model(&messageStatusCockroach).Where("message_id = ?", messageStatus.MessageId).Update(&messageStatusCockroach)
	if db.Error != nil {
		return errors.New("Error to update the message from database, ")
	}

	return nil
}
func (r *messageStatusCockroachRepository) UpdateStatusByMessageId(messageId string, status string) (*entity.MessageStatus, error) {
	db := r.connection.GetConnection()
	defer r.connection.CloseConnection()
	var messageStatusCockroach db_model.MessageStatus

	db.Model(messageStatusCockroach).Where("message_id = ?", messageId).Update("status", status)
	if db.Error != nil {
		return nil, errors.New("Error to update the status from database, ")
	}

	db.Where("message_id = ?", messageId).
		Find(&messageStatusCockroach)

	if db.Error != nil {
		return nil, errors.New("Error to get the message from database, ")
	}

	var email *entity.Email
	json.Unmarshal([]byte(messageStatusCockroach.Payload), &email)
	messageStatusFound := &entity.MessageStatus{
		MessageId: messageStatusCockroach.MessageId,
		Order:     entity.Order{Id: messageStatusCockroach.OrderId, Channel: messageStatusCockroach.Channel},
		Event:     entity.Event{Label: messageStatusCockroach.Event, AppClient: messageStatusCockroach.AppClient},
		Payload:   *email,
		Status:    messageStatusCockroach.Status,
		Type:      messageStatusCockroach.Type,
	}
	return messageStatusFound, nil
}

func (r *messageStatusCockroachRepository) GetByMessageStatusAndOrderId(event *entity.Event, order *entity.Order) (*entity.MessageStatus, error) {
	var messageStatusCockroach db_model.MessageStatus
	db := r.connection.GetConnection()
	defer r.connection.CloseConnection()

	db.Where("order_id = ?", order.Id).
		Find(&messageStatusCockroach)

	if db.Error != nil {
		return nil, errors.New("Error to get the message from database, ")
	}

	var email *entity.Email
	json.Unmarshal([]byte(messageStatusCockroach.Payload), &email)

	messageStatus := &entity.MessageStatus{
		MessageId: messageStatusCockroach.MessageId,
		Order:     entity.Order{Id: messageStatusCockroach.OrderId, Channel: messageStatusCockroach.Channel},
		Event:     entity.Event{Label: messageStatusCockroach.Event, AppClient: messageStatusCockroach.AppClient},
		Payload:   *email,
		Status:    messageStatusCockroach.Status,
		Type:      messageStatusCockroach.Type,
	}

	return messageStatus, nil
}

package kafka_consumer

type kafkaConsumerSendEmail struct {

}


/*func (k *kafkaShopConsumer) EcommercePaymentCreated() {
	topic := constant.EcommercePaymentCreated
	payload := &midas.EcommercePaymentCreated{}

	k.listenTopic(topic, payload, func(message kafka.Message) error {
		log.Printf("%s payload: %s", topic, utils.EntityToJson(payload))
		return k.usecase.EcommercePaymentCreated(payload, string(message.Key))
	})
}*/
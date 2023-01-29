package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hugosrc/shortlink/internal/util"
	"github.com/spf13/viper"
)

func NewProducer(config *viper.Viper) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetString("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": config.GetString("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanisms":   config.GetString("KAFKA_SASL_MECHANISMS"),
		"sasl.username":     config.GetString("KAFKA_SASL_USERNAME"),
		"sasl.password":     config.GetString("KAFKA_SASL_PASSWORD"),
		"compression.type":  config.GetString("KAFKA_PRODUCER_COMPRESSION_TYPE"),
		"retries":           config.GetString("KAFKA_PRODUCER_RETRIES"),
		"linger.ms":         config.GetString("KAFKA_PRODUCER_LINGER_MS"),
		"acks":              config.GetString("KAFKA_PRODUCER_ACKS"),
		"batch.size":        config.GetString("KAFKA_PRODUCER_BATCH_SIZE"),
	})
	if err != nil {
		return nil, util.WrapErrorf(err, util.ErrCodeUnknown, "couldn't create kafka producer")
	}

	return producer, nil
}

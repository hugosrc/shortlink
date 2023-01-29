package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hugosrc/shortlink/internal/core/domain"
	"github.com/hugosrc/shortlink/internal/util"
)

type KafkaMetricsProducer struct {
	topic    string
	producer *kafka.Producer
}

func NewKafkaMetricsProducer(topic string, producer *kafka.Producer) *KafkaMetricsProducer {
	return &KafkaMetricsProducer{
		topic:    topic,
		producer: producer,
	}
}

func (p *KafkaMetricsProducer) Produce(metrics *domain.LinkMetrics) error {
	metricsBytes, err := json.Marshal(metrics)
	if err != nil {
		return util.WrapErrorf(err, util.ErrCodeUnknown, "unable to marshal metrics data")
	}

	topicName := &p.topic
	if err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     topicName,
			Partition: kafka.PartitionAny,
		},
		Value: metricsBytes,
		Key:   []byte(metrics.ShortURL),
	}, nil); err != nil {
		return util.WrapErrorf(err, util.ErrCodeUnknown, "couldn't produce message")
	}

	return nil
}

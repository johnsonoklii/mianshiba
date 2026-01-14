package contract

import "context"

type KafkaProducer interface {
	SendMessage(ctx context.Context, topic string, key []byte, value []byte) error
	Close() error
}

type KafkaConsumer interface {
	Consume(ctx context.Context, topics []string, handler func(ctx context.Context, message *KafkaMessage) error) error
	Close() error
}

type KafkaMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
}

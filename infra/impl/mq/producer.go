package kafka

import (
	"context"
	"fmt"
	"mianshiba/conf"
	mq "mianshiba/infra/contract/mq"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type kafkaProducer struct {
	producer *kafka.Producer
}

func NewProducer(ctx context.Context) (mq.KafkaProducer, error) {
	// 使用配置创建Kafka生产者
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Global.Kafka.Brokers,
		"client.id":         "mianshiba-producer",
		"acks":              "1", // 至少一个副本确认
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %v", err)
	}

	// 启动一个goroutine来处理Kafka的事件
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &kafkaProducer{
		producer: p,
	}, nil
}

func (kp *kafkaProducer) SendMessage(ctx context.Context, topic string, key []byte, value []byte) error {
	// 创建消息
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}

	// 发送消息到Kafka
	return kp.producer.Produce(msg, nil)
}

func (kp *kafkaProducer) Close() error {
	// 等待所有消息发送完成
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()

	return nil
}

package kafka

import (
	"context"
	"fmt"
	"mianshiba/conf"
	mq "mianshiba/infra/contract/mq"
	"mianshiba/pkg/logs"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type kafkaConsumer struct {
	consumer *kafka.Consumer
}

// NewConsumer 创建一个新的Kafka消费者实例
func NewConsumer(ctx context.Context) (mq.KafkaConsumer, error) {
	// 构建Kafka配置
	config := &kafka.ConfigMap{
		"bootstrap.servers":  conf.Global.Kafka.Brokers,
		"group.id":           conf.Global.Kafka.GroupID,
		"auto.offset.reset":  "earliest", // 从最早的偏移量开始消费
		"enable.auto.commit": true,       // 自动提交偏移量
	}

	// 创建Kafka消费者
	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %v", err)
	}

	return &kafkaConsumer{
		consumer: c,
	}, nil
}

// Consume 开始消费指定主题的消息
func (kc *kafkaConsumer) Consume(ctx context.Context, topics []string, handler func(ctx context.Context, message *mq.KafkaMessage) error) error {
	// 订阅主题
	if err := kc.consumer.SubscribeTopics(topics, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topics: %v", err)
	}

	// 消费循环
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 从Kafka获取消息，超时时间为1秒
			msg, err := kc.consumer.ReadMessage(1000)
			if err != nil {
				// 处理超时或其他错误
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue // 超时，继续下一次循环
				}
				return fmt.Errorf("failed to read message: %v", err)
			}

			// 构建自定义KafkaMessage结构
			kafkaMsg := &mq.KafkaMessage{
				Topic:     *msg.TopicPartition.Topic,
				Partition: msg.TopicPartition.Partition,
				Offset:    int64(msg.TopicPartition.Offset),
				Key:       msg.Key,
				Value:     msg.Value,
			}

			// 调用消息处理器处理消息
			if err := handler(ctx, kafkaMsg); err != nil {
				logs.Errorf("Failed to handle message from topic %s, partition %d, offset %d: %v",
					kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset, err)
				// 消息处理失败时，根据需要进行重试或其他处理
			}
		}
	}
}

// Close 关闭Kafka消费者
func (kc *kafkaConsumer) Close() error {
	if kc.consumer != nil {
		kc.consumer.Close()
	}
	return nil
}

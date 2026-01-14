package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mianshiba/application"
	"mianshiba/application/agent/handler"
	"mianshiba/conf"
	cmq "mianshiba/infra/contract/mq"
	mq "mianshiba/infra/impl/mq"
	"mianshiba/pkg/logs"
)

func main() {
	// 1. 加载配置
	if err := conf.LoadConfig("./config.yaml"); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	logs.SetLevel(logs.LevelInfo)

	ctx := context.Background()
	if err := application.Init(ctx); err != nil {
		panic("InitializeInfra failed, err=" + err.Error())
	}

	// 3. 创建Kafka消费者
	consumer, err := mq.NewConsumer(ctx)
	if err != nil {
		logs.Errorf("Failed to create kafka consumer: %v", err)
		os.Exit(1)
	}
	defer consumer.Close()

	// 4. 定义要消费的主题
	topics := []string{conf.Global.Kafka.ResumeTopic}

	// 5. 启动消费循环
	go func() {
		if err := consumer.Consume(ctx, topics, handleMessage); err != nil {
			logs.Errorf("Failed to consume messages: %v", err)
			os.Exit(1)
		}
	}()

	logs.Infof("Kafka consumer started, consuming topics: %v", topics)

	// 6. 等待终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logs.Infof("Shutting down kafka consumer...")
}

// handleMessage 处理收到的Kafka消息
func handleMessage(ctx context.Context, msg *cmq.KafkaMessage) error {
	// 1. 将Kafka消息转换为领域事件
	domainEvent, err := handler.ConvertToResumeParseDomainEvent(msg)
	if err != nil {
		logs.Errorf("Failed to convert message to domain event: %v", err)
		return err
	}

	// 2. 调用事件处理器处理领域事件
	if err := handler.ResumeHandlerSVC.HandleResumeParseEvent(ctx, domainEvent); err != nil {
		logs.Errorf("Failed to handle ResumeParseEvent: %v", err)
		return err
	}

	// 3. 记录处理完成日志
	logs.Infof("Processed message: topic=%s, partition=%d, offset=%d, key=%s, value=%s",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	return nil
}

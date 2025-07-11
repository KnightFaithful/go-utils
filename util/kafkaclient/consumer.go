package kafkaclient

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

func singleConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 创建消费者
	consumer, err := sarama.NewConsumer([]string{"xxx"}, config)
	if err != nil {
		log.Fatalf("消费者创建失败: %v", err)
	}
	defer consumer.Close()

	// 订阅主题分区（这里消费分区 0）
	partitionConsumer, err := consumer.ConsumePartition("yyy", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("分区消费者创建失败: %v", err)
	}
	defer partitionConsumer.Close()

	fmt.Println("开始...")

	// 消费消息
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("收到消息: topic=%s, partition=%d, offset=%d, key=%s, value=%s\n",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("消费错误: %v\n", err)
		}
	}
}

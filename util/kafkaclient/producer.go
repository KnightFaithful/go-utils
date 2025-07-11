package kafkaclient

import (
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
)

func asyncProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("异步生产者创建失败: %v", err)
	}
	defer producer.Close()

	// 监听发送结果
	go func() {
		for {
			select {
			case success := <-producer.Successes():
				log.Printf("异步发送成功: topic=%s, partition=%d, offset=%d\n",
					success.Topic, success.Partition, success.Offset)
			case err := <-producer.Errors():
				log.Printf("异步发送失败: %v\n", err.Err)
			}
		}
	}()

	// 发送消息
	msg := &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("Async Message"),
	}
	producer.Input() <- msg

	// 防止主线程退出
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}

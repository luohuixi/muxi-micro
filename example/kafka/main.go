package main

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/muxi-Infra/muxi-micro/pkg/kafka"
)

// 同步生产者示例
func SyncProducer() {
	p := kafka.NewSyncConfig([]string{"localhost:9092"})
	//可自行改配置如：（可选）
	//p.GetConfig().Producer.RequiredAcks = sarama.WaitForAll

	err := p.CreateProducer()
	if err != nil {
		log.Fatal(err)
	}

	partition, offset, err := p.SendMessage("study2", "testkey8", "testvalue8")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("partition: %d, offset: %d\n", partition, offset)

	err = p.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// 异步生产者示例
func AsyncProducer() {
	p := kafka.NewAsyncConfig([]string{"localhost:9092"})
	//可自行改配置如：（可选）
	p.GetConfig().Producer.Return.Successes = true

	err := p.CreateProducer()
	if err != nil {
		log.Fatal(err)
	}

	//异步获取成功消息
	ch := p.GetSuccess()
	//或 ch:= p.GetError()
	go func() {
		for msg := range ch {
			log.Printf("partition: %d, offset: %d, key: %s, value: %s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
		}
	}()

	p.SendMessage("study2", "testkey5", "testvalue5")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	err = p.Close()
	if err != nil {
		log.Fatal(err)
	}

}

// 消费者示例
func Consumer() {
	// 先构造 handler 实现 ConsumerGroupHandler 接口
	// setup, cleanup, consumeclaim 都要用户自行实现
	fn := func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
		for msg := range claim.Messages() {
			log.Println("partition: ", msg.Partition, "offset: ", msg.Offset, "key: ", string(msg.Key), "value: ", string(msg.Value))
			session.MarkMessage(msg, "")
		}
		return nil // 正常退出
	}
	handler := kafka.NewHandler(nil, nil, fn)

	// 构造消费者组
	c := kafka.NewConsumerConfig([]string{"localhost:9092"}, "test-group", []string{"study2"}, handler)
	//可自行改配置如：（可选）
	//c.GetConfig().Consumer.Offsets.Initial = sarama.OffsetOldest

	err := c.CreateConsumerGroup()
	if err != nil {
		log.Fatal(err)
	}

	// 消费消息
	ctx, cancel := context.WithCancel(context.Background())

	// 取消消费
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()

	err = c.Consume(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	SyncProducer()
	//AsyncProducer()
	Consumer()
}

package kafka

import (
	"github.com/IBM/sarama"
)

type Sync_producer interface {
	CreateProducer() error                                      //创建生产者
	GetConfig() *sarama.Config                                  //获取配置
	SendMessage(topic, key, value string) (int32, int64, error) //发送消息
	Close() error                                               //关闭生产者
}

type Async_producer interface {
	Sync_producer
	GetSuccess() <-chan *sarama.ProducerMessage //异步获取成功消息
	GetError() <-chan *sarama.ProducerError     //异步获取错误消息
}

// 同步生成者
type SyncProducer struct {
	host     []string            //kafka的地址
	config   *sarama.Config      //配置
	producer sarama.SyncProducer //生产者
}

func NewSyncConfig(host []string) Sync_producer {
	config := sarama.NewConfig()
	// 默认是false
	config.Producer.Return.Successes = true // 必须设为 true

	// 返回允许用户自己改配置
	return &SyncProducer{
		host:   host,
		config: config,
	}
}

func (p *SyncProducer) GetConfig() *sarama.Config {
	return p.config
}

func (p *SyncProducer) CreateProducer() error {
	producer, err := sarama.NewSyncProducer(p.host, p.config)
	if err != nil {
		return err
	}

	p.producer = producer
	return nil
}

func (p *SyncProducer) SendMessage(topic, key, value string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return -1, -1, err
	}
	return partition, offset, nil
}

func (p *SyncProducer) Close() error {
	return p.producer.Close()
}

// 异步生成者
type AsyncProducer struct {
	host     []string             //kafka的地址
	config   *sarama.Config       //配置
	producer sarama.AsyncProducer //生产者
}

func NewAsyncConfig(host []string) Async_producer {
	config := sarama.NewConfig()
	// 默认是 Waitforlocal
	config.Producer.RequiredAcks = sarama.NoResponse

	// 返回允许用户自己改配置
	return &AsyncProducer{
		host:   host,
		config: config,
	}
}

func (p *AsyncProducer) GetConfig() *sarama.Config {
	return p.config
}

func (p *AsyncProducer) GetSuccess() <-chan *sarama.ProducerMessage {
	return p.producer.Successes()
}

func (p *AsyncProducer) GetError() <-chan *sarama.ProducerError {
	return p.producer.Errors()
}

func (p *AsyncProducer) CreateProducer() error {
	producer, err := sarama.NewAsyncProducer(p.host, p.config)
	if err != nil {
		return err
	}

	p.producer = producer
	return nil
}

func (p *AsyncProducer) SendMessage(topic, key, value string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	p.producer.Input() <- msg

	//异步不能在这里查错
	return -1, -1, nil
}

func (p *AsyncProducer) Close() error {
	return p.producer.Close()
}

package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type Consumer interface {
	CreateConsumerGroup() error        //创建消费者组
	GetConfig() *sarama.Config         //获取配置
	Consume(ctx context.Context) error //消费信息
	Close() error                      //关闭消费者
}

// 实现 ConsumerGroupHandler 接口
type Handler struct {
	setupFunc        func(sarama.ConsumerGroupSession) error
	cleanupFunc      func(sarama.ConsumerGroupSession) error
	consumeclaimFunc func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

func (c *Handler) Setup(session sarama.ConsumerGroupSession) error {
	if c.setupFunc == nil {
		return nil
	}
	return c.setupFunc(session)
}

func (c *Handler) Cleanup(session sarama.ConsumerGroupSession) error {
	if c.cleanupFunc == nil {
		return nil
	}
	return c.cleanupFunc(session)
}

func (c *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if c.consumeclaimFunc == nil {
		return nil
	}
	return c.consumeclaimFunc(session, claim)
}

func NewHandler(setupFunc, cleanupFunc func(sarama.ConsumerGroupSession) error, consumeclaimFunc func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error) *Handler {
	return &Handler{
		setupFunc:        setupFunc,
		cleanupFunc:      cleanupFunc,
		consumeclaimFunc: consumeclaimFunc,
	}
}

// 构造消费者结构体
type ConsumerGroup struct {
	config        *sarama.Config       //配置
	consumerGroup sarama.ConsumerGroup //消费者组
	host          []string             //kafka的地址
	group         string               //消费者组的名称
	topic         []string             //消费的主题
	handler       *Handler             //消费者组处理函数(用于实现sarama.ConsumerGroupHandler接口)
}

func NewConsumerConfig(host []string, group string, topic []string, handler *Handler) Consumer {
	config := sarama.NewConfig()
	// 返回允许用户自己改配置
	return &ConsumerGroup{
		config:  config,
		handler: handler,
		host:    host,
		group:   group,
		topic:   topic,
	}
}

func (c *ConsumerGroup) GetConfig() *sarama.Config {
	return c.config
}

func (c *ConsumerGroup) CreateConsumerGroup() error {
	group, err := sarama.NewConsumerGroup(c.host, c.group, c.config)
	if err != nil {
		return err
	}
	c.consumerGroup = group
	return nil
}

// 会一直阻塞消费直到出错或取消上下文
func (c *ConsumerGroup) Consume(ctx context.Context) error {
	if err := c.consumerGroup.Consume(ctx, c.topic, c.handler); err != nil {
		return err
	}
	return nil
}

func (c *ConsumerGroup) Close() error {
	return c.consumerGroup.Close()
}

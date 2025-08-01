package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/fs/sonyflake"
	"github.com/farseer-go/fs/trace"
)

type registerEvent struct {
	eventName    string
	client       IClient
	traceManager trace.IManager
}

func (receiver *registerEvent) Publish(message any) error {
	var jsonContent string
	switch message.(type) {
	case string:
		jsonContent = message.(string)
	default:
		b, _ := snc.Marshal(message)
		jsonContent = string(b)
	}
	_, err := receiver.client.Publish(receiver.eventName, jsonContent)
	return err
}

func (receiver *registerEvent) PublishAsync(message any) {
	var jsonContent string
	switch msg := message.(type) {
	case string:
		jsonContent = msg
	default:
		b, _ := snc.Marshal(message)
		jsonContent = string(b)
	}
	go receiver.client.Publish(receiver.eventName, jsonContent)
}

type registerSubscribe struct {
	eventName string
	client    IClient
	consumers map[string]core.ConsumerFunc
}

// RegisterEvent 注册core.IEvent实现
func RegisterEvent(redisConfigName, eventName string) *registerSubscribe {
	redisClient := container.Resolve[IClient](redisConfigName)
	// 注册仓储
	container.Register(func() core.IEvent {
		return &registerEvent{
			eventName:    eventName,
			client:       redisClient,
			traceManager: container.Resolve[trace.IManager](),
		}
	}, eventName)

	sub := &registerSubscribe{
		eventName: eventName,
		client:    redisClient,
		consumers: make(map[string]core.ConsumerFunc),
	}
	go sub.subscribe()
	return sub
}

// RegisterSubscribe 注册订阅者
func (receiver *registerSubscribe) RegisterSubscribe(subscribeName string, consumerFunc core.ConsumerFunc) *registerSubscribe {
	if _, exists := receiver.consumers[subscribeName]; exists {
		panic("RegisterSubscribe已存在相同的订阅者名称：" + subscribeName)
	}
	receiver.consumers[subscribeName] = consumerFunc
	return receiver
}

func (receiver *registerSubscribe) subscribe() {
	server := fmt.Sprintf("redis订阅/%s", receiver.client.Original().String())
	for message := range receiver.client.Subscribe(receiver.eventName) {
		// InitContext 初始化同一协程上下文，避免在同一协程中多次初始化
		asyncLocal.InitContext()
		eventArgs := core.EventArgs{
			Id:         strconv.FormatInt(sonyflake.GenerateId(), 10),
			CreateAt:   time.Now().UnixMilli(),
			Message:    message.Payload,
			ErrorCount: 0,
			EventName:  message.Channel,
		}

		// 同时订阅消费
		for subscribeName, consumerFunc := range receiver.consumers {
			// 创建一个事件消费入口
			eventTraceContext := container.Resolve[trace.IManager]().EntryEventConsumer(server, receiver.eventName, subscribeName)
			exception.Try(func() {
				consumerFunc(message.Payload, eventArgs)
			})
			container.Resolve[trace.IManager]().Push(eventTraceContext, nil)
		}
		asyncLocal.Release()
	}
}

package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/color"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
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
	switch message := message.(type) {
	case string:
		jsonContent = message
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
	for {
		for message := range receiver.client.Subscribe(receiver.eventName) {
			eventArgs := core.EventArgs{
				Id:         strconv.FormatInt(sonyflake.GenerateId(), 10),
				CreateAt:   time.Now().UnixMilli(),
				Message:    message.Payload,
				ErrorCount: 0,
				EventName:  message.Channel,
			}

			// 同时订阅消费
			for subscribeName, consumerFunc := range receiver.consumers {
				// InitContext 初始化同一协程上下文，避免在同一协程中多次初始化
				asyncLocal.InitContext()
				// 创建一个事件消费入口
				traceContext := container.Resolve[trace.IManager]().EntryEventConsumer(server, receiver.eventName, subscribeName)
				exception.Try(func() {
					consumerFunc(message.Payload, eventArgs)
				}).CatchException(func(exp any) {
					if traceContext.IsIgnore() { // 如果忽略了链路,则要在这里打印错误日志
						lstLogs := []string{fmt.Sprintf("%s,%s 异常: %v", server, receiver.eventName, exp)}
						for index, exceptionStackDetail := range trace.GetCallerInfo() {
							lstLogs = append(lstLogs, fmt.Sprintf("\t%d、%s:%s %s", index+1, exceptionStackDetail.ExceptionCallFile, color.Yellow(exceptionStackDetail.ExceptionCallLine), color.Red(exceptionStackDetail.ExceptionCallFuncName)))
						}
						flog.Error(strings.Join(lstLogs, "\n") + "\n")
					}
				})
				container.Resolve[trace.IManager]().Push(traceContext, nil)
				asyncLocal.Release()
			}
		}
		// channel 关闭说明连接断开，等待后重新订阅
		flog.Warningf("%s,事件: %s 连接断开，3秒后重新订阅...", server, receiver.eventName)
		time.Sleep(3 * time.Second)
	}
}

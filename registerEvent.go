package redis

import (
	"encoding/json"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/snowflake"
	"github.com/farseer-go/fs/trace"
	"strconv"
	"time"
)

type registerEvent struct {
	eventName    string
	client       IClient
	traceManager trace.IManager
}

func (receiver *registerEvent) Publish(message any) error {
	traceDetail := receiver.traceManager.TraceRedis("Publish", receiver.eventName, "")
	var jsonContent string
	switch message.(type) {
	case string:
		jsonContent = message.(string)
	default:
		b, _ := json.Marshal(message)
		jsonContent = string(b)
	}
	_, err := receiver.client.Publish(receiver.eventName, jsonContent)
	defer func() { traceDetail.End(err) }()
	return err
}

// RegisterEvent 注册core.IEvent实现
func RegisterEvent(redisConfigName, eventName string, fns ...core.ConsumerFunc) {
	client := container.Resolve[IClient](redisConfigName)
	// 注册仓储
	container.Register(func() core.IEvent {
		return &registerEvent{
			eventName: eventName,
			client:    client,
		}
	}, eventName)

	go subscribe(client, eventName, fns)
}

func subscribe(client IClient, eventName string, fns []core.ConsumerFunc) {
	for message := range client.Subscribe(eventName) {
		eventArgs := core.EventArgs{
			Id:         strconv.FormatInt(snowflake.GenerateId(), 10),
			CreateAt:   time.Now().UnixMilli(),
			Message:    message.Payload,
			ErrorCount: 0,
			EventName:  message.Channel,
		}

		// 同时订阅消费
		for i := 0; i < len(fns); i++ {
			fns[i](message.Payload, eventArgs)
		}
	}
}

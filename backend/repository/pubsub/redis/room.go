package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"chat-app-demo/lib"
)

type Payload struct {
	LastMessageId int64 `json:"lastMsgId"`
}

func (p *Payload) ToJson() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (rps *RedisPubSub) PublishRoom(roomId string, msgId int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	payload := Payload{LastMessageId: msgId}
	_, err := rps.rdb.Publish(ctx, roomId, payload.ToJson()).Result()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}

func (rps *RedisPubSub) SubscribeRoom(roomId string, lastMsgId chan<- int64) (unsubscribe func(), err error) {
	sub := rps.rdb.Subscribe(context.Background(), roomId)

	unsubscribe = func() {
		sub.Close()
	}

	ch := sub.Channel()

	go func() {
		defer func() {
			sub.Close()

			// Catch panic in case this method sends a LastMessageId to closed channel.
			recover()
		}()

		for msg := range ch {
			var payload Payload
			err := json.Unmarshal([]byte(msg.Payload), &payload)
			if err != nil {
				fmt.Println(lib.ErrBuilder(err).Error())
				return
			}

			lastMsgId <- payload.LastMessageId
		}
	}()

	return unsubscribe, nil
}

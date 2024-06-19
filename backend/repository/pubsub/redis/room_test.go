package redis

import (
	"sync"
	"testing"
)

func TestRoom(t *testing.T) {
	defer DeleteAll()

	pubsub := NewRedisPubSub()

	lastMsgId := make(chan int64, 1)

	unsubscribe, err := pubsub.SubscribeRoom("r1", lastMsgId)
	if err != nil {
		t.Fatal(err)
	}
	defer unsubscribe()

	n := 5

	wg := sync.WaitGroup{}
	wg.Add(n)

	total := int64(0)
	go func() {
		for id := range lastMsgId {
			total += id
			wg.Done()
		}
	}()

	for i := range n {
		err := pubsub.PublishRoom("r1", int64(i))
		if err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()

	if total != 10 {
		t.Fatal("Expected 15, got", total)
	}

}

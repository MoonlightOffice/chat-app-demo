package pubsub

type IPubSub interface {
	PublishRoom(roomId string, msgId int64) error

	SubscribeRoom(roomId string, lastMsgId chan<- int64) (unsubscribe func(), err error)
}

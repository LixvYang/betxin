package mq

type Broker interface {
	publish(topic string, msg interface{}) error
	subscribe(topic string) (<-chan interface{}, error)
	unsubscribe(topic string, sub <-chan interface{}) error
	close()
	broadcast(msg interface{}, subscribers []chan interface{})
	setConditions(capacity int)
}

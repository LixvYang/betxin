package mq

type Client struct {
	bro *MQ
}

func NewMQClient() *Client {
	return &Client{
		bro: NewMQ(),
	}
}

func (c *Client) SetConditions(capacity int) {
	c.bro.setConditions(capacity)
}

func (c *Client) Publish(topic string, msg interface{}) error {
	return c.bro.publish(topic, msg)
}

func (c *Client) Subscribe(topic string) (<-chan interface{}, error) {
	return c.bro.subscribe(topic)
}

func (c *Client) Unsubscribe(topic string, sub <-chan interface{}) error {
	return c.bro.unsubscribe(topic, sub)
}

func (c *Client) Close() {
	c.bro.close()
}

func (c *Client) GetPayLoad(sub <-chan interface{}) interface{} {
	for val := range sub {
		if val != nil {
			return val
		}
	}
	return nil
}

func (c *Client) GetQueue(title string) []chan interface{} {
	return c.bro.topics[title]
}

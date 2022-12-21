package mq

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

type Message struct {
	ID   float64
	Name string
}

func TestClient(t *testing.T) {
	var m Message

	b := NewMQClient()
	b.SetConditions(100)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		topic := "Golang"

		m.ID = rand.Float64()
		m.Name = strconv.Itoa(rand.Int())

		payload := m
		ch, err := b.Subscribe(topic)
		if err != nil {
			t.Fatal(err)
		}

		wg.Add(1)
		go func() {
			e := b.GetPayLoad(ch).(Message)
			fmt.Println(e)
			if e != payload {
				t.Fatal(topic, " expected ", payload, " but get", e)
			}
			if err := b.Unsubscribe(topic, ch); err != nil {
				t.Fatal(err)
			}
			wg.Done()
		}()

		if err := b.Publish(topic, payload); err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()
}

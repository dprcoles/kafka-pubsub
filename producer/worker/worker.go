package worker

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"sync"
	"time"
)

type postAdded struct {
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Content    string    `json:"content"`
	OccurredOn time.Time `json:"occurred_on"`
}

func Create(messagesPerSecond int, publisher message.Publisher, wg *sync.WaitGroup, closeCh chan struct{}) {
	ticker := time.NewTicker(time.Duration(int(time.Second) / messagesPerSecond))

	for {
		select {
		case <-closeCh:
			ticker.Stop()
			wg.Done()
			return

		case <-ticker.C:
		}

		msgPayload := postAdded{
			OccurredOn: time.Now(),
			Author:     gofakeit.Username(),
			Title:      gofakeit.Sentence(rand.Intn(5) + 1),
			Content:    gofakeit.Sentence(rand.Intn(10) + 5),
		}

		payload, err := json.Marshal(msgPayload)
		if err != nil {
			panic(err)
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)

		middleware.SetCorrelationID(watermill.NewShortUUID(), msg)
		err = publisher.Publish("posts_published", msg)
		if err != nil {
			fmt.Println("cannot publish message:", err)
			continue
		}
	}
}

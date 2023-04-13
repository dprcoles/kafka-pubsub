package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"time"
)

type postAdded struct {
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	OccurredOn time.Time `json:"occurred_on"`
}

type feedStorage interface {
	AddToFeed(title, author string, time time.Time) error
}

type printFeedStorage struct{}

func (printFeedStorage) AddToFeed(title, author string, time time.Time) error {
	fmt.Printf("Adding to feed: %s by %s @%s\n", title, author, time)
	return nil
}

type FeedGenerator interface {
	UpdateFeed(message *message.Message) error
}

type feedGenerator struct {
	feedStorage feedStorage
}

func NewFeedGenerator() FeedGenerator {
	return &feedGenerator{printFeedStorage{}}
}

func (f *feedGenerator) UpdateFeed(message *message.Message) error {
	event := postAdded{}
	if err := json.Unmarshal(message.Payload, &event); err != nil {
		return err
	}

	err := f.feedStorage.AddToFeed(event.Title, event.Author, event.OccurredOn)
	if err != nil {
		return errors.Wrap(err, "cannot update feed")
	}

	return nil
}

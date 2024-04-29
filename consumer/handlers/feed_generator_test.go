package handlers

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockFeedStorage struct {
	mock.Mock
}

func (m *MockFeedStorage) AddToFeed(title, author string, time time.Time) error {
	args := m.Called(title, author, time)
	return args.Error(0)
}

func TestFeedGenerator_UpdateFeed(t *testing.T) {
	storage := new(MockFeedStorage)
	author := "author"
	title := "title"
	occurredOn := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	storage.On("AddToFeed", title, author, occurredOn).Return(nil)

	feed := &feedGenerator{storage}
	msg := &message.Message{Payload: []byte(`{"title":"` + title + `","author":"` + author + `","occurred_on":"` + occurredOn.Format(time.RFC3339) + `"}`)}
	err := feed.UpdateFeed(msg)

	assert.NoError(t, err)
	storage.AssertCalled(t, "AddToFeed", title, author, occurredOn)

	var payload postAdded
	err = json.Unmarshal(msg.Payload, &payload)
	assert.NoError(t, err)

	assert.Equal(t, title, payload.Title)
	assert.Equal(t, author, payload.Author)
	assert.Equal(t, occurredOn, payload.OccurredOn)
}

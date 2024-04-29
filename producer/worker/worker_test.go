package worker

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"time"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(topic string, messages ...*message.Message) error {
	args := m.Called(topic, messages)
	return args.Error(0)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreate(t *testing.T) {
	publisher := new(MockPublisher)
	wg := &sync.WaitGroup{}
	closeCh := make(chan struct{})
	messagesPerSecond := 1

	publisher.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]*message.Message")).Return(nil)

	wg.Add(1)
	go Create(messagesPerSecond, publisher, wg, closeCh)

	time.Sleep(1 * time.Second)

	// Simulate closing of the worker
	close(closeCh)
	wg.Wait()

	publisher.AssertCalled(t, "Publish", "posts_published", mock.AnythingOfType("[]*message.Message"))

	msg := publisher.Calls[0].Arguments.Get(1).([]*message.Message)[0]
	var payload postAdded
	err := json.Unmarshal(msg.Payload, &payload)
	assert.NoError(t, err)

	assert.NotEmpty(t, payload.Title)
	assert.NotEmpty(t, payload.Author)
	assert.NotEmpty(t, payload.Content)
	assert.NotZero(t, payload.OccurredOn)
}

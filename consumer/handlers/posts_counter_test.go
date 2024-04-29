package handlers

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockCountStorage struct {
	mock.Mock
}

func (m *MockCountStorage) CountAdd() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCountStorage) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestPostsCounter_Count(t *testing.T) {
	storage := new(MockCountStorage)
	storage.On("CountAdd").Return(int64(1), nil)

	counter := &postsCounter{storage}

	msg := message.NewMessage("1", []byte("dummy message"))
	returnedMessages, err := counter.Count(msg)

	assert.NoError(t, err)
	storage.AssertCalled(t, "CountAdd")

	assert.Len(t, returnedMessages, 1)

	var payload postsCountUpdated
	err = json.Unmarshal(returnedMessages[0].Payload, &payload)
	assert.NoError(t, err)

	assert.Equal(t, int64(1), payload.NewCount)
}

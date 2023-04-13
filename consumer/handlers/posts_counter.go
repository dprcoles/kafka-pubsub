package handlers

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"sync/atomic"
)

type postsCountUpdated struct {
	NewCount int64 `json:"new_count"`
}

type countStorage interface {
	CountAdd() (int64, error)
	Count() (int64, error)
}

type memoryCountStorage struct {
	count *int64
}

func (m memoryCountStorage) CountAdd() (int64, error) {
	return atomic.AddInt64(m.count, 1), nil
}

func (m memoryCountStorage) Count() (int64, error) {
	return atomic.LoadInt64(m.count), nil
}

type PostsCounter interface {
	Count(msg *message.Message) ([]*message.Message, error)
}

type postsCounter struct {
	countStorage countStorage
}

func NewPostsCounter() PostsCounter {
	return &postsCounter{memoryCountStorage{new(int64)}}
}

func (p *postsCounter) Count(msg *message.Message) ([]*message.Message, error) {
	newCount, err := p.countStorage.CountAdd()
	if err != nil {
		return nil, errors.Wrap(err, "cannot add count")
	}

	producedMsg := postsCountUpdated{NewCount: newCount}
	b, err := json.Marshal(producedMsg)
	if err != nil {
		return nil, err
	}

	return []*message.Message{message.NewMessage(watermill.NewUUID(), b)}, nil
}

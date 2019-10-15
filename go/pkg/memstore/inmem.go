package memstore

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

var (
	atomic_id uint64
)

// IdSorter sorts planets by name.
type IdSorter []Message

func (id IdSorter) Len() int           { return len(id) }
func (id IdSorter) Swap(i, j int)      { id[i], id[j] = id[j], id[i] }
func (id IdSorter) Less(i, j int) bool { return id[i].Id < id[j].Id }

type Message struct {
	Id   uint64 `json:"id,omitempty" form:"id" binding:"required"`
	Text string `json:"text" form:"message" binding:"required"`
}

type Store struct {
	Messages map[uint64]Message `json:"id" form:"id" binding:"required"`
	sync.RWMutex
}

func (s *Store) GetMessages() (error, []Message) {
	s.RLock()
	data := make([]Message, 0, len(s.Messages))
	for _, value := range s.Messages {
		data = append(data, value)
	}
	s.RUnlock()
	sort.Sort(IdSorter(data))
	return nil, data
}

func (s *Store) GetMessage(id uint64) (error, Message) {
	s.RLock()
	if val, ok := s.Messages[id]; ok {
		s.RUnlock()
		return nil, val
	} else {
		s.RUnlock()
		return fmt.Errorf("There is no message for id: %d", id), Message{}
	}
}

func (s *Store) UpdateMessage(id uint64, message string) error {
	var msg Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return err
	}
	msg.Id = id
	s.Lock()
	s.Messages[id] = msg
	s.Unlock()
	return nil
}

func (s *Store) NewMessage(message string) error {
	var msg Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return err
	}
	id := atomic.AddUint64(&atomic_id, 1)
	msg.Id = id
	s.Lock()
	s.Messages[id] = msg
	s.Unlock()
	return nil
}

func (s *Store) DeleteMessage(id uint64) error {
	s.Lock()
	delete(s.Messages, id)
	s.Unlock()
	return nil
}

package memstore

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/peterahl/storytel/go/pkg/models"
)

var (
	atomic_id uint64
)

// IdSorter sorts planets by name.
type IdSorter []models.Message

func (id IdSorter) Len() int           { return len(id) }
func (id IdSorter) Swap(i, j int)      { id[i], id[j] = id[j], id[i] }
func (id IdSorter) Less(i, j int) bool { return id[i].Id < id[j].Id }

type Store struct {
	Messages map[uint64]models.Message `json:"id" form:"id" binding:"required"`
	sync.RWMutex
}

func (s *Store) GetMessages() (error, []models.Message) {
	s.RLock()
	data := make([]models.Message, 0, len(s.Messages))
	for _, value := range s.Messages {
		data = append(data, value)
	}
	s.RUnlock()
	sort.Sort(IdSorter(data))
	return nil, data
}

func (s *Store) GetMessage(id uint64) (error, models.Message) {
	s.RLock()
	if val, ok := s.Messages[id]; ok {
		s.RUnlock()
		return nil, val
	} else {
		s.RUnlock()
		return fmt.Errorf("There is no message for id: %d", id), models.Message{}
	}
}

func (s *Store) UpdateMessage(msg models.Message) error {
	s.Lock()
	s.Messages[msg.GetId()] = msg
	s.Unlock()
	return nil
}

func (s *Store) NewMessage(msg models.Message) error {
	id := atomic.AddUint64(&atomic_id, 1)
	msg.Id = id
	s.Lock()
	s.Messages[id] = msg
	s.Unlock()
	return nil
}

func (s *Store) DeleteMessage(msg models.Message) error {
	s.Lock()
	delete(s.Messages, msg.Id)
	s.Unlock()
	return nil
}

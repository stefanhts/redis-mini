package data

import (
	"errors"
	"sync"
)

type Store struct {
	lists map[string][]string
	lock  sync.Mutex
}

func NewStore() *Store {
	return &Store{lists: make(map[string][]string)}
}

func (s *Store) LLen(key string) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	list, ok := s.lists[key]
	if !ok {
		return 0
	}

	return len(list)
}

func (s *Store) LPush(key string, values ...string) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	list := s.lists[key]
	for _, value := range values {
		list = append([]string{value}, list...)
	}
	s.lists[key] = list

	return len(values)
}

func (s *Store) LPop(key string, num int64) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	list, ok := s.lists[key]
	if !ok {
		return []string{}, errors.New("key not found")
	}

	if len(list) == 0 {
		return []string{}, nil
	}

	if int(num) > len(list) {
		num = int64(len(list))
	}

	values := list[:num]
	s.lists[key] = list[num:]
	return values, nil
}

func (s *Store) LPos(key, value string) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	list, ok := s.lists[key]
	if !ok {
		return -1, errors.New("key not found")
	}

	for i, v := range list {
		if v == value {
			return i, nil
		}
	}

	return -1, nil
}

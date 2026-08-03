package user

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotFound     = errors.New("user not found")
	ErrNameConflict = errors.New("user name already exists")
)

type Store interface {
	Get(ctx context.Context, id ID) (*User, error)
	List(ctx context.Context) []User
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id ID) (*User, error)
}

type user struct {
	Name     string
	Password string
}

type MemoryStore struct {
	users  map[ID]user
	nextID ID
	mu     sync.RWMutex
}

func NewMemoryStore() Store {
	return &MemoryStore{
		users: map[ID]user{
			1: {Name: "admin", Password: "password"},
			2: {Name: "tito", Password: "password"},
		},
		nextID: 3,
	}
}

func (s *MemoryStore) isNameConflict(name string) error {
	for _, user := range s.users {
		if name == user.Name {
			return ErrNameConflict
		}
	}
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, id ID) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &User{
		ID:       id,
		Name:     user.Name,
		Password: user.Password,
	}, nil
}

func (s *MemoryStore) List(ctx context.Context) []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))

	for id, user := range s.users {
		users = append(users, User{
			ID:       id,
			Name:     user.Name,
			Password: user.Password,
		})
	}

	return users
}

func (s *MemoryStore) Create(ctx context.Context, _user *User) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.isNameConflict(_user.Name); err != nil {
		return nil, err
	}

	id := s.nextID
	s.nextID++

	s.users[id] = user{
		Name:     _user.Name,
		Password: _user.Password,
	}

	createUser := *_user
	createUser.ID = id

	return &createUser, nil
}

func (s *MemoryStore) Update(ctx context.Context, _user *User) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldUser, ok := s.users[_user.ID]
	if !ok {
		return nil, ErrNotFound
	}

	if _user.Name != oldUser.Name {
		if err := s.isNameConflict(_user.Name); err != nil {
			return nil, err
		}
	}

	s.users[_user.ID] = user{
		Name:     _user.Name,
		Password: _user.Password,
	}

	return _user, nil
}

func (s *MemoryStore) Delete(ctx context.Context, id ID) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}

	delete(s.users, id)

	return &User{
		ID:       id,
		Name:     user.Name,
		Password: user.Password,
	}, nil
}

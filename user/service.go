package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) hashUserPassword(user User) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hash)

	return &user, nil
}

func (s *Service) List(ctx context.Context) []User {
	return s.store.List(ctx)
}

func (s *Service) Get(ctx context.Context, id ID) (*User, error) {
	return s.store.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, user *User) (*User, error) {
	securedUser, err := s.hashUserPassword(*user)
	if err != nil {
		return nil, err
	}

	return s.store.Create(ctx, securedUser)
}

func (s *Service) Update(ctx context.Context, user *User) (*User, error) {
	securedUser, err := s.hashUserPassword(*user)
	if err != nil {
		return nil, err
	}

	return s.store.Update(ctx, securedUser)
}

func (s *Service) Delete(ctx context.Context, id ID) (*User, error) {
	return s.store.Delete(ctx, id)
}

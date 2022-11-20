package user

import "context"

type storage interface {
	InsertUser(ctx context.Context, user *User) (string, error)
	FindUser(ctx context.Context, id string) (*User, error)
	InsertTransaction(ctx context.Context, ur *UserRequest) error
}

type Service struct {
	storage storage
}

func NewService(storage storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateUser(ctx context.Context, user *User) (string, error) {
	return "", nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func (s *Service) MakeTransaction(ctx context.Context, ur *UserRequest) error {

	return nil
}

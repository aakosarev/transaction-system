package user

import "context"

type storage interface {
	InsertUser(ctx context.Context, user *User) (string, error)
	FindUser(ctx context.Context, id string) (*User, error)
	InsertTransaction(ctx context.Context, trr *TransactionRequest) error
}

type Service struct {
	storage storage
}

func NewService(storage storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateUser(ctx context.Context, user *User) (string, error) {
	userID, err := s.storage.InsertUser(ctx, user)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	user, err := s.storage.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) MakeTransaction(ctx context.Context, trr *TransactionRequest) error {
	err := s.storage.InsertTransaction(ctx, trr)
	if err != nil {
		return err
	}
	return nil
}

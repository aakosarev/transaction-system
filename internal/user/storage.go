package user

import (
	"context"
	"github.com/aakosarev/transaction-system/pkg/client/postgresql"
)

type Storage struct {
	client postgresql.Client
}

func NewStorage(client postgresql.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) InsertUser(ctx context.Context, user *User) (string, error) {
	return "", nil
}

func (s *Storage) FindUser(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func (s *Storage) InsertTransaction(ctx context.Context, ur *UserRequest) error {
	return nil
}

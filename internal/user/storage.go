package user

import (
	"context"
	"errors"
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
	query := `INSERT INTO users (name) VALUES ($1) RETURNING id`
	err := s.client.QueryRow(ctx, query, user.ID).Scan(&user.ID)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (s *Storage) FindUser(ctx context.Context, id string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE id = $1`
	err := s.client.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Balance)
	if err != nil {
		return nil, errors.New("user not found")
	}
	// transfer the balance to rubles ?
	return &user, nil
}

func (s *Storage) InsertTransaction(ctx context.Context, trr *TransactionRequest) error {
	query := `INSERT INTO transactions (user_id, tr_type, amount) VALUES ($1, $2, $3)`
	_, err := s.client.Exec(ctx, query, trr.UserID, trr.TransactionType, trr.Amount)
	if err != nil {
		return err
	}
	return nil
}

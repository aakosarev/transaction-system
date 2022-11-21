package user

import (
	"context"
	"log"
	"time"
)

type storage interface {
	InsertUser(ctx context.Context, user *User) (string, error)
	FindUser(ctx context.Context, id string) (*User, error)
	InsertTransaction(ctx context.Context, trr *TransactionRequest) error
	HandleOneTransaction(ctx context.Context, tr *UserTransaction, u *User) error
	FindOpenTransactionsForUser(ctx context.Context, u *User) ([]*UserTransaction, error)
	FindAllUsers(ctx context.Context) ([]*User, error)
	GetUserBalance(ctx context.Context, userID string) (int, error)
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

func (s *Service) HandleTransactions(ctx context.Context, user *User) {
	for {
		openTransactions, err := s.storage.FindOpenTransactionsForUser(ctx, user)
		if err != nil {
			log.Fatal(err)
		}
		for _, openTransaction := range openTransactions {
			err = s.storage.HandleOneTransaction(ctx, openTransaction, user)
			if err != nil {
				log.Fatal(err)
			}
			newBalance, err := s.storage.GetUserBalance(ctx, user.ID)
			if err != nil {
				log.Fatal(err)
			}
			user.Balance = newBalance
			time.Sleep(5 * time.Second)
		}
	}
}

func (s *Service) GetAllUsers(ctx context.Context) ([]*User, error) {
	users, err := s.storage.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

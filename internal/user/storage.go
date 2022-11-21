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
	err := s.client.QueryRow(ctx, query, user.Name).Scan(&user.ID)
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

func (s *Storage) HandleOneTransaction(ctx context.Context, tr *UserTransaction, u *User) error {
	if tr.Amount > 0 && tr.Type == "replenishment" {
		newBalance := tr.Amount + u.Balance
		tx, err := s.client.Begin(ctx)
		defer tx.Rollback(ctx)
		if err != nil {
			return err
		}
		query := `UPDATE users SET balance = $1 WHERE id = $2`
		_, err = tx.Exec(ctx, query, newBalance, u.ID)
		if err != nil {
			return err
		}
		query = `UPDATE transactions SET tr_status = $1 WHERE id = $2`
		_, err = tx.Exec(ctx, query, "closed", tr.ID)
		if err != nil {
			return err
		}
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	} else if tr.Amount > 0 && tr.Amount <= u.Balance && tr.Type == "write-off" {
		newBalance := u.Balance - tr.Amount
		tx, err := s.client.Begin(ctx)
		defer tx.Rollback(ctx)
		if err != nil {
			return err
		}
		query := `UPDATE users SET balance = $1 WHERE id = $2`
		_, err = tx.Exec(ctx, query, newBalance, u.ID)
		if err != nil {
			return err
		}
		query = `UPDATE transactions SET tr_status = $1 WHERE id = $2`
		_, err = tx.Exec(ctx, query, "closed", tr.ID)
		if err != nil {
			return err
		}
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	} else {
		query := `UPDATE transactions SET tr_status = $1 WHERE id = $2`
		_, err := s.client.Exec(ctx, query, "rejected", tr.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *Storage) FindOpenTransactionsForUser(ctx context.Context, u *User) ([]*UserTransaction, error) {
	var transactions []*UserTransaction
	query := `SELECT * FROM transactions WHERE user_id = $1 AND tr_status = $2 ORDER BY created_at`
	transactionsDB, err := s.client.Query(ctx, query, u.ID, "open")
	if err != nil {
		return nil, err
	}
	for transactionsDB.Next() {
		transaction := &UserTransaction{}
		err = transactionsDB.Scan(&transaction.ID, &transaction.UserID, &transaction.Type, &transaction.Status, &transaction.Amount, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *Storage) FindAllUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	query := `SELECT * FROM users`
	usersDB, err := s.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for usersDB.Next() {
		user := &User{}
		err = usersDB.Scan(&user.ID, &user.Name, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) GetUserBalance(ctx context.Context, userID string) (int, error) {
	query := `SELECT balance FROM users WHERE id = $1`
	var userBalance int
	err := s.client.QueryRow(ctx, query, userID).Scan(&userBalance)
	if err != nil {
		return -1, err
	}
	return userBalance, nil
}

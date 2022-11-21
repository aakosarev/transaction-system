package user

import (
	"context"
	"log"
)

var _users []*User

type Worker struct {
	service service
}

func NewWorker(service service) *Worker {
	return &Worker{
		service: service,
	}
}

func (w *Worker) StartTransactionProcessing(ctx context.Context) {
	_users, err := w.service.GetAllUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range _users {
		go w.service.HandleTransactions(ctx, user)
	}
}

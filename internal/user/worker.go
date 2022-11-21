package user

import (
	"context"
	"fmt"
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
		log.Fatal("here _1", err)
	}
	for _, user := range _users {
		fmt.Println("q", user)
		go w.service.HandleTransactions(ctx, user)
	}
}

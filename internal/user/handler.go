package user

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type service interface {
	CreateUser(ctx context.Context, user *User) (string, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	MakeTransaction(ctx context.Context, ur *UserRequest) error
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/user/new", h.CreateUser)
	router.HandlerFunc(http.MethodGet, "/user/:uuid", h.GetUserByID)
	router.HandlerFunc(http.MethodPost, "/user/transaction", h.MakeTransaction)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) MakeTransaction(w http.ResponseWriter, r *http.Request) {

}

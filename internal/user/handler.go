package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _validator = validator.New()

type service interface {
	CreateUser(ctx context.Context, user *User) (string, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	MakeTransaction(ctx context.Context, trr *TransactionRequest) error
	HandleTransactions(ctx context.Context, user *User)
	GetAllUsers(ctx context.Context) ([]*User, error)
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
	router.GET("/user/:uuid", h.GetUserByID)
	router.HandlerFunc(http.MethodPost, "/user/transaction", h.MakeTransaction)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	errValidation := _validator.Struct(user)
	if errValidation != nil || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}
	userID, err := h.service.CreateUser(r.Context(), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"user_id":"%s"}`, userID)))
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	uuid := params.ByName("uuid")

	user, err := h.service.GetUserByID(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func (h *Handler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	trRequest := TransactionRequest{}
	err := json.NewDecoder(r.Body).Decode(&trRequest)
	errValidation := _validator.Struct(trRequest)
	if errValidation != nil || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}
	err = h.service.MakeTransaction(r.Context(), &trRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"The transaction is queued"}`))
}

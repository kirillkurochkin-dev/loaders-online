package handler

import (
	"context"
	"github.com/gorilla/mux"
	"loaders-online/internal/entity/dto"
	"net/http"
)

type UserService interface {
	Register(ctx context.Context, user *dto.UserSignInDto) error
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	public := r.PathPrefix("/api").Subrouter()
	{
		//public
		public.HandleFunc("/register", h.register).Methods(http.MethodPost)
	}

	return r
}

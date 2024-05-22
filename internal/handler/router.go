package handler

import (
	"context"
	"github.com/gorilla/mux"
	"loaders-online/internal/entity/dto"
	"net/http"
)

type TaskService interface {
	CreateTask(ctx context.Context, taskCr *dto.CreateTaskDto) error
}

type UserService interface {
	Register(ctx context.Context, user *dto.UserSignUpDto) error
	Login(ctx context.Context, user *dto.UserSignInDto) (string, error)
}

type Handler struct {
	userService UserService
	taskService TaskService
}

func NewHandler(userService UserService, taskService TaskService) *Handler {
	return &Handler{
		userService: userService,
		taskService: taskService,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	public := r.PathPrefix("/api").Subrouter()
	{
		//public
		public.HandleFunc("/register", h.register).Methods(http.MethodPost)
		public.HandleFunc("/login", h.login).Methods(http.MethodPost)
		public.HandleFunc("/tasks", h.tasks).Methods(http.MethodPost)
	}

	return r
}

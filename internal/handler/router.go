package handler

import (
	"context"
	"github.com/gorilla/mux"
	"loaders-online/internal/entity/dto"
	"loaders-online/internal/handler/middleware"
	"net/http"
)

type TaskService interface {
	CreateTask(ctx context.Context, taskCr *dto.CreateTaskDto) error
	GetCompletedTasks(ctx context.Context, id int) ([]dto.TaskCompletedDto, error)
	GetUncompletedTasks(ctx context.Context, id int) ([]dto.TaskUncompletedDto, error)
	AssignTasks(ctx context.Context, taskId int, loaderId int) error
	UpdateTask(ctx context.Context, task *dto.TaskUncompletedDto) error
}

type UserService interface {
	Register(ctx context.Context, user *dto.UserSignUpDto) error
	Login(ctx context.Context, user *dto.UserSignInDto) (string, error)
	CreateCustomer(ctx context.Context, customer *dto.CustomerSignUpDto) error
	GetCustomerById(ctx context.Context, id int) (dto.CustomerOutputDto, error)
	GetAssignedLoaders(ctx context.Context, id int) ([]dto.LoaderOutputDto, error)
	GetLoaderById(ctx context.Context, id int) (*dto.LoaderOutputDto, error)
	UpdateLoader(ctx context.Context, loader *dto.LoaderOutputDto) error
	UpdateCustomer(ctx context.Context, outputDto dto.CustomerUpdateDto) error
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
	r.Use()
	r.Use(middleware.LoggingMiddleware)

	public := r.PathPrefix("/api").Subrouter()
	{
		//public
		public.HandleFunc("/register", h.register).Methods(http.MethodPost)
		public.HandleFunc("/login", h.login).Methods(http.MethodPost)
		public.HandleFunc("/tasks", h.tasksPublic).Methods(http.MethodPost)

		//protected (customer, loader)
		customerLoaderOnly := public.PathPrefix("").Subrouter()
		customerLoaderOnly.Use(middleware.JWTMiddleware)
		customerLoaderOnly.Use(middleware.RoleMiddleware("customer", "loader"))
		customerLoaderOnly.HandleFunc("/me", h.me).Methods(http.MethodGet)
		customerLoaderOnly.HandleFunc("/tasks", h.tasks).Methods(http.MethodGet)

		customerOnly := public.PathPrefix("").Subrouter()
		customerOnly.Use(middleware.JWTMiddleware)
		customerOnly.Use(middleware.RoleMiddleware("customer"))
		customerOnly.HandleFunc("/start", h.start).Methods(http.MethodPost)
	}

	return r
}

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"loaders-online/internal/entity/dto"
	"loaders-online/pkg/game"
	"loaders-online/pkg/util"
	"net/http"
	"strconv"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var user dto.UserSignUpDto
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.LogHandler("register", "error decoding body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userService.Register(r.Context(), &user)
	if err != nil {
		util.LogHandler("register", "error registering user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var user dto.UserSignInDto
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.LogHandler("login", "error decoding body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.userService.Login(r.Context(), &user)

	if err != nil {
		util.LogHandler("login", "error login", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	id, role, err := getAuth(r)
	if err != nil {
		util.LogHandler("me", "error getting auth data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch role {
	case "customer":
		customer, err := h.userService.GetCustomerById(r.Context(), id)
		if err != nil {
			util.LogHandler("me", "error getting customer by id", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		customer.RegisteredLoaders, err = h.userService.GetAssignedLoaders(r.Context(), id)

		fmt.Println(customer.RegisteredLoaders)

		b, err := json.Marshal(customer)
		if err != nil {
			util.LogHandler("me", "error marshalling customer", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
		w.WriteHeader(http.StatusOK)
	case "loader":
		loader, err := h.userService.GetLoaderById(r.Context(), id)
		if err != nil {
			util.LogHandler("me", "error getting loader by id", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(loader)
		if err != nil {
			util.LogHandler("me", "error marshalling loader", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
		w.WriteHeader(http.StatusOK)
	default:
		util.LogHandler("me", "role not found in context", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) tasks(w http.ResponseWriter, r *http.Request) {
	id, role, err := getAuth(r)
	if err != nil {
		util.LogHandler("me", "error getting auth data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch role {
	case "customer":
		tasks, err := h.taskService.GetUncompletedTasks(r.Context(), id)
		if err != nil {
			util.LogHandler("me", "error getting uncompleted tasks", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(tasks)
		if err != nil {
			util.LogHandler("me", "error marshalling tasks", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
		w.WriteHeader(http.StatusOK)
	case "loader":
		tasks, err := h.taskService.GetCompletedTasks(r.Context(), id)
		if err != nil {
			util.LogHandler("me", "error getting completed tasks", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(tasks)
		if err != nil {
			util.LogHandler("me", "error marshalling tasks", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
		w.WriteHeader(http.StatusOK)
	default:
		util.LogHandler("tasks", "role not found in context", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) start(w http.ResponseWriter, r *http.Request) {
	id, _, err := getAuth(r)
	if err != nil {
		util.LogHandler("start", "error getting auth", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tasks, err := h.taskService.GetUncompletedTasks(r.Context(), id)
	if err != nil {
		util.LogHandler("start", "error getting uncompleted tasks", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var loaderInputIdDto dto.LoaderInputIdDto
	if err := json.NewDecoder(r.Body).Decode(&loaderInputIdDto); err != nil {
		util.LogHandler("start", "error decoding body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loaders, err := h.getLoadersByIds(r.Context(), loaderInputIdDto.Loaders)
	if err != nil {
		util.LogHandler("start", "error getting loaders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	customer, err := h.userService.GetCustomerById(r.Context(), id)
	if err != nil {
		util.LogHandler("start", "error getting customer", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = game.CalcMoney(loaders, &customer)
	if err != nil {
		util.LogHandler("start", "error calculating money", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateCustomer(r.Context(), dto.CustomerUpdateDto{
		CustomerID:     id,
		CurrentCapital: customer.CurrentCapital,
	})
	if err != nil {
		util.LogHandler("start", "error updating customer", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := h.processTasksAndLoaders(r.Context(), tasks, loaders)
	if err != nil {
		util.LogHandler("start", "error processing tasks and loaders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{Result: result}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		util.LogHandler("start", "error marshalling response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *Handler) getLoadersByIds(ctx context.Context, loaderIds []int) ([]dto.LoaderOutputDto, error) {
	var loaders []dto.LoaderOutputDto
	for _, id := range loaderIds {
		loader, err := h.userService.GetLoaderById(ctx, id)
		if err != nil {
			return nil, err
		}
		loaders = append(loaders, *loader)
	}
	return loaders, nil
}

func (h *Handler) processTasksAndLoaders(ctx context.Context, tasks []dto.TaskUncompletedDto, loaders []dto.LoaderOutputDto) (string, error) {
	flag := false

	for len(tasks) > 0 && len(loaders) > 0 {
		for i := 0; i < len(loaders); i++ {
			game.Recalculate(&loaders[i])
			if loaders[i].Fatigue != 100 {
				taskIndex := 0
				tasks[taskIndex].Weight -= loaders[i].MaxWeight
				if tasks[taskIndex].Weight < 0 {
					tasks[taskIndex].Weight = 0
				}

				err := h.taskService.AssignTasks(ctx, tasks[taskIndex].TaskID, loaders[i].LoaderID)
				if err != nil {
					return "", err
				}

				err = h.taskService.UpdateTask(ctx, &tasks[taskIndex])
				if err != nil {
					return "", err
				}

				if tasks[taskIndex].Weight <= 0 {
					tasks[taskIndex].Completed = true
					tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)
				}
			}

			loaders[i] = game.DoJob(loaders[i])
			err := h.userService.UpdateLoader(ctx, &loaders[i])
			if err != nil {
				return "", err
			}

			if loaders[i].Fatigue >= 100 {
				loaders = append(loaders[:i], loaders[i+1:]...)
				i--
			}
		}

		if len(tasks) == 0 {
			flag = true
			break
		}
	}

	if flag {
		return "win", nil
	}
	return "lose", nil
}

type Response struct {
	Result string `json:"result"`
}

func getAuth(r *http.Request) (int, string, error) {
	userID, ok := r.Context().Value("userID").(string)
	id, err := strconv.Atoi(userID)
	if !ok || id == 0 {
		return 0, "", errors.New("error getting user id")
	}
	role, ok := r.Context().Value("role").(string)
	if !ok || role == "" {
		return 0, "", errors.New("error getting user role")
	}
	if err != nil {
		return 0, "", errors.New("error getting auth data")
	}

	return id, role, nil
}

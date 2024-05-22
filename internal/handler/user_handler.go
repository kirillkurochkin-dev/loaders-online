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
	id, role, err := getAuth(w, r)
	if err != nil {
		util.LogHandler("me", "error getting auth data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch role {
	case "customer":
		customer, err := h.userService.GetCustomerById(context.Background(), id)
		if err != nil {
			util.LogHandler("me", "error getting customer by id", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		customer.RegisteredLoaders, err = h.userService.GetAssignedLoaders(context.Background(), id)

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
		loader, err := h.userService.GetLoaderById(context.Background(), id)
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
	id, role, err := getAuth(w, r)
	if err != nil {
		util.LogHandler("me", "error getting auth data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch role {
	case "customer":
		tasks, err := h.taskService.GetUncompletedTasks(context.Background(), id)
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
		tasks, err := h.taskService.GetCompletedTasks(context.Background(), id)
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
	id, _, err := getAuth(w, r)
	tasks, err := h.taskService.GetUncompletedTasks(context.Background(), id)
	if err != nil {
		util.LogHandler("start", "error getting uncompleted tasks", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var loaderInputIdDto dto.LoaderInputIdDto
	err = json.NewDecoder(r.Body).Decode(&loaderInputIdDto)
	if err != nil {
		util.LogHandler("start", "error decoding body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	loadersIds := loaderInputIdDto.Loaders

	var loaders []dto.LoaderOutputDto
	for i := 0; i < len(loadersIds); i++ {
		loader, err := h.userService.GetLoaderById(r.Context(), loadersIds[i])
		if err != nil {
			util.LogHandler("start", "error getting loader", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loaders = append(loaders, *loader)
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

	flag := false

	j := 0
	for {
		if len(tasks) == 0 {
			flag = true
			break
		}
		if len(loaders) == 0 {
			break
		}

		for i := 0; i < len(loaders); i++ {
			game.Recalculate(&loaders[i])
			if loaders[i].Fatigue != 100 {
				tasks[j].Weight -= loaders[i].MaxWeight
				if tasks[j].Weight < 0 {
					tasks[j].Weight = 0
				}
				fmt.Println(tasks[j].Weight)
				err := h.taskService.AssignTasks(context.Background(), tasks[j].TaskID, loaders[i].LoaderID)
				if err != nil {
					util.LogHandler("start", "error assigning tasks", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				err = h.taskService.UpdateTask(context.Background(), &tasks[j])
				if err != nil {
					util.LogHandler("start", "error updating task", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if tasks[j].Weight <= 0 {
					tasks[j].Completed = true

					tasks = append(tasks[:j], tasks[j+1:]...)
					//j++ не делаем, само сдвинется
				}
			}
			loaders[i] = game.DoJob(loaders[i])
			err = h.userService.UpdateLoader(context.Background(), &loaders[i])
			if err != nil {
				util.LogHandler("start", "error updating loader", err)
				return
			}
			if loaders[i].Fatigue >= 100 {
				loaders = append(loaders[:i], loaders[i+1:]...)
				i--
			}
		}
	}

	fmt.Println(flag)
}

func getAuth(w http.ResponseWriter, r *http.Request) (int, string, error) {
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

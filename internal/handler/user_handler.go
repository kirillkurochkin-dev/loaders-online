package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"loaders-online/internal/entity/dto"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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

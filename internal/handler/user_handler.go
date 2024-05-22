package handler

import (
	"encoding/json"
	"loaders-online/internal/entity/dto"
	"loaders-online/pkg/util"
	"net/http"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var user dto.UserSignInDto
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
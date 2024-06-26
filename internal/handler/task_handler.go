package handler

import (
	"encoding/json"
	"loaders-online/internal/entity/dto"
	"loaders-online/pkg/util"
	"net/http"
)

func (h *Handler) tasksPublic(w http.ResponseWriter, r *http.Request) {
	var task dto.CreateTaskDto
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		util.LogHandler("tasksPublic", "json decode error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.taskService.CreateTask(r.Context(), &task)
	if err != nil {
		util.LogHandler("tasksPublic", "create task error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

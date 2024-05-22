package repository

import (
	"context"
	"database/sql"
	"loaders-online/internal/entity/dto"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (t *TaskRepository) CreateTask(ctx context.Context, task *dto.TaskGeneratedDto) error {
	_, err := t.db.ExecContext(ctx, "INSERT INTO tasks (customer_id, task_name, weight, completed) VALUES ($1, $2, $3, $4)",
		task.CustomerID, task.TaskName, task.Weight, false)
	return err
}

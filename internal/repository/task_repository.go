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

func (t *TaskRepository) GetUncompletedTasksForCustomer(ctx context.Context, id int) ([]dto.TaskUncompletedDto, error) {
	var tasks []dto.TaskUncompletedDto
	rows, err := t.db.QueryContext(ctx, "SELECT task_id, task_name, weight, completed FROM tasks WHERE completed = false AND customer_id = $1", id)
	if err != nil {
		return tasks, err
	}
	for rows.Next() {
		var task dto.TaskUncompletedDto
		err := rows.Scan(&task.TaskID, &task.TaskName, &task.Weight, &task.Completed)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (t *TaskRepository) GetCompletedTasksForLoader(ctx context.Context, id int) ([]dto.TaskCompletedDto, error) {
	var tasks []dto.TaskCompletedDto
	rows, err := t.db.QueryContext(ctx, "select t.task_id, t.customer_id, t.task_name, t.weight, t.completed from tasks as t inner join loaders_tasks as lt on t.task_id = lt.task_id where lt.loader_id = $1", id)
	if err != nil {
		return tasks, err
	}
	for rows.Next() {
		var task dto.TaskCompletedDto
		err := rows.Scan(&task.TaskID, &task.CustomerID, &task.TaskName, &task.Weight, &task.Completed)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

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

func (t *TaskRepository) UpdateTask(ctx context.Context, task *dto.TaskUncompletedDto) error {
	_, err := t.db.ExecContext(ctx, "UPDATE tasks SET task_name = $1, weight = $2, completed = $3 WHERE task_id = $4",
		task.TaskName, task.Weight, task.Completed, task.TaskID)
	return err
}

func (t *TaskRepository) AssignTask(ctx context.Context, taskId int, loaderId int) error {
	_, err := t.db.ExecContext(ctx, "INSERT INTO loaders_tasks (loader_id, task_id) VALUES ($1, $2) ON CONFLICT (loader_id, task_id) DO NOTHING",
		loaderId, taskId)
	return err
}

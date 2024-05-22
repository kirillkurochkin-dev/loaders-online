package service

import (
	"context"
	"loaders-online/internal/entity/dto"
	utils "loaders-online/pkg/util"
	"math/rand"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *dto.TaskGeneratedDto) error
	GetCompletedTasksForLoader(ctx context.Context, id int) ([]dto.TaskCompletedDto, error)
	GetUncompletedTasksForCustomer(ctx context.Context, id int) ([]dto.TaskUncompletedDto, error)
}

type TaskService struct {
	taskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) CreateTask(ctx context.Context, taskCr *dto.CreateTaskDto) error {
	var tasks []dto.TaskGeneratedDto
	for i := 0; i < taskCr.Count; i++ {
		nameId := rand.Intn(24)
		weight := float64(rand.Intn(71) + 10)
		task := dto.TaskGeneratedDto{
			TaskName:   utils.Items[nameId],
			CustomerID: taskCr.CustomerID,
			Weight:     weight,
		}
		err := s.taskRepository.CreateTask(ctx, &task)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}
	return nil

}

func (t *TaskService) GetUncompletedTasks(ctx context.Context, id int) ([]dto.TaskUncompletedDto, error) {
	return t.taskRepository.GetUncompletedTasksForCustomer(ctx, id)
}

func (t *TaskService) GetCompletedTasks(ctx context.Context, id int) ([]dto.TaskCompletedDto, error) {
	return t.taskRepository.GetCompletedTasksForLoader(ctx, id)
}

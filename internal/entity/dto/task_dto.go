package dto

type TaskGeneratedDto struct {
	TaskName   string  `json:"task_name" validate:"required,max=100"`
	CustomerID int     `json:"customer_id" validate:"required"`
	Weight     float64 `json:"weight" validate:"required,gte=0,lte=80"`
}

type CreateTaskDto struct {
	Count      int `json:"count" validate:"required"`
	CustomerID int `json:"customer_id" validate:"required"`
}

type TaskUncompletedDto struct {
	TaskID    int     `json:"task_id" validate:"required"`
	TaskName  string  `json:"task_name" validate:"required,max=100"`
	Weight    float64 `json:"weight" validate:"required,gte=0,lte=80"`
	Completed bool    `json:"completed"`
}

type TaskCompletedDto struct {
	TaskID     int     `json:"task_id" validate:"required"`
	CustomerID int     `json:"customer_id" validate:"required"`
	TaskName   string  `json:"task_name" validate:"required,max=100"`
	Weight     float64 `json:"weight" validate:"required,gte=0,lte=80"`
	Completed  bool    `json:"completed"`
}

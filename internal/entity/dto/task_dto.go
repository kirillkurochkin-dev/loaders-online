package dto

type TaskGeneratedDto struct {
	TaskName   string  `json:"task_name" validate:"required,max=100"`
	CustomerID int     `json:"customer_id" validate:"required"`
	Weight     float64 `json:"weight" validate:"required,gte=10,lte=80"`
}

type CreateTaskDto struct {
	Count      int `json:"count" validate:"required"`
	CustomerID int `json:"customer_id" validate:"required"`
}

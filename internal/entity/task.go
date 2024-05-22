package entity

type Task struct {
	TaskID     int     `json:"task_id" validate:"required"`
	CustomerID int     `json:"customer_id" validate:"required"`
	TaskName   string  `json:"task_name" validate:"required,max=100"`
	Weight     float64 `json:"weight" validate:"required,gte=0,lte=80"`
	Completed  bool    `json:"completed"`
}

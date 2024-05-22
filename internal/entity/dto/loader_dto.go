package dto

type LoaderOutputDto struct {
	LoaderID    int     `json:"loader_id" validate:"required"`
	MaxWeight   float64 `json:"max_weight" validate:"required,gte=5,lte=30"`
	Drunkenness bool    `json:"drunkenness" validate:"required"`
	Fatigue     float64 `json:"fatigue" validate:"required,gte=0,lte=100"`
	Salary      float64 `json:"salary" validate:"required,gte=10000,lte=30000"`
}

type LoaderInputIdDto struct {
	Loaders []int `json:"loaders"`
}

/*
{
	"loaders": [2, 3, 4, 5, 6, 7, 8, 9, 10]
}
*/

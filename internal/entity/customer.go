package entity

type Customer struct {
	CustomerID      int     `json:"customer_id" validate:"required"`
	StartingCapital float64 `json:"starting_capital" validate:"required,gte=10000,lte=100000"`
	CurrentCapital  float64 `json:"current_capital" validate:"required,gte=10000,lte=100000"`
}

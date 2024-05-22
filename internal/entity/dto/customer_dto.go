package dto

type CustomerSignUpDto struct {
	CustomerID      int     `json:"customer_id" validate:"required"`
	StartingCapital float64 `json:"starting_capital" validate:"required,gte=10000,lte=100000"`
}

type CustomerOutputDto struct {
	CurrentCapital    float64           `json:"starting_capital" validate:"required,gte=10000,lte=100000"`
	RegisteredLoaders []LoaderOutputDto `json:"registered_loaders"`
}

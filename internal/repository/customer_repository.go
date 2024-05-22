package repository

import (
	"context"
	"database/sql"
	"loaders-online/internal/entity/dto"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (c CustomerRepository) CreateCustomer(ctx context.Context, customer *dto.CustomerSignUpDto) error {
	_, err := c.db.ExecContext(ctx, "INSERT INTO customers (customer_id, starting_capital, current_capital) VALUES ($1, $2, $3)",
		customer.CustomerID, customer.StartingCapital, customer.StartingCapital)
	return err
}

func (c CustomerRepository) GetCustomerById(ctx context.Context, id int) (dto.CustomerOutputDto, error) {
	var customer dto.CustomerOutputDto
	err := c.db.QueryRowContext(ctx, "SELECT current_capital FROM customers WHERE customer_id = $1", id).
		Scan(&customer.CurrentCapital)
	if err != nil {
		return customer, err
	}

	var loaders []dto.LoaderOutputDto
	rows, err := c.db.Query("select l.loader_id, l.max_weight, l.drunkenness, l.fatigue, l.salary from loaders as l " +
		"inner join loaders_tasks as lt on l.loader_id = lt.loader_id inner join tasks as t on lt.task_id = t.task_id " +
		"where customer_id = 0;\n ")
	if err != nil {
		return customer, err
	}

	for rows.Next() {
		var loader dto.LoaderOutputDto
		rows.Scan(&loader.LoaderID, &loader.MaxWeight, &loader.Drunkenness, &loader.Fatigue, &loader.Salary)
		loaders = append(loaders, loader)
	}

	customer.RegisteredLoaders = loaders
	return customer, err
}

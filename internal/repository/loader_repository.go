package repository

import (
	"context"
	"database/sql"
	"fmt"
	"loaders-online/internal/entity/dto"
)

type LoaderRepository struct {
	db *sql.DB
}

func NewLoaderRepository(db *sql.DB) *LoaderRepository {
	return &LoaderRepository{db: db}
}

func (l *LoaderRepository) GetLoaderById(ctx context.Context, id int) (*dto.LoaderOutputDto, error) {
	var loader dto.LoaderOutputDto
	row := l.db.QueryRowContext(ctx, "SELECT max_weight, drunkenness, fatigue, salary FROM loaders WHERE loader_id = $1", id).
		Scan(&loader.MaxWeight, &loader.Drunkenness, &loader.Fatigue, &loader.Salary)
	loader.LoaderID = id
	return &loader, row
}

func (l *LoaderRepository) GetAssignedLoaders(ctx context.Context, id int) ([]dto.LoaderOutputDto, error) {
	var loaders []dto.LoaderOutputDto
	rows, err := l.db.QueryContext(ctx, "select l.loader_id, l.max_weight, l.drunkenness, l.fatigue, l.salary from loaders as l inner join loaders_tasks as lt on l.loader_id = lt.loader_id inner join tasks as t on lt.task_id = t.task_id inner join customers as c on t.customer_id = c.customer_id where c.customer_id = $1", id)
	if err != nil {
		return loaders, err
	}
	fmt.Println(rows)
	for rows.Next() {
		fmt.Println("ASDASD")
		var loader dto.LoaderOutputDto
		err := rows.Scan(&loader.LoaderID, &loader.MaxWeight, &loader.Drunkenness, &loader.Fatigue, &loader.Salary)
		if err != nil {
			return loaders, err
		}
		loaders = append(loaders, loader)
	}
	return loaders, err
}

func (l *LoaderRepository) CreateLoader(ctx context.Context, loader *dto.LoaderOutputDto) error {
	_, err := l.db.ExecContext(ctx, "INSERT INTO loaders (loader_id, max_weight, drunkenness, fatigue, salary) VALUES ($1, $2, $3, $4, $5)",
		loader.LoaderID, loader.MaxWeight, loader.Drunkenness, loader.Fatigue, loader.Salary)
	return err
}

func (l *LoaderRepository) UpdateLoader(ctx context.Context, loader *dto.LoaderOutputDto) error {
	_, err := l.db.ExecContext(ctx, "UPDATE loaders SET max_weight = $1, drunkenness = $2, fatigue = $3 WHERE loader_id = $4",
		loader.MaxWeight, loader.Drunkenness, loader.Fatigue, loader.LoaderID)
	return err
}

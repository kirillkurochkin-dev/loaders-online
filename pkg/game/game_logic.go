package game

import (
	"errors"
	"loaders-online/internal/entity/dto"
)

func Recalculate(loader *dto.LoaderOutputDto) *dto.LoaderOutputDto {
	if loader.Drunkenness {
		loader.MaxWeight = loader.MaxWeight * (100 - (loader.Fatigue + 50)) / 100
		loader.Drunkenness = false
	} else {
		loader.MaxWeight = loader.MaxWeight * (100 - (loader.Fatigue)) / 100
	}

	if loader.MaxWeight < 5 {
		loader.MaxWeight = 5
	}
	if loader.MaxWeight > 30 {
		loader.MaxWeight = 30
	}

	return loader
}
func DoJob(loader dto.LoaderOutputDto) dto.LoaderOutputDto {
	loader.Fatigue += 20
	if loader.Fatigue > 100 {
		loader.Fatigue = 100
	}
	return loader
}
func CalcMoney(loaders []dto.LoaderOutputDto, customer *dto.CustomerOutputDto) error {
	var sum float64
	for _, loader := range loaders {
		sum += loader.Salary
	}
	if customer.CurrentCapital < sum {
		return errors.New("not enough capital")
	} else {
		customer.CurrentCapital -= sum
	}
	return nil
}

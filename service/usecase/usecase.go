package usecase

import (
	"bd_tp/models"
	serviceRepo "bd_tp/service/repository"
)

type Usecase struct {
	repository *serviceRepo.Repository
}

func NewServiceUseCase(sR *serviceRepo.Repository) *Usecase {
	return &Usecase{
		repository: sR,
	}
}

func (sU *Usecase) Clear() error {
	return sU.repository.Clear()
}

func (sU *Usecase) GetStatus() (*models.Status, error) {
	return sU.repository.GetStatus()
}
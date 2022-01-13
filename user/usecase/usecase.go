package usecase

import (
	userRepo "bd_tp/user/repository"
	"bd_tp/models"
	"fmt"
)

type Usecase struct {
	userRepo *userRepo.Repository
}

func NewUserUsecase (uR *userRepo.Repository) *Usecase {
	return &Usecase{
		userRepo: uR,
	}
}

func (uR *Usecase) CreateUser (u *models.User) (*models.User, bool, error) {
	resultUser,isNew,err := uR.userRepo.CreateUser(u)
	if err != nil {
		fmt.Println(err)
		return nil,isNew, err
	}
	return resultUser, isNew, nil
}

func (uR *Usecase) ProfileInfo (nickname string) (*models.User,error) {
	resultUser, err := uR.userRepo.ProfileInfo(nickname)
	if err != nil {
		return nil, err
	}
	return resultUser, nil
}

func (uR *Usecase) UpdateProfile (u *models.User) (*models.User, bool, error) {
	resultUser, isFound, err := uR.userRepo.UpdateProfile(u)
	if err != nil {
		return nil, isFound, err
	}
	return resultUser, isFound, nil
}
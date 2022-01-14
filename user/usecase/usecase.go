package usecase

import (
	"bd_tp/models"
	userRepo "bd_tp/user/repository"
)

type Usecase struct {
	userRepo *userRepo.Repository
}

func NewUserUsecase (uR *userRepo.Repository) *Usecase {
	return &Usecase{
		userRepo: uR,
	}
}

func (uR *Usecase) CreateUser (u *models.User) ([]models.User, bool, error) {
	users,isNew,err := uR.userRepo.CreateUser(u)
	if err != nil {
		
		return nil,isNew, err
	}
	return users, isNew, nil
}

func (uR *Usecase) ProfileInfo (nickname string) (*models.User,error) {
	resultUser, err := uR.userRepo.ProfileInfo(nickname)
	if err != nil {
		return nil, err
	}
	return resultUser, nil
}

func (uR *Usecase) UpdateProfile (u *models.User) (*models.User, bool, error) {
	oldProfile,err := uR.userRepo.ProfileInfo(u.Nickname) 
	if err != nil {
		return nil,false,err
	}
	if u.Fullname == "" {
		u.Fullname = oldProfile.Fullname
	}
	if u.Email == "" {
		u.Email = oldProfile.Email
	}
	if u.About == "" {
		u.About = oldProfile.About
	}
	resultUser, isFound, err := uR.userRepo.UpdateProfile(u)
	if err != nil {
		return nil, isFound, err
	}
	return resultUser, isFound, nil
}
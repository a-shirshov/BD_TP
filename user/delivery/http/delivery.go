package delivery

import (
	"bd_tp/response"
	userUseCase "bd_tp/user/usecase"
	"net/http"
	"strings"
)

type UserDelivery struct {
	userUsecase *userUseCase.Usecase
}

func NewUserDelivery(uU *userUseCase.Usecase) *UserDelivery {
	return &UserDelivery{
		userUsecase: uU,
	}
}

func (uD *UserDelivery) CreateUser(w http.ResponseWriter, r *http.Request) {
	u, err := response.GetUserFromRequest(r.Body)

	path := r.URL.Path
	split := strings.Split(path, "/")
	nickname := split[len(split)-2]
	u.Nickname = nickname

	if err != nil {

		return
	}
	users, isNew, err := uD.userUsecase.CreateUser(u)
	if err != nil {

		return
	}
	var statusCode int
	if isNew {
		statusCode = 201
		userResponse := &response.UserResponse{
			Nickname: users[0].Nickname,
			Fullname: users[0].Fullname,
			About:    users[0].About,
			Email:    users[0].Email,
		}
		response.SendResponse(w, statusCode, userResponse)
		return
	} else {
		statusCode = 409
		var usersResponse []response.UserResponse
		for _, user := range users {
			userResponse := &response.UserResponse{
				Nickname: user.Nickname,
				Fullname: user.Fullname,
				About:    user.About,
				Email:    user.Email,
			}
			usersResponse = append(usersResponse, *userResponse)
		}
		response.SendResponse(w, statusCode, usersResponse)
		return
	}
}

func (uD *UserDelivery) ProfileInfo(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	split := strings.Split(path, "/")
	nickname := split[len(split)-2]

	user, err := uD.userUsecase.ProfileInfo(nickname)
	if err != nil {
		statusCode := 404
		errorResponse := &response.Error{
			Message: "Can't find user with nickname:" + nickname,
		}
		response.SendResponse(w, statusCode, errorResponse)
		return
	}
	statusCode := 200
	userResponse := &response.UserResponse{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		About:    user.About,
		Email:    user.Email,
	}
	response.SendResponse(w, statusCode, userResponse)
}

func (uD *UserDelivery) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	u, err := response.GetUserFromRequest(r.Body)
	if err != nil {
		return
	}
	path := r.URL.Path
	split := strings.Split(path, "/")
	nickname := split[len(split)-2]
	u.Nickname = nickname

	user, isFound, err := uD.userUsecase.UpdateProfile(u)

	if err != nil {
		var statusCode int
		var errorResponse *response.Error
		if !isFound {
			errorResponse = &response.Error{
				Message: "Can't find user with nickname:" + nickname,
			}
			statusCode = 404
		} else {
			errorResponse = &response.Error{
				Message: "Conficts with other users",
			}
			statusCode = 409
		}
		response.SendResponse(w, statusCode, errorResponse)
		return
	}
	statusCode := 200
	userResponse := &response.UserResponse{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		About:    user.About,
		Email:    user.Email,
	}
	response.SendResponse(w, statusCode, userResponse)
}

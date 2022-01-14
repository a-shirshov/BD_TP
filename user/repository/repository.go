package repository

import (
	models "bd_tp/models"
	"strings"
	"errors"
	sql "github.com/jmoiron/sqlx"
)

const (
	createUserQuery = `insert into "user" (nickname, fullname, about, email) values ($1, $2, $3, $4)`
	findUserByNicknameQuery = `select * from "user" where nickname = $1`
	findUserByEmailQuery = `select * from "user" where email = $1`
	updateUserQuery = `update "user" set fullname = $1,about = $2,email = $3 where nickname = $4 `
	getUserID = `select id from "user" where nickname = $1`
)

type Repository struct {
	db *sql.DB
}

func NewRepository (db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (uR *Repository) CreateUser(u *models.User) ([]models.User, bool, error) {
	isNew := true
	query := createUserQuery
	var users []models.User
	rows, err := uR.db.Query(query,u.Nickname,u.Fullname,u.About,u.Email)
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint`) {
			query := findUserByNicknameQuery
			userByNickname := models.User{}
			err := uR.db.Get(&userByNickname,query,u.Nickname)
			if err == nil {
				users = append(users, userByNickname)
			}
			query = findUserByEmailQuery
			userByEmail := models.User{}
			err = uR.db.Get(&userByEmail,query,u.Email)
			if err == nil {
				if userByNickname.Email != userByEmail.Email{
					users = append(users, userByEmail)
				}
			}
			isNew = false
			return users, isNew, nil
		}

	}
	defer rows.Close()
	users = append(users, *u)
	return users, isNew, nil
}

func (uR* Repository) ProfileInfo (nickname string) (*models.User, error) {
	query := findUserByNicknameQuery
	user := models.User{}
	err := uR.db.Get(&user,query,nickname)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (uR *Repository) UpdateProfile (u *models.User) (*models.User, bool, error) {
	query := updateUserQuery
	user := models.User{}
	isFound := true
	result,err := uR.db.Exec(query,u.Fullname,u.About,u.Email,u.Nickname)
	if err != nil {
		
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "user_email_key"`) {
			return nil,isFound,err
		}
		isFound := false
		return nil, isFound, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		isFound := false
		return nil, isFound, errors.New("no user with with nickname")
	}
	
	query = findUserByNicknameQuery
	err = uR.db.Get(&user,query,u.Nickname)
	if err != nil {
		
		return nil, isFound, err
	}
	return &user, isFound, nil
}

func (uR *Repository) GetIdByNickname (nickname string) (int, error) {
	query := getUserID
	var userID int
	err := uR.db.Get(&userID,query,nickname)
	if err != nil {
		return 0, err
	}
	return userID,nil
}
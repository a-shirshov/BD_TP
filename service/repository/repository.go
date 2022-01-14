package repository

import (
	"bd_tp/models"
	sql "github.com/jmoiron/sqlx"
)

const (
	clearQuery = `truncate forum_user, vote, post, thread, forum, "user"`
	GetStatusUserQuery   = `select count(*) from "user"`
	GetStatusForumQuery  = `select count(*) from forum`
	GetStatusThreadQuery = `select count(*) from thread`
	GetStatusPostQuery   = `select count(*) from post`
)

type Repository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (sR *Repository) Clear() error {
	query := clearQuery
	_, err := sR.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (sR *Repository) GetStatus() (*models.Status, error) {
	var query string
	status := &models.Status{}
	query = GetStatusUserQuery
	err := sR.db.QueryRow(query).Scan(&status.User)
	if err != nil {
		return nil, err
	}
	query = GetStatusForumQuery
	err = sR.db.QueryRow(query).Scan(&status.Forum)
	if err != nil {
		return nil, err
	}
	query = GetStatusThreadQuery
	err = sR.db.QueryRow(query).Scan(&status.Thread)
	if err != nil {
		return nil, err
	}
	query = GetStatusPostQuery
	err = sR.db.QueryRow(query).Scan(&status.Post)
	if err != nil {
		return nil, err
	}
	return status, nil
}
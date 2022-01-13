package repository

import (
	sql "github.com/jmoiron/sqlx"
	"bd_tp/models"
	//"strings"
	//"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewPostRepository (db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

const (
	FindPostByID = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,p.thread_id as thread,p.created from "post" as p
    join "thread" as t on t.id = p.thread_id
    join "forum" as f on f.id = t.forum_id
    join "user" as u on u.id = p.user_id
    where p.id = $1;`
)

func (pR *Repository) GetPostByID(id int) (*models.Post,error) {
	query := FindPostByID
	var post models.Post
	err := pR.db.Get(&post,query,id)
	if err != nil {
		return nil,err
	}
	return &post,err
}
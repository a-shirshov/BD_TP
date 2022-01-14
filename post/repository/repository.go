package repository

import (
	"bd_tp/models"

	sql "github.com/jmoiron/sqlx"
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

	UpdatePostByIdQuery = `update post set message = $1,edited = true 
	where id = $2 
	returning id;` 
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

func (pR *Repository) UpdatePost(post *models.Post) (*models.Post, error) {
	query := UpdatePostByIdQuery
	var updatedPost models.Post
	err := pR.db.QueryRow(query,post.Message,post.ID).Scan(&updatedPost.ID)
	if err != nil {
		
		return nil,err
	}
	query = FindPostByID
	err = pR.db.Get(&updatedPost,query,&updatedPost.ID)
	if err != nil {
		
		return nil,err
	}
	return &updatedPost,nil
}
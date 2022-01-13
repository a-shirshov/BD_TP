package repository

import (
	"bd_tp/models"
	"fmt"

	sql "github.com/jmoiron/sqlx"
)

const (
	selectPostsWithIdQuery = `select p.id, p.edited, f.slug as forum, p.created from "post" as p
    join "thread" as t on p.thread_id = t.id 
    join "forum" as f on f.id = t.forum_id
    where t.id = $1 and p.id = $2;`

	selectPostsWithSlugQuery = `select p.id, p.edited, f.slug as forum, p.created,t.id as thread from "post" as p
    join "thread" as t on p.thread_id = t.id 
    join "forum" as f on f.id = t.forum_id
    where t.slug = $1 and p.id = $2;`

	createPostsWithIdQuery = `insert into "post" (parent,message,user_id,thread_id,created) 
    select $1,$2,u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = $3 AND t.id = $4 returning id;`

	createPostsWithSlugQuery = `insert into "post" (parent,message,user_id,thread_id,created) 
    select $1,$2,u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = $3 AND t.slug = $4 returning id;`

	findThreadWithIdQuery = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,t.slug,t.created from "thread" as t
    join "user" as u on u.id = t.user_id
    join "forum" as f on f.id = t.forum_id
    where t.id = $1;`

	findThreadWithSlugQuery = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,t.slug,t.created from "thread" as t
    join "user" as u on u.id = t.user_id
    join "forum" as f on f.id = t.forum_id
    where t.slug = $1;`

	updateThreadWithIdQuery = `update "thread" set title = $1,message = $2
    where id = $3;`

	updateThreadWithSlugQuery = `update "thread" set title = $1,message = $2
    where slug = $3;`

	insertVoteWithIdQuery = `insert into "vote" (voice,user_id,thread_id)
    select $1, u.id, t.id from "user" as u, "thread" as t
    where u.nickname = $2 and t.id = $3;`

	insertVoteWithSlugQuery = `insert into "vote" (voice,user_id,thread_id)
    select $1, u.id, t.id from "user" as u, "thread" as t
    where u.nickname = $2 and t.slug = $3;`
)

type Repository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (fR *Repository) CreatePostsWithID (posts []models.Post, id int) ([]models.Post,int,error) {
	query := createPostsWithIdQuery
	var postID int
	for index,post := range posts {
		err := fR.db.Get(&postID,query,post.Parent,post.Message,post.Author,id)
		if err != nil {
			fmt.Println(err)
			return nil,500,err
		}
		posts[index].ID = postID
		fmt.Println(posts[index].ID)
		fmt.Println(posts[index])
	}

	query = selectPostsWithIdQuery 
	for index,post := range posts {
		err := fR.db.Get(&(posts[index]),query,id,post.ID)
		if err != nil {
			fmt.Println("Here")
			fmt.Println(err)
			return nil,500,err
		}
		posts[index].Thread = id
	}
	return posts, 201, nil
}

func (fR *Repository) CreatePostsWithSlug (posts []models.Post, slug string) ([]models.Post,int,error) {
	query := createPostsWithSlugQuery
	var postID int
	for index,post := range posts {
		err := fR.db.Get(&postID,query,post.Parent,post.Message,post.Author,slug)
		if err != nil {
			fmt.Println(err)
			return nil,500,err
		}
		posts[index].ID = postID
		fmt.Println(posts[index].ID)
		fmt.Println(posts[index])
	}

	query = selectPostsWithSlugQuery 
	for index,post := range posts {
		err := fR.db.Get(&(posts[index]),query,slug,post.ID)
		if err != nil {
			fmt.Println("Here")
			fmt.Println(err)
			return nil,500,err
		}
	}
	return posts, 201, nil
}

func (fR *Repository) ThreadDetailsByID (id int) (*models.Thread, error) {
	query := findThreadWithIdQuery
	var thread models.Thread
	err := fR.db.Get(&thread,query,id)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	//thread.ID = id
	return &thread,nil
}

func (fR *Repository) ThreadDetailsBySlug (slug string) (*models.Thread,error) {
	query := findThreadWithSlugQuery
	var thread models.Thread
	err := fR.db.Get(&thread,query,slug)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	//thread.Slug = slug
	return &thread,nil
}

func (fR *Repository) ThreadDetailsUpdateByID (threadInfo *models.Thread, id int) (*models.Thread, error) {
	query := updateThreadWithIdQuery
	var thread models.Thread
	_,err := fR.db.Query(query,threadInfo.Title,threadInfo.Message,id)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	query = findThreadWithIdQuery
	err = fR.db.Get(&thread,query,id)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	thread.ID = id
	return &thread,nil
}

func (fR *Repository) ThreadDetailsUpdateBySlug (threadInfo *models.Thread, slug string) (*models.Thread, error) {
	query := updateThreadWithSlugQuery
	var thread models.Thread
	_,err := fR.db.Query(query,threadInfo.Title,threadInfo.Message,slug)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	query = findThreadWithSlugQuery
	err = fR.db.Get(&thread,query,slug)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	//thread.ID = id
	return &thread,nil
}

func (fR *Repository) ThreadVoteByID (vote *models.Vote, id int) (*models.Thread, error) {
	query := insertVoteWithIdQuery
	var thread models.Thread
	_,err := fR.db.Query(query,vote.Voice,vote.Nickname,id)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	query = findThreadWithIdQuery
	err = fR.db.Get(&thread,query,id)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	thread.ID = id
	return &thread,nil
}

func (fR *Repository) ThreadVoteBySlug (vote *models.Vote, slug string) (*models.Thread, error) {
	query := insertVoteWithSlugQuery
	var thread models.Thread
	_,err := fR.db.Query(query,vote.Voice,vote.Nickname,slug)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	query = findThreadWithSlugQuery
	err = fR.db.Get(&thread,query,slug)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	//thread.ID = id
	return &thread,nil
}


package repository

import (
	"bd_tp/models"
	"errors"
	"fmt"
	"strings"
	"time"

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

	createPostsWithIdQuery = `insert into "post" (parent,message,user_id,thread_id,created,forum) 
    select $1,$2,u.id,t.id,$5,$6 from "user" as u, "thread" as t 
    where u.nickname = $3 AND t.id = $4  returning id;`

	createPostsWithSlugQuery = `insert into "post" (parent,message,user_id,thread_id,created,forum) 
    select $1,$2,u.id,t.id,$5,$6 from "user" as u, "thread" as t 
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

	updateVoteWithIdQuery = `update "vote" set voice = $1 
	where thread_id = $2 
	and user_id = $3`

	getPostsByThreadFlatDescSince = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 and p.id < $2
	order by p.id desc limit $3`

	getPostsByThreadFlatDesc = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 
	order by p.id desc limit $2`

	getPostsByThreadFlatAscSince = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 and p.id > $2
	order by p.id asc limit $3`

	getPostsByThreadFlatAsc = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 
	order by p.id asc limit $2`
)

type Repository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (fR *Repository) CreatePostsWithID(posts []models.Post, id int) ([]models.Post, int, error) {
	query := createPostsWithIdQuery
	var postID int
	created := time.Now()
	for index, post := range posts {
		fmt.Println("post = ",post)
		err := fR.db.Get(&postID, query, post.Parent, post.Message, post.Author, id, created, post.Forum)
		if err != nil {
			fmt.Println("In create err",err)
			return nil, 500, err
		}
		posts[index].ID = postID
	}

	query = selectPostsWithIdQuery
	for index, post := range posts {
		err := fR.db.Get(&(posts[index]), query, id, post.ID)
		if err != nil {
			fmt.Println("In select err",err)
			return nil, 500, err
		}
		posts[index].Thread = id
	}
	return posts, 201, nil
}

func (fR *Repository) CreatePostsWithSlug(posts []models.Post, slug string) ([]models.Post, int, error) {
	query := createPostsWithSlugQuery
	var postID int
	created := time.Now()
	for index, post := range posts {
		err := fR.db.Get(&postID, query, post.Parent, post.Message, post.Author, slug, created,post.Forum)
		if err != nil {

			return nil, 500, err
		}
		posts[index].ID = postID
	}

	query = selectPostsWithSlugQuery
	for index, post := range posts {
		err := fR.db.Get(&(posts[index]), query, slug, post.ID)
		if err != nil {

			return nil, 500, err
		}
	}
	return posts, 201, nil
}

func (fR *Repository) ThreadDetailsByID(id int) (*models.Thread, error) {
	query := findThreadWithIdQuery
	var thread models.Thread
	err := fR.db.Get(&thread, query, id)
	if err != nil {

		return nil, err
	}
	//thread.ID = id
	return &thread, nil
}

func (fR *Repository) ThreadDetailsBySlug(slug string) (*models.Thread, error) {
	query := findThreadWithSlugQuery
	var thread models.Thread
	err := fR.db.Get(&thread, query, slug)
	if err != nil {

		return nil, err
	}
	//thread.Slug = slug
	return &thread, nil
}

func (fR *Repository) ThreadDetailsUpdateByID(threadInfo *models.Thread, id int) (*models.Thread, error) {
	query := updateThreadWithIdQuery
	var thread models.Thread
	_, err := fR.db.Query(query, threadInfo.Title, threadInfo.Message, id)
	if err != nil {

		return nil, err
	}
	query = findThreadWithIdQuery
	err = fR.db.Get(&thread, query, id)
	if err != nil {

		return nil, err
	}
	thread.ID = id
	return &thread, nil
}

func (fR *Repository) ThreadDetailsUpdateBySlug(threadInfo *models.Thread, slug string) (*models.Thread, error) {
	query := updateThreadWithSlugQuery
	var thread models.Thread
	_, err := fR.db.Query(query, threadInfo.Title, threadInfo.Message, slug)
	if err != nil {

		return nil, err
	}
	query = findThreadWithSlugQuery
	err = fR.db.Get(&thread, query, slug)
	if err != nil {

		return nil, err
	}
	//thread.ID = id
	return &thread, nil
}

func (fR *Repository) ThreadVoteByID(vote *models.Vote, id int,userId int) (*models.Thread, error) {
	query := insertVoteWithIdQuery
	var thread models.Thread
	_, err := fR.db.Query(query, vote.Voice, vote.Nickname, id)
	if err != nil {
		if strings.Contains(err.Error(),"duplicate") {
			query := updateVoteWithIdQuery 
			_, err = fR.db.Query(query,vote.Voice,id,userId)
			if err != nil {
				fmt.Println(err)
				return nil,err
			}
		}
	}
	query = findThreadWithIdQuery
	err = fR.db.Get(&thread, query, id)
	if err != nil {

		return nil, err
	}
	thread.ID = id
	return &thread, nil
}

func (fR *Repository) ThreadVoteBySlug(vote *models.Vote, slug string,userId int) (*models.Thread, error) {
	query := insertVoteWithSlugQuery
	var thread models.Thread
	_, err := fR.db.Query(query, vote.Voice, vote.Nickname, slug)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	query = findThreadWithSlugQuery
	err = fR.db.Get(&thread, query, slug)
	if err != nil {

		return nil, err
	}
	//thread.ID = id
	return &thread, nil
}

func (fR *Repository) ThreadGetPosts(id int, limit, since, sort, desc string) ([]models.Post,error) {
	if sort == "" || sort == "flat" {
		return fR.threadGetPostsFlat(id, limit, since,sort, desc)
	}
	return nil,errors.New("недописал")
}

func (fR *Repository) threadGetPostsFlat(id int, limit, since, sort, desc string) ([]models.Post, error) {
	var query string
	var rows *sql.Rows
	var err error
	if desc == "true" {
		if since != "" {
			query = getPostsByThreadFlatDescSince
			rows, err = fR.db.Queryx(query,id,since,limit)
		} else {
			query = getPostsByThreadFlatDesc
			rows, err = fR.db.Queryx(query,id,limit)
		}
	} else {
		if since != "" {
			query = getPostsByThreadFlatAscSince
			rows, err = fR.db.Queryx(query,id,since,limit)
		} else {
			query = getPostsByThreadFlatAsc
			rows, err = fR.db.Queryx(query,id,limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		post := models.Post{}
		rows.Scan(&post.ID,&post.Parent,&post.Author,&post.Message,&post.Edited,&post.Forum,&post.Thread,&post.Created)
		posts = append(posts, post)
	}
	if err != nil {

		return nil, err
	}
	return posts, nil
}

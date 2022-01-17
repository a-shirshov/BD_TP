package repository

import (
	"bd_tp/models"
	"errors"
	"strconv"
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

	selectPostsNew = `select id,parent,message,edited,forum,thread_id,created from "post`

	createPostsWithIdQuery = `insert into "post" (parent,user_id,message,thread_id,created,forum_id) 
    select $1,u.id,$2,t.id,$3,$4 from "user" as u, "thread" as t 
    where u.nickname = $5 AND t.id = $6  returning id;`

	createPostsWithSlugQuery = `insert into "post" (parent,user_id,message,thread_id,created,forum_id) 
    select $1,u.id,$2,t.id,$3,$4 from "user" as u, "thread" as t 
    where u.nickname = $5 AND t.slug = $6  returning id;`

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

	getPostsByThreadTreeDescSince = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 and 
	p.path < (select path from post where id =$2)
	order by p.path desc, id desc limit $3`

	getPostsByThreadTreeDesc = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 
	order by p.path desc, id desc limit $2`

	getPostsByThreadTreeAscSince = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 and 
	p.path > (select path FROM post where id = $2) 
	order by path, id limit $3`

	getPostsByThreadTreeAsc = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1
	order by p.path, id limit $2`

	getPostsByThreadParentTreeDescSince = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created 
	from post as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.path[1] in 
	(select id from post where thread_id = $1 and parent = 0 and path[1] <
	(select path[1] from post where id = $2) order by id desc limit $3)
	order by p.path[1] desc, p.path, id`

	getPostsByThreadParentTreeDesc = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created 
	from post as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.path[1] in 
	(select id from post where thread_id = $1 and parent = 0 
	order by id desc limit $2)
	order by p.path[1] desc, p.path, id`

	getPostsByThreadParentTreeAsc = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created 
	from post as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.path[1] in 
	(select id from post where thread_id = $1 and parent = 0 order by id limit $2)
	order by p.path, id`

	getPostsByThreadParentTreeAscSince = `select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created 
	from post as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.path[1] in 
	(select id from post where thread_id = $1 and parent = 0 and path[1] > 
	(select path[1] from post where id = $2) 
	order by id limit $3)
	order by p.path, id`

	createPostsNewQuery = `insert into post (parent, user_id, message, forum, thread_id, created) values ($1,$2,$3,$4,$5,$6) returning id`
)

type Repository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (fR *Repository) CreatePostsWithID(posts []models.Post, threadId int, forumId int) ([]models.Post, int, error) {
	query := createPostsWithIdQuery
	var postID int
	created := time.Now()
	for index, post := range posts {
		err := fR.db.Get(&postID, query, post.Parent, post.Message, created, forumId, post.Author, threadId)
		if err != nil {
			return nil, 404, err
		}
		posts[index].ID = postID
	}

	query = selectPostsWithIdQuery
	for index, post := range posts {
		err := fR.db.Get(&(posts[index]), query, threadId, post.ID)
		if err != nil {
			return nil, 404, err
		}
		posts[index].Thread = threadId
	}
	return posts, 201, nil
}

func (fR *Repository) CreatePostsWithSlug(posts []models.Post, threadSlug string, forumId int) ([]models.Post, int, error) {
	query := createPostsWithSlugQuery
	var postID int
	created := time.Now()
	for index, post := range posts {
		err := fR.db.Get(&postID, query, post.Parent, post.Message, created, forumId, post.Author, threadSlug)
		if err != nil {

			return nil, 404, err
		}
		posts[index].ID = postID
	}

	query = selectPostsWithSlugQuery
	for index, post := range posts {
		err := fR.db.Get(&(posts[index]), query, threadSlug, post.ID)
		if err != nil {

			return nil, 404, err
		}
	}
	return posts, 201, nil
}

func (fR *Repository) CreatePostsNew(threadId int, forum string, posts []models.Post, users []int) ([]models.Post, error) {
	created := time.Now()
	for index := range posts {
		if posts[index].Parent != 0 {
			id := -1
			err := fR.db.QueryRow(`select id from post where thread = $1 and id = $2`, threadId, posts[index].Parent).Scan(&id)
			if err != nil {
				return nil,err
			}
		}
		query := createPostsNewQuery
		err := fR.db.Get(&(posts[index].ID), query, posts[index].Parent, users[index], posts[index].Message, forum, threadId, created)
		if err != nil {
			return nil, err
		}
	}

	var newPosts []models.Post
	query := selectPostsNew
	for index := range posts {
		var newPost models.Post
		err := fR.db.Get(&newPost, query, posts[index].Parent, users[index], posts[index].Message, forum, threadId, created)
		if err != nil {
			return nil, err
		}
		newPost.Author = posts[index].Author
		newPosts = append(newPosts, newPost)
	}
	return newPosts, nil
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

func (fR *Repository) ThreadDetails(slug_or_id string) (*models.Thread, error) {
	var query string
	_, err := strconv.Atoi(slug_or_id)
	if err != nil {
		query = findThreadWithSlugQuery
	} else {
		query = findThreadWithIdQuery
	}
	var thread models.Thread
	err = fR.db.Get(&thread, query, slug_or_id)
	if err != nil {

		return nil, err
	}
	//thread.ID = id
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

func (fR *Repository) ThreadVoteByID(vote *models.Vote, id int, userId int) (*models.Thread, error) {
	query := insertVoteWithIdQuery
	var thread models.Thread
	_, err := fR.db.Exec(query, vote.Voice, vote.Nickname, id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			query := updateVoteWithIdQuery
			_, err = fR.db.Exec(query, vote.Voice, id, userId)
			if err != nil {

				return nil, err
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

func (fR *Repository) ThreadVoteBySlug(vote *models.Vote, slug string, userId int) (*models.Thread, error) {
	query := insertVoteWithSlugQuery
	var thread models.Thread
	_, err := fR.db.Exec(query, vote.Voice, vote.Nickname, slug)
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

func (fR *Repository) ThreadGetPosts(id int, limit, since, sort, desc string) ([]models.Post, error) {
	if sort == "" || sort == "flat" {
		return fR.threadGetPostsFlat(id, limit, since, sort, desc)
	} else if sort == "tree" {
		return fR.threadGetPostsTree(id, limit, since, desc)
	} else if sort == "parent_tree" {
		return fR.threadGetPostsParentTree(id, limit, since, desc)
	}
	return nil, errors.New("недописал")
}

func (fR *Repository) threadGetPostsFlat(id int, limit, since, sort, desc string) ([]models.Post, error) {
	var query string
	var rows *sql.Rows
	var err error
	if desc == "true" {
		if since != "" {
			query = getPostsByThreadFlatDescSince
			rows, err = fR.db.Queryx(query, id, since, limit)
		} else {
			query = getPostsByThreadFlatDesc
			rows, err = fR.db.Queryx(query, id, limit)
		}
	} else {
		if since != "" {
			query = getPostsByThreadFlatAscSince
			rows, err = fR.db.Queryx(query, id, since, limit)
		} else {
			query = getPostsByThreadFlatAsc
			rows, err = fR.db.Queryx(query, id, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		post := models.Post{}
		rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.Edited, &post.Forum, &post.Thread, &post.Created)
		posts = append(posts, post)
	}
	if err != nil {

		return nil, err
	}
	return posts, nil
}

func (tR *Repository) threadGetPostsTree(id int, limit string, since string, desc string) ([]models.Post, error) {
	var err error
	var rows *sql.Rows
	if desc == "true" {
		if since != "" {
			query := getPostsByThreadTreeDescSince
			rows, err = tR.db.Queryx(query, id, since, limit)
		} else {
			query := getPostsByThreadTreeDesc
			rows, err = tR.db.Queryx(query, id, limit)
		}
	} else {
		if since != "" {
			query := getPostsByThreadTreeAscSince
			rows, err = tR.db.Queryx(query, id, since, limit)
		} else {
			query := getPostsByThreadTreeAsc
			rows, err = tR.db.Queryx(query, id, limit)
		}
	}

	if err != nil {

		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		post := models.Post{}
		rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.Edited, &post.Forum, &post.Thread, &post.Created)
		posts = append(posts, post)
	}
	if err != nil {

		return nil, err
	}
	return posts, nil
}

func (tR *Repository) threadGetPostsParentTree(id int, limit string, since string, desc string) ([]models.Post, error) {

	var rows *sql.Rows
	var err error

	if desc == "true" {
		if since != "" {
			query := getPostsByThreadParentTreeDescSince
			rows, err = tR.db.Queryx(query, id, since, limit)
		} else {
			query := getPostsByThreadParentTreeDesc
			rows, err = tR.db.Queryx(query, id, limit)
		}
	} else {
		if since != "" {
			query := getPostsByThreadParentTreeAscSince
			rows, err = tR.db.Queryx(query, id, since, limit)
		} else {
			query := getPostsByThreadParentTreeAsc
			rows, err = tR.db.Queryx(query, id, limit)
		}
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		post := models.Post{}
		rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.Edited, &post.Forum, &post.Thread, &post.Created)
		posts = append(posts, post)
	}
	if err != nil {

		return nil, err
	}
	return posts, nil
}

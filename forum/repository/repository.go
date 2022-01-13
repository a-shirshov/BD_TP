package repository

import (
	sql "github.com/jmoiron/sqlx"
	"bd_tp/models"
	"strings"
	"fmt"
)

const (
	createForumQuery = `insert into "forum" (title,slug,user_id) values ($1,$2,$3)`
	findForumQuery = `select title from "forum" where slug = $1`
	findForumJoinQuery = `select f.title, u.nickname from "forum" as f join "user" as u on f.user_id = u.id where f.slug= $1`
	createForumBranchQuery = `insert into "thread" (title,message,created,user_id,forum_id) values ($1,$2,$3,$4,$5)`
	GetForumIdAndTitle = `select id,title from "forum" where slug = $1`
	getBranchInfo = `select id,title,message,slug,created from "thread" where title = $1`
	getThreadsByForumDesc = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,f.slug,t.created from "thread" as t 
	join "forum" as f on t.forum_id = f.id  
	join "user" as u on u.id = t.user_id where f.slug = $1 
	ORDER BY t.title DESC
	LIMIT $2;`

	getThreadsByForumAsc = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,f.slug,t.created from "thread" as t 
	join "forum" as f on t.forum_id = f.id  
	join "user" as u on u.id = t.user_id where f.slug = $1 
	ORDER BY t.title ASC
	LIMIT $2;`
)

type Repository struct {
	db *sql.DB
}

func NewForumRepository (db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (fR *Repository) CreateForum (f *models.Forum, userID int) (*models.Forum, int, error) {
	fmt.Println("here")
	query := createForumQuery
	_,err := fR.db.Query(query,f.Title,f.Slug,userID)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "forum_slug_key"`) {
			fmt.Println(err.Error())
			query := findForumQuery
			forum := models.Forum{}
			err := fR.db.Get(&forum,query,f.Slug)
			if err != nil {
				return nil, 500, err
			}
			fmt.Println(err)
			forum.Slug = f.Slug
			forum.User = f.User
			fmt.Println(forum)
			return &forum, 409, nil
		}
		return nil, 500, err
	}
	f.Threads = 0
	f.Posts = 0
	fmt.Println("here return")
	return f,201,nil
}

func (fR *Repository) ForumDetails (slug string) (*models.Forum, error) {
	query := findForumJoinQuery
	forum := models.Forum{}
	err := fR.db.Get(&forum,query,slug)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(err)
	fmt.Println(forum)
	forum.Slug = slug
	return &forum, nil
}

func (fR *Repository) GetIdAndTitleBySlug (slug string) (*models.IdAndTitleForum, error) {
	query := GetForumIdAndTitle
	forumInfo := models.IdAndTitleForum{} 
	err := fR.db.Get(&forumInfo,query,slug)
	if err != nil {
		return nil, err
	}
	return &forumInfo,nil
}

func (fR *Repository) ForumSlugCreate (th *models.Thread, dopForumInfo *models.IdAndTitleForum, userId int) (*models.Thread, int, error) {
	query := createForumBranchQuery
	_,err := fR.db.Query(query,th.Title,th.Message,th.Created,userId,dopForumInfo.ID)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "thread_title_key"`) {
			fmt.Println(err.Error())
			query := getBranchInfo
			thread := models.Thread{}
			err := fR.db.Get(&thread,query,th.Title)
			if err != nil {
				fmt.Println(err)
				return nil, 500, err
			}
			fmt.Println(err)
			thread.Author = th.Author
			fmt.Println(thread)
			return &thread, 409, nil
		}
		fmt.Println(err)
		return nil, 500, err
	}
	th.Votes = 0
	return th, 200,nil
}

func (fR *Repository) GetThreadsByForum (info *models.ForumThreadsRequest) ([]models.Thread, error) {
	var query string
	if strings.ToLower(info.Desc) == "desc" {
		query = getThreadsByForumDesc
	} else {
		query = getThreadsByForumAsc
	}
	var threads []models.Thread
	rows,err := fR.db.Queryx(query,info.Slug,info.Limit)
	fmt.Println(&rows)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		thread := models.Thread{}
		rows.Scan(&thread.ID,&thread.Title,&thread.Author,&thread.Forum,&thread.Message,&thread.Votes,&thread.Slug,&thread.Created)
		fmt.Println(thread)
		threads = append(threads, thread)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return threads,nil
}
package repository

import (
	"bd_tp/models"
	"strings"

	sql "github.com/jmoiron/sqlx"
)

const (
	createForumQuery         = `insert into "forum" (title,slug,user_id) values ($1,$2,$3)`
	findForumQuery           = `select title,slug from "forum" where slug = $1`
	findForumJoinQuery       = `select f.id, f.title, f.slug, u.nickname from "forum" as f join "user" as u on f.user_id = u.id where f.slug= $1`
	findForumJoinIdQuery     = `select f.id, f.title, f.slug, u.nickname from "forum" as f join "user" as u on f.user_id = u.id where f.id= $1`
	createForumBranchQuery   = `insert into "thread" (title,message,created,slug,user_id,forum_id) values ($1,$2,$3,$4,$5,$6) returning id`
	GetForumIdAndTitle       = `select id,title from "forum" where slug = $1`
	getBranchInfo            = `select id,title,message,slug,created from "thread" where slug = $1`
	getBranchInfoWithoutSlug = `select id,title,message,created from "thread" where slug = $1`
	getThreadsByForumDesc    = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,t.slug,t.created from "thread" as t 
	join "forum" as f on t.forum_id = f.id  
	join "user" as u on u.id = t.user_id where f.slug = $1 
	ORDER BY t.created DESC
	LIMIT $2;`

	getThreadsByForumAsc = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,t.slug,t.created from "thread" as t 
	join "forum" as f on t.forum_id = f.id  
	join "user" as u on u.id = t.user_id where f.slug = $1 
	ORDER BY t.created ASC
	LIMIT $2;`

	getThreadsByForumDescSince = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,t.slug,t.created from "thread" as t 
	join "forum" as f on t.forum_id = f.id  
	join "user" as u on u.id = t.user_id where f.slug = $1 and t.created <=$2
	ORDER BY t.created DESC
	LIMIT $3;`

	getThreadsByForumAscSince = `select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,t.slug,t.created from "thread" as t 
	join "forum" as f on t.forum_id = f.id  
	join "user" as u on u.id = t.user_id where f.slug = $1 and t.created >=$2
	ORDER BY t.created ASC
	LIMIT $3;`

	createForumBranchQueryWithoutTimestamp = `insert into "thread" (title,message,slug,user_id,forum_id) values ($1,$2,$3,$4,$5) returning id`
)

type Repository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (fR *Repository) CreateForum(f *models.Forum, userID int) (*models.Forum, int, error) {
	query := createForumQuery
	_, err := fR.db.Query(query, f.Title, f.Slug, userID)
	if err != nil {

		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "forum_slug_key"`) {
			query := findForumQuery
			forum := models.Forum{}
			err := fR.db.Get(&forum, query, f.Slug)
			if err != nil {
				return nil, 500, err
			}
			return &forum, 409, nil
		}
		return nil, 500, err
	}
	f.Threads = 0
	f.Posts = 0
	return f, 201, nil
}

func (fR *Repository) ForumDetails(slug string) (*models.Forum, error) {
	query := findForumJoinQuery
	forum := models.Forum{}
	err := fR.db.Get(&forum, query, slug)
	if err != nil {

		return nil, err
	}
	return &forum, nil
}

func (fR *Repository) ForumDetailsById(id int) (*models.Forum, error) {
	query := findForumJoinIdQuery
	forum := models.Forum{}
	err := fR.db.Get(&forum, query, id)
	if err != nil {

		return nil, err
	}
	return &forum, nil
}

func (fR *Repository) GetIdAndTitleBySlug(slug string) (*models.IdAndTitleForum, error) {
	query := GetForumIdAndTitle
	forumInfo := models.IdAndTitleForum{}
	err := fR.db.Get(&forumInfo, query, slug)
	if err != nil {
		return nil, err
	}
	return &forumInfo, nil
}

func (fR *Repository) ForumSlugCreate(th *models.Thread, dopForumInfo *models.Forum, userId int) (*models.Thread, int, error) {
	query := createForumBranchQuery
	var threadId int
	err := fR.db.Get(&threadId, query, th.Title, th.Message, th.Created, th.Slug, userId, dopForumInfo.ID)
	if err != nil {

		query := getBranchInfo
		thread := models.Thread{}
		err := fR.db.Get(&thread, query, th.Slug)
		if err != nil {

			return nil, 500, err
		}

		thread.Author = th.Author
	
		return &thread, 409, nil
	}
	th.ID = threadId
	return th, 201, nil
}

func (fR *Repository) ForumSlugCreateWithoutTimeStamp(th *models.Thread, dopForumInfo *models.Forum, userId int) (*models.Thread, int, error) {
	query := createForumBranchQueryWithoutTimestamp
	var threadId int
	err := fR.db.Get(&threadId, query, th.Title, th.Message, th.Slug, userId, dopForumInfo.ID)
	if err != nil {
		query := getBranchInfo
		thread := models.Thread{}
		err := fR.db.Get(&thread, query, th.Slug)
		if err != nil {

			return nil, 500, err
		}

		thread.Author = th.Author
		return &thread, 409, nil
	}
	th.ID = threadId
	return th, 201, nil
}

func (fR *Repository) GetThreadsByForum(slug, limit, since, desc string) ([]models.Thread, error) {
	var query string
	var rows *sql.Rows
	var err error
	if strings.ToLower(desc) == "true" {
		if since == "" {
			query = getThreadsByForumDesc
			rows, err = fR.db.Queryx(query, slug, limit)
		} else {
			query = getThreadsByForumDescSince
			rows, err = fR.db.Queryx(query, slug, since, limit)
		}
	} else {
		if since == "" {
			query = getThreadsByForumAsc
			rows, err = fR.db.Queryx(query, slug, limit)
		} else {
			query = getThreadsByForumAscSince
			rows, err = fR.db.Queryx(query, slug, since, limit)
		}
	}
	var threads []models.Thread

	if err != nil {

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		thread := models.Thread{}
		rows.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		threads = append(threads, thread)
	}
	if err != nil {

		return nil, err
	}
	return threads, nil
}

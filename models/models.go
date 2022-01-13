package models

type User struct {
	ID int `db:"id"`
	Nickname string `db:"nickname"`
	Fullname string `db:"fullname"`
	About string  `db:"about"`
	Email string `db:"email"`
}

type Forum struct {
	ID int `db:"id"`
	Title string `db:"title"`
	User string	`db:"nickname"`
	Slug string `db:"slug"`
	Posts int
	Threads int
}

type Thread struct {
	ID int `db:"id" json:"id,omitempty"`
	Title string `db:"title" json:"title,omitempty"`
	Author string `db:"author" json:"author,omitempty"`
	Forum string `db:"forum" json:"forum,omitempty"`
	Message string `db:"message" json:"message,omitempty"`
	Votes int `db:"votes" json:"votes,omitempty"`
	Slug string `db:"slug" json:"slug,omitempty"`
	Created string `db:"created" json:"created,omitempty"`
}

type ForumThreadsRequest struct {
	Slug string 
	Limit string 
	Since string 
	Desc string 
}

type Threads struct {
	Threads []Thread `json:"threads,omitempty"`
}

type Post struct {
	ID int `db:"id"`
	Parent int `db:"parent"` 
	Author string `db:"author"` 
	Message string `db:"message"`
	Edited bool `db:"edited"`
	Forum string `db:"forum"`
	Thread int `db:"thread"`
	Created string `db:"created"`
}

type Posts struct {
	Posts []Post `json:"posts,omitempty"`
}

type Vote struct {
	Nickname string `json:"nickname,omitempty"`
	Voice int 	`json:"voice,omitempty"`
}

type PostsRelated struct {
	Related []string `json:"related,omitempty"`
}

type FullPostInfo struct {
	Post Post
	Author User
	Thread Thread
	Forum Forum
}
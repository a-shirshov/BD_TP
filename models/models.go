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
	Posts int `db:"posts"`
	Threads int `db:"threads"`
}

type Thread struct {
	ID int `db:"id"` 
	Title string `db:"title"` 
	Author string `db:"author"` 
	Forum string `db:"forum"`
	Message string `db:"message"`
	Votes int `db:"votes"`
	Slug string `db:"slug"`
	Created string `db:"created"`
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
	Post *Post
	Author *User
	Thread *Thread
	Forum *Forum
}

type PostMessage struct {
	Message string `db:"message"`
}

type Status struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	Post   int `json:"post"`
}

package response

type Error struct {
	Message string `json:"message,omitempty"`
}

type UserResponse struct {
	ID *int `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About string `json:"about,omitempty"`
	Email string `json:"email,omitempty"`
}

type ForumResponse struct {
	ID *int `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	User string `json:"user,omitempty"`
	Slug string `json:"slug,omitempty"`
	Posts *int `json:"posts,omitempty"`
	Threads *int `json:"threads,omitempty"`
}

type ThreadResponse struct {
	ID *int `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Forum string `json:"forum,omitempty"`
	Message string `json:"message,omitempty"`
	Votes *int `json:"votes,omitempty"`
	Slug string `json:"slug,omitempty"`
	Created string `json:"created,omitempty"`
}

type ForumThreadsRequest struct {
	Limit string `json:"limit,omitempty"`
	Since string `json:"since,omitempty"`
	Desc string `json:"desc,omitempty"`
}

type ThreadsResponse struct {
	Threads []ThreadResponse `json:""`
}

type PostsRequest struct {
	Posts []PostRequest `json:"posts,omitempty"`
}

type UsersResponse struct {
	Users [] UserResponse `json:""`
}

type PostRequest struct {
	Parent int `json:"parent,omitempty"`
	Author string `json:"author,omitempty"`
	Message string `json:"message,omitempty"`
}

type VoteRequest struct {
	Nickname string `json:"nickname,omitempty"`
	Voice *int `json:"voice,omitempty"`
}

type PostRelated struct {
	Related []string `json:"related,omitempty"`
}

type FullPostInfo struct {
	Post PostResponse `json:"post,omitempty"`	
	Author UserResponse `json:"author,omitempty"`
	Thread ThreadResponse `json:"thread,omitempty"`
	Forum ForumResponse `json:"forum,omitempty"`
}

type PostResponse struct {
	ID *int `json:"id,omitempty"`
	Parent *int `json:"parent,omitempty"`
	Author string `json:"author,omitempty"`
	Message string `json:"message,omitempty"`
	Edited *bool `json:"isEdited,omitempty"`
	Forum string `json:"forum,omitempty"`
	Thread *int `json:"thread,omitempty"`
	Created string `json:"created,omitempty"`
}
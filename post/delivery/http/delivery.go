package delivery

import (
	postUsecase "bd_tp/post/usecase"
	"bd_tp/response"
	"net/http"
	"strconv"
	"strings"
)

type PostDelivery struct {
	PostUsecase *postUsecase.Usecase
}

func NewPostDelivery(pU *postUsecase.Usecase) *PostDelivery {
	return &PostDelivery{
		PostUsecase: pU,
	}
}

func (pD *PostDelivery) PostDetails (w http.ResponseWriter, r* http.Request) {
	path := r.URL.Path
	split := strings.Split(path,"/")
	id := split[len(split)-2]

	q := r.URL.Query()
	var related string
	if len(q["related"]) > 0 {
		related = q["related"][0]
	}

	fullPost,err := pD.PostUsecase.PostDetails(id,related)
	if err != nil {
		errorResponse := &response.Error{
			Message: "Error",
		}
		response.SendResponse(w,404,errorResponse)
		return
	}

	var postResponse *response.PostResponse
	if fullPost.Post != nil{
		postResponse = &response.PostResponse{
			ID: &fullPost.Post.ID,
			Parent: &fullPost.Post.Parent,
			Author: fullPost.Post.Author,
			Message: fullPost.Post.Message,
			Edited: &fullPost.Post.Edited,
			Forum: fullPost.Post.Forum,
			Thread: &fullPost.Post.Thread,
			Created: fullPost.Post.Created,
		}
	}

	var authorResponse *response.UserResponse
	if fullPost.Author != nil {
		authorResponse = &response.UserResponse{
			Nickname: fullPost.Author.Nickname,
			Fullname: fullPost.Author.Fullname,
			About: fullPost.Author.About,
			Email: fullPost.Author.Email,
		}
	}
	
	var threadResponse *response.ThreadResponse
	if fullPost.Thread != nil {
		threadResponse = &response.ThreadResponse{
			ID: &fullPost.Thread.ID,
			Title: fullPost.Thread.Title,
			Author: fullPost.Thread.Author,
			Forum: fullPost.Thread.Forum,
			Message: fullPost.Thread.Message,
			Votes: &fullPost.Thread.Votes,
			Slug: fullPost.Thread.Slug,
			Created: fullPost.Thread.Created,
		}
	}

	var forumResponse *response.ForumResponse
	if fullPost.Forum != nil {
		forumResponse = &response.ForumResponse{
			Title: fullPost.Forum.Title,
			User: fullPost.Forum.User,
			Slug: fullPost.Forum.Slug,
			Posts: &fullPost.Forum.Posts,
			Threads: &fullPost.Forum.Threads,
		}
	}

	fullPostResponse := response.FullPostInfo {
		Post: postResponse,
		Author: authorResponse,
		Thread: threadResponse,
		Forum: forumResponse,
	}
	response.SendResponse(w,200,fullPostResponse)
}

func (pD *PostDelivery) UpdatePost (w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	split := strings.Split(path,"/")
	id := split[len(split)-2]

	post, err := response.GetPostFromRequest(r.Body)
	if err != nil {
		return 
	}

	idInt,err := strconv.Atoi(id)
	if err != nil {
		return
	}
	post.ID = idInt
	updatedPost, err := pD.PostUsecase.UpdatePost(post)
	if err != nil {
		errorResponse := &response.Error{
			Message: "Error",
		}
		response.SendResponse(w,404,errorResponse)
		return
	}
	postResponse := &response.PostResponse{
		ID: &updatedPost.ID,
		Parent: &updatedPost.Parent,
		Author: updatedPost.Author,
		Message: updatedPost.Message,
		Edited: &updatedPost.Edited,
		Forum: updatedPost.Forum,
		Thread: &updatedPost.Thread,
		Created: updatedPost.Created,
	}
	response.SendResponse(w,200,postResponse)
}
package delivery

import (
	postUsecase "bd_tp/post/usecase"
	"bd_tp/response"
	"fmt"
	"net/http"
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
	postsRequest, err := response.GetPostRelatedFromRequest(r.Body)
	fmt.Println(postsRequest.Related)
	if err != nil {
		return
	}
	path := r.URL.Path
	split := strings.Split(path,"/")
	id := split[len(split)-2]

	requestString := strings.Join(postsRequest.Related," ")
	fmt.Println(requestString)

	wannaUser := strings.Contains(requestString,"user")
	wannaForum := strings.Contains(requestString,"forum")
	wannaThread := strings.Contains(requestString,"thread")

	fullPost,err := pD.PostUsecase.PostDetails(id,wannaUser,wannaForum,wannaThread)
	if err != nil {
		errorResponse := &response.Error{
			Message: "Error",
		}
		response.SendResponse(w,404,errorResponse)
	}

	postResponse := &response.PostResponse{
		ID: &fullPost.Post.ID,
		Parent: &fullPost.Post.Parent,
		Author: fullPost.Post.Author,
		Message: fullPost.Post.Message,
		Edited: &fullPost.Post.Edited,
		Forum: fullPost.Post.Forum,
		Thread: &fullPost.Post.Thread,
		Created: fullPost.Post.Created,
	}

	authorResponse := &response.UserResponse{
		Nickname: fullPost.Author.Nickname,
		Fullname: fullPost.Author.Fullname,
		About: fullPost.Author.About,
		Email: fullPost.Author.Email,
	}

	threadResponse := &response.ThreadResponse{
		ID: &fullPost.Thread.ID,
		Title: fullPost.Thread.Title,
		Author: fullPost.Thread.Author,
		Forum: fullPost.Thread.Forum,
		Message: fullPost.Thread.Message,
		Votes: &fullPost.Thread.Votes,
		Slug: fullPost.Thread.Slug,
		Created: fullPost.Thread.Created,
	}

	forumResponse := &response.ForumResponse{
		Title: fullPost.Forum.Title,
		User: fullPost.Forum.User,
		Slug: fullPost.Forum.Slug,
		Posts: &fullPost.Forum.Posts,
		Threads: &fullPost.Forum.Threads,
	}

	fullPostResponse := response.FullPostInfo {
		Post: *postResponse,
		Author: *authorResponse,
		Thread: *threadResponse,
		Forum: *forumResponse,
	}
	response.SendResponse(w,200,fullPostResponse)
}
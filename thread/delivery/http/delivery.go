package delivery

import (
	"bd_tp/models"
	"bd_tp/response"
	threadUsecase "bd_tp/thread/usecase"
	"fmt"
	"net/http"
	"strings"
)

type ThreadDelivery struct {
	threadU *threadUsecase.Usecase
}

func NewThreadDelivery (tU *threadUsecase.Usecase) *ThreadDelivery {
	return &ThreadDelivery{
		threadU: tU,
	}
}

func (tD *ThreadDelivery) CreatePosts (w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	postsRequest, err := response.GetPostsFromRequest(r.Body)
	fmt.Println(postsRequest)
	if err != nil {
		return
	}

	path := r.URL.Path
	split := strings.Split(path,"/")
	slug_or_id := split[len(split)-2]

	posts,code,err := tD.threadU.CreatePosts(postsRequest,slug_or_id)
	fmt.Println(posts)
	if err != nil {
		errorResponse := &response.Error{
			Message: "Some troubles",
		}
		response.SendResponse(w,code,errorResponse)
	}
	var postsResponse models.Posts
	postsResponse.Posts = posts
	response.SendResponse(w,code,postsResponse)
}

func (tD *ThreadDelivery) ThreadDetails (w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	split := strings.Split(path,"/")
	slug_or_id := split[len(split)-2]

	thread,err := tD.threadU.ThreadDetails(slug_or_id)
	if err != nil {
		errorResponse := &response.Error{
			Message: "Id or slug problem",
		}
		response.SendResponse(w,404,errorResponse)
		return
	}
	threadResponse := &response.ThreadResponse{
		ID: &thread.ID,
		Title: thread.Title,
		Author: thread.Author,
		Forum: thread.Forum,
		Message: thread.Message,
		Votes: &thread.Votes,
		Slug: thread.Slug,
		Created: thread.Created,
	}
	response.SendResponse(w,200,threadResponse)
}

func (tD *ThreadDelivery) ThreadDetailsUpdate (w http.ResponseWriter, r *http.Request) {
	threadRequest, err := response.GetThreadUpdateFromRequest(r.Body)
	if err != nil {
		return
	}
	path := r.URL.Path
	split := strings.Split(path,"/")
	slug_or_id := split[len(split)-2]

	thread,err := tD.threadU.ThreadDetailsUpdate(threadRequest,slug_or_id)
	if err != nil {
		errorResponse := &response.Error{
			Message: "Id or slug problem",
		}
		response.SendResponse(w,404,errorResponse)
		return
	}
	threadResponse := &response.ThreadResponse{
		ID: &thread.ID,
		Title: thread.Title,
		Author: thread.Author,
		Forum: thread.Forum,
		Message: thread.Message,
		Votes: &thread.Votes,
		Slug: thread.Slug,
		Created: thread.Created,
	}
	response.SendResponse(w,200,threadResponse)
}

func (tD *ThreadDelivery) ThreadVote (w http.ResponseWriter, r *http.Request) {
	voteRequest, err := response.GetVoteFromRequest(r.Body)
	if err != nil {
		return
	}
	path := r.URL.Path
	split := strings.Split(path,"/")
	slug_or_id := split[len(split)-2]

	thread,err := tD.threadU.ThreadVote(voteRequest,slug_or_id)

	if err != nil {
		errorResponse := &response.Error{
			Message: "Id or slug problem",
		}
		response.SendResponse(w,404,errorResponse)
		return
	}

	threadResponse := &response.ThreadResponse{
		ID: &thread.ID,
		Title: thread.Title,
		Author: thread.Author,
		Forum: thread.Forum,
		Message: thread.Message,
		Votes: &thread.Votes,
		Slug: thread.Slug,
		Created: thread.Created,
	}
	response.SendResponse(w,200,threadResponse)
}
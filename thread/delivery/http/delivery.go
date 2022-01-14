package delivery

import (
	"bd_tp/response"
	threadUsecase "bd_tp/thread/usecase"
	"net/http"
	"strings"
	"fmt"
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
	postsRequest, err := response.GetPostsFromRequest(r.Body)
	if err != nil {
		return
	}

	path := r.URL.Path
	split := strings.Split(path,"/")
	slug_or_id := split[len(split)-2]

	posts,code,err := tD.threadU.CreatePosts(postsRequest,slug_or_id)
	if err != nil {
		fmt.Println(err)
		errorResponse := &response.Error{
			Message: "Some troubles",
		}
		response.SendResponse(w,code,errorResponse)
		return
	}
	var postsResponse []response.PostResponse
	for index := range posts {
		postResponse := &response.PostResponse{
			ID: &posts[index].ID,
			Parent: &posts[index].Parent,
			Author: posts[index].Author,
			Message: posts[index].Message,
			Edited: &posts[index].Edited,
			Forum: posts[index].Forum,
			Thread: &posts[index].Thread,
			Created: posts[index].Created,
		}
		postsResponse = append(postsResponse, *postResponse)
	}
	if len(postsResponse) == 0 {
		response.SendResponse(w, code, []response.PostResponse{})
		return
	}
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

func (tD *ThreadDelivery) ThreadGetPosts (w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	split := strings.Split(path,"/")
	slug_or_id := split[len(split)-2]

	q := r.URL.Query()
	var limit string
	var since string
	var sort string
	var desc string
	if len(q["limit"]) > 0 {
		limit = q["limit"][0]
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["sort"]) > 0 {
		sort = q["sort"][0]
	}
	if len(q["desc"]) > 0 {
		desc = q["desc"][0]
	}

	posts, err := tD.threadU.ThreadGetPosts(slug_or_id, limit,since,sort,desc)
	if err != nil {
		errorResponse := &response.Error{
			Message: "Problems",
		}
		response.SendResponse(w,404,errorResponse)
		return
	}

	var postsResponse []response.PostResponse
	for index := range posts {
		postResponse := &response.PostResponse{
			ID: &posts[index].ID,
			Parent: &posts[index].Parent,
			Author: posts[index].Author,
			Message: posts[index].Message,
			Edited: &posts[index].Edited,
			Forum: posts[index].Forum,
			Thread: &posts[index].Thread,
			Created: posts[index].Created,
		}
		postsResponse = append(postsResponse, *postResponse)
	}
	if len(postsResponse) == 0 {
		response.SendResponse(w, 200, []response.ThreadResponse{})
		return
	}
	response.SendResponse(w,200,postsResponse)

}
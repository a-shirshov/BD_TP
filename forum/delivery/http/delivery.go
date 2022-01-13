package delivery

import (
	forumUsecase "bd_tp/forum/usecase"
	"bd_tp/response"
	"fmt"
	"net/http"
	"strings"
)

type ForumDelivery struct {
	ForumUsecase *forumUsecase.Usecase
}

func NewForumDelivery(fU *forumUsecase.Usecase) *ForumDelivery {
	return &ForumDelivery{
		ForumUsecase: fU,
	}
}

func (fD *ForumDelivery) CreateForum (w http.ResponseWriter, r *http.Request) {
	f, err := response.GetForumFromRequest(r.Body)
	if err != nil {
		return
	}

	forum, code, err := fD.ForumUsecase.CreateForum(f)
	if err != nil {
		if code == 404 {
			errorResponse := &response.Error{
				Message: "No user with this nickname:"+f.User,
			}
			response.SendResponse(w,code,errorResponse)
		}
		return
	}

	forumResponse := &response.ForumResponse{
		Title: forum.Title,
		User: forum.User,
		Slug: forum.Slug,
		Posts: &forum.Posts,
		Threads: &forum.Threads,
	}
	response.SendResponse(w,code,forumResponse)
}

func (fD *ForumDelivery) ForumDetails (w http.ResponseWriter, r *http.Request) {
	
	path := r.URL.Path
	split := strings.Split(path,"/")
	slug := split[len(split)-2]
	
	forum, err := fD.ForumUsecase.ForumDetails(slug)
	if err != nil {
		errorResponse := &response.Error{
			Message: "No forum with this slug:"+slug,
		}
		response.SendResponse(w,404,errorResponse)
		return
	}
	forumResponse := &response.ForumResponse{
		Title: forum.Title,
		User: forum.User,
		Slug: forum.Slug,
		Posts: &forum.Posts,
		Threads: &forum.Threads,
	}
	response.SendResponse(w,200,forumResponse)
}

func (fD *ForumDelivery) ForumSlugCreate (w http.ResponseWriter, r *http.Request) {
	th, err := response.GetThreadFromRequest(r.Body)
	if err != nil {
		fmt.Println(err)
		return 
	}
	path := r.URL.Path
	split := strings.Split(path,"/")
	slug := split[len(split)-2]
	th.Slug = slug

	thread,code, err := fD.ForumUsecase.ForumSlugCreate(th)
	fmt.Println(thread)
	if err != nil {
		fmt.Println(err)
		errorResponse := &response.Error{
			Message: "Mistake with author or slug",
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
	response.SendResponse(w,code,threadResponse)
}

func (fD *ForumDelivery) GetThreadsByForum (w http.ResponseWriter, r *http.Request) {
	Info,err := response.GetThreadsQueryInfo(r.Body)
	if err != nil {
		return
	}

	path := r.URL.Path
	split := strings.Split(path,"/")
	slug := split[len(split)-2]
	Info.Slug = slug

	threads,err := fD.ForumUsecase.GetThreadsByForum(Info)
	if err != nil {
		fmt.Println(err)
		errorResponse := &response.Error{
			Message: "No forum with with slug:"+Info.Slug,
		}
		response.SendResponse(w,404,errorResponse)
		return
	}
	fmt.Println(threads)
	var threadsResponse response.ThreadsResponse
	for _,thread := range threads {
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
		threadsResponse.Threads = append(threadsResponse.Threads, *threadResponse)
	}
	response.SendResponse(w,200,threadsResponse)
}
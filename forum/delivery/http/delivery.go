package delivery

import (
	forumUsecase "bd_tp/forum/usecase"
	"bd_tp/response"
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
		
		return 
	}
	path := r.URL.Path
	split := strings.Split(path,"/")
	slug := split[len(split)-2]
	thread,code, err := fD.ForumUsecase.ForumSlugCreate(th,slug)
	if err != nil {
		
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

	path := r.URL.Path
	split := strings.Split(path,"/")
	slug := split[len(split)-2]

	q := r.URL.Query()
	var limit string
	var since string
	var desc string
	if len(q["limit"]) > 0 {
		limit = q["limit"][0]
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["desc"]) > 0 {
		desc = q["desc"][0]
	}

	threads,err := fD.ForumUsecase.GetThreadsByForum(slug,limit,since,desc)
	if err != nil {
		
		errorResponse := &response.Error{
			Message: "No forum with with slug:"+slug,
		}
		response.SendResponse(w,404,errorResponse)
		return
	}
	var threadsResponse []response.ThreadResponse
	for index := range threads {
		threadResponse := &response.ThreadResponse{
			ID: &threads[index].ID,
			Title: threads[index].Title,
			Author: threads[index].Author,
			Forum: threads[index].Forum,
			Message: threads[index].Message,
			Votes: &threads[index].Votes,
			Slug: threads[index].Slug,
			Created: threads[index].Created,
		}
		threadsResponse = append(threadsResponse, *threadResponse)
	}
	if len(threadsResponse) == 0 {
		response.SendResponse(w, 200, []response.ThreadResponse{})
		return
	}
	response.SendResponse(w,200,threadsResponse)
}

func (fD *ForumDelivery) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	split := strings.Split(path,"/")
	slug := split[len(split)-2]

	q := r.URL.Query()
	var limit string
	var since string
	var desc string
	if len(q["limit"]) > 0 {
		limit = q["limit"][0]
	}
	if len(q["since"]) > 0 {
		since = q["since"][0]
	}
	if len(q["desc"]) > 0 {
		desc = q["desc"][0]
	}

	users, err := fD.ForumUsecase.GetForumUsers(slug,limit,since,desc)
	if err != nil {
		
		errorResponse := &response.Error{
			Message: "No forum with with slug:"+slug,
		}
		response.SendResponse(w,404,errorResponse)
		return
	}
	var usersResponse []response.UserResponse
	for _, user := range users {
		userResponse := &response.UserResponse{
			Nickname: user.Nickname,
			Fullname: user.Fullname,
			About:    user.About,
			Email:    user.Email,
		}
		usersResponse = append(usersResponse, *userResponse)
	}
	if len(usersResponse) == 0 {
		response.SendResponse(w, 200, []response.ThreadResponse{})
		return
	}
	response.SendResponse(w,200,usersResponse)
}
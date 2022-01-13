package register

import (
	"github.com/gorilla/mux"
	userD "bd_tp/user/delivery/http"
	forumD "bd_tp/forum/delivery/http"
	threadD "bd_tp/thread/delivery/http"
	postD "bd_tp/post/delivery/http"
)

func UserEndpoints(r *mux.Router, userD *userD.UserDelivery) {
	r.HandleFunc("/{nickname}/create",userD.CreateUser).Methods("POST")
	r.HandleFunc("/{nickname}/profile",userD.ProfileInfo).Methods("GET")
	r.HandleFunc("/{nickname}/profile",userD.UpdateProfile).Methods("POST")
}

func ForumEndpoints(r *mux.Router, forumD *forumD.ForumDelivery) {
	r.HandleFunc("/create",forumD.CreateForum).Methods("POST")
	r.HandleFunc("/{slug}/details",forumD.ForumDetails).Methods("GET")
	r.HandleFunc("/{slug}/create", forumD.ForumSlugCreate).Methods("POST")
	r.HandleFunc("/{slug}/threads", forumD.GetThreadsByForum).Methods("GET")
}

func ThreadEndpoints(r *mux.Router, threadD *threadD.ThreadDelivery) {
	r.HandleFunc("/{slug_or_id}/create",threadD.CreatePosts).Methods("POST") 
	r.HandleFunc("/{slug_or_id}/details",threadD.ThreadDetails).Methods("GET")
	r.HandleFunc("/{slug_or_id}/details",threadD.ThreadDetailsUpdate).Methods("POST")
	r.HandleFunc("/{slug_or_id}/vote",threadD.ThreadVote).Methods("POST")
}

func PostEndpoints(r *mux.Router, postD *postD.PostDelivery) {
	r.HandleFunc("/{id}/details",postD.PostDetails).Methods("GET")
}
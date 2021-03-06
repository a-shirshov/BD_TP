package main

import (
	"bd_tp/register"
	userDelivery "bd_tp/user/delivery/http"
	userRepo "bd_tp/user/repository"
	userUsecase "bd_tp/user/usecase"

	forumDelivery "bd_tp/forum/delivery/http"
	forumRepo "bd_tp/forum/repository"
	forumUsecase "bd_tp/forum/usecase"

	threadDelivery "bd_tp/thread/delivery/http"
	threadRepo "bd_tp/thread/repository"
	threadUsecase "bd_tp/thread/usecase"

	postDelivery "bd_tp/post/delivery/http"
	postRepo "bd_tp/post/repository"
	postUsecase "bd_tp/post/usecase"

	serviceDelivery "bd_tp/service/delivery/http"
	serviceRepo "bd_tp/service/repository"
	serviceUsecase "bd_tp/service/usecase"

	"bd_tp/utils"
	"fmt"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {

	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		os.Exit(1)
	}

	db, err := utils.InitPostgresDB()
	if err != nil {
		fmt.Printf("%s",err)
		return
	}

	userR := userRepo.NewRepository(db)
	forumR := forumRepo.NewForumRepository(db)
	threadR := threadRepo.NewThreadRepository(db)
	postR := postRepo.NewPostRepository(db)
	serviceR := serviceRepo.NewServiceRepository(db)

	userU := userUsecase.NewUserUsecase(userR)
	forumU := forumUsecase.NewForumUsecase(forumR,userR,threadR)
	threadU := threadUsecase.NewThreadUsecase(threadR,postR,userR,forumR)
	postU := postUsecase.NewPostUsecase(postR,userR,forumR,threadR)
	serviceU := serviceUsecase.NewServiceUseCase(serviceR)

	userD := userDelivery.NewUserDelivery(userU)
	forumD := forumDelivery.NewForumDelivery(forumU)
	threadD := threadDelivery.NewThreadDelivery(threadU)
	postD := postDelivery.NewPostDelivery(postU)
	serviceD := serviceDelivery.NewServiceDelivery(serviceU)
	

	r := mux.NewRouter()
	rApi := r.PathPrefix("/api").Subrouter()
	userRouter := rApi.PathPrefix("/user").Subrouter()
	register.UserEndpoints(userRouter,userD)
	forumRouter := rApi.PathPrefix("/forum").Subrouter()
	register.ForumEndpoints(forumRouter,forumD)
	threadRouter := rApi.PathPrefix("/thread").Subrouter()
	register.ThreadEndpoints(threadRouter,threadD)
	postRouter := rApi.PathPrefix("/post").Subrouter()
	register.PostEndpoints(postRouter,postD)
	serviceRouter := rApi.PathPrefix("/service").Subrouter()
	register.ServiceEndpoints(serviceRouter,serviceD)
	
	err = http.ListenAndServe(":5000", r)
	if err != nil {
		fmt.Printf("%s",err)
		return 
	}
}
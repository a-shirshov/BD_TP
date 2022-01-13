package usecase

import (
	forumRepo "bd_tp/forum/repository"
	userRepo "bd_tp/user/repository"
	"bd_tp/models"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

type Usecase struct {
	forumRepo *forumRepo.Repository
	userRepo *userRepo.Repository
}

func NewForumUsecase (fR *forumRepo.Repository, uR *userRepo.Repository) *Usecase {
	return &Usecase{
		forumRepo: fR,
		userRepo: uR,
	}
}

func (fU *Usecase) CreateForum(f *models.Forum) (*models.Forum, int, error) {
	userId, err := fU.userRepo.GetIdByNickname(f.User)
	if err != nil {
		return nil,404,err
	}
	forum, code, err := fU.forumRepo.CreateForum(f,userId)
	if err != nil {
		return nil, code, err
	}
	return forum,code,err
}

func (fU *Usecase) ForumDetails(slug string) (*models.Forum,error) {
	forum,err := fU.forumRepo.ForumDetails(slug)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (fU *Usecase) ForumSlugCreate (th *models.Thread) (*models.Thread,int,error) {
	userId, err := fU.userRepo.GetIdByNickname(th.Author)
	if err != nil {
		fmt.Println(err)
		return nil,404,err
	}
	forumInfo, err := fU.forumRepo.GetIdAndTitleBySlug(th.Slug)
	if err !=nil {
		fmt.Println(err)
		return nil,404,err
	}
	slug := forumInfo.Title+"/thread" + uuid.NewV4().String()
	th.Slug = slug
	thread, code, err := fU.forumRepo.ForumSlugCreate(th,forumInfo,userId)
	if err != nil {
		fmt.Println(err)
		return nil, code, err
	}
	thread.Forum = forumInfo.Title
	return thread,code,err
}

func (fU *Usecase) GetThreadsByForum (info *models.ForumThreadsRequest) ([]models.Thread, error) {
	//Прверка на наличие форума
	_, err := fU.forumRepo.GetIdAndTitleBySlug(info.Slug)
	if err !=nil {
		fmt.Println(err)
		return nil,err
	}
	//Здесь форум есть, но мб нет постов - не ошибка
	threads, err := fU.forumRepo.GetThreadsByForum(info)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return threads,nil
}	
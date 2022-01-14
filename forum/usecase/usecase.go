package usecase

import (
	forumRepo "bd_tp/forum/repository"
	"bd_tp/models"
	threadRepo "bd_tp/thread/repository"
	userRepo "bd_tp/user/repository"
)

type Usecase struct {
	forumRepo *forumRepo.Repository
	userRepo  *userRepo.Repository
	threadRepo *threadRepo.Repository
}

func NewForumUsecase(fR *forumRepo.Repository, uR *userRepo.Repository, tR *threadRepo.Repository) *Usecase {
	return &Usecase{
		forumRepo: fR,
		userRepo:  uR,
		threadRepo: tR,
	}
}

func (fU *Usecase) CreateForum(f *models.Forum) (*models.Forum, int, error) {
	user, err := fU.userRepo.ProfileInfo(f.User)
	if err != nil {
		return nil, 404, err
	}
	if f.Slug != "" {
		oldForum, err := fU.forumRepo.ForumDetails(f.Slug)
		if err == nil {
			return oldForum, 409, nil
		}
	}
	forum, code, err := fU.forumRepo.CreateForum(f, user.ID)
	forum.User = user.Nickname
	if err != nil {
		return nil, code, err
	}
	return forum, code, err
}

func (fU *Usecase) ForumDetails(slug string) (*models.Forum, error) {
	forum, err := fU.forumRepo.ForumDetails(slug)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (fU *Usecase) ForumSlugCreate(th *models.Thread, slug string) (*models.Thread, int, error) {
	userId, err := fU.userRepo.GetIdByNickname(th.Author)
	if err != nil {

		return nil, 404, err
	}

	forum, err := fU.forumRepo.ForumDetails(slug)
	if err != nil {

		return nil, 404, err
	}

	if th.Slug != "" {
		oldThread, err := fU.threadRepo.ThreadDetailsBySlug(th.Slug)
		if err != nil {
		} else {
			return oldThread,409,nil
		}
	}

	if th.Created == "" {
		thread, code, err := fU.forumRepo.ForumSlugCreateWithoutTimeStamp(th, forum, userId)
		if err != nil {

			return nil, code, err
		}
		thread.Forum = forum.Slug
		return thread, code, err
	}

	thread, code, err := fU.forumRepo.ForumSlugCreate(th, forum, userId)
	if err != nil {

		return nil, code, err
	}
	thread.Forum = forum.Slug
	return thread, code, err
}

func (fU *Usecase) GetThreadsByForum(slug, limit, since, desc string) ([]models.Thread, error) {
	//Прверка на наличие форума
	_, err := fU.forumRepo.GetIdAndTitleBySlug(slug)
	if err != nil {

		return nil, err
	}
	if limit == "" {
		limit = "100"
	}
	//Здесь форум есть, но мб нет постов - не ошибка
	threads, err := fU.forumRepo.GetThreadsByForum(slug, limit, since, desc)
	if err != nil {

		return nil, err
	}
	return threads, nil
}

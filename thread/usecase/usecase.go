package usecase

import (
	"bd_tp/models"
	threadRepo "bd_tp/thread/repository"
	"strconv"
)

type Usecase struct {
	threadRepo *threadRepo.Repository
}

func NewThreadUsecase (tR *threadRepo.Repository) *Usecase {
	return &Usecase{
		threadRepo: tR,
	}
}

func (fR *Usecase) CreatePosts(postsInput []models.Post, slug_or_id string) ([]models.Post, int, error) {
	id, err := strconv.Atoi(slug_or_id) 
	var posts []models.Post
	if err != nil {
		posts, code, err := fR.threadRepo.CreatePostsWithSlug(postsInput,slug_or_id)
		if err != nil {
			return nil, code, err
		}
		return posts, code, nil
	}
	posts, code, err := fR.threadRepo.CreatePostsWithID(postsInput,id)
	if err != nil {
		return nil, code, err
	}
	return posts, code, nil
}

func (fR *Usecase) ThreadDetails(slug_or_id string) (*models.Thread,error) {
	id, err := strconv.Atoi(slug_or_id) 
	if err != nil {
		thread, err := fR.threadRepo.ThreadDetailsBySlug(slug_or_id)
		if err != nil {
			return nil,err
		}
		return thread, nil
	}
	thread, err := fR.threadRepo.ThreadDetailsByID(id)
	if err != nil {
		return nil, err
	}
	return thread,nil
}

func (fR *Usecase) ThreadDetailsUpdate(threadInfo *models.Thread, slug_or_id string) (*models.Thread, error) {
	id, err := strconv.Atoi(slug_or_id) 
	if err != nil {
		thread, err := fR.threadRepo.ThreadDetailsUpdateBySlug(threadInfo,slug_or_id)
		if err != nil {
			return nil, err
		}
		return thread,nil
	}
	thread, err := fR.threadRepo.ThreadDetailsUpdateByID(threadInfo,id)
	if err != nil {
		return nil, err
	}
	return thread,nil
}

func (fR *Usecase) ThreadVote(vote *models.Vote, slug_or_id string) (*models.Thread, error) {
	id, err := strconv.Atoi(slug_or_id) 
	if err != nil {
		thread, err := fR.threadRepo.ThreadVoteBySlug(vote,slug_or_id)
		if err != nil {
			return nil, err
		}
		return thread,nil
	}
	thread, err := fR.threadRepo.ThreadVoteByID(vote,id)
	if err != nil {
		return nil, err
	}
	return thread,nil
}
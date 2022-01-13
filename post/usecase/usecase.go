package usecase

import (
	"bd_tp/models"
	postRepo "bd_tp/post/repository"
	userRepo "bd_tp/user/repository"
	forumRepo "bd_tp/forum/repository"
	threadRepo "bd_tp/thread/repository"
	//"fmt"
	"strconv"
)

type Usecase struct {
	postRepo *postRepo.Repository
	userRepo *userRepo.Repository
	forumRepo *forumRepo.Repository
	threadRepo *threadRepo.Repository
}

func NewPostUsecase (pR *postRepo.Repository, uR *userRepo.Repository, fR *forumRepo.Repository, tR *threadRepo.Repository) *Usecase {
	return &Usecase{
		postRepo: pR,
		userRepo: uR,
		forumRepo: fR,
		threadRepo: tR,
	}
}

func (pU *Usecase) PostDetails(idStr string,wannaUser,wannaForum,wannaThread bool) (*models.FullPostInfo, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil,err
	}

	var FullPostInfo models.FullPostInfo
	post, err := pU.postRepo.GetPostByID(id)
	if err != nil {
		return nil,err
	}
	FullPostInfo.Post = *post

	if wannaUser {
		user,err := pU.userRepo.ProfileInfo(post.Author)
		if err != nil {
			return nil,err
		}
		FullPostInfo.Author = *user
	}

	if wannaForum {
		forum,err := pU.forumRepo.ForumDetails(post.Forum)
		if err != nil {
			return nil,err
		}
		FullPostInfo.Forum = *forum
	}

	if wannaThread {
		thread,err := pU.threadRepo.ThreadDetailsByID(post.Thread)
		if err != nil {
			return nil,err
		}
		FullPostInfo.Thread = *thread
	}

	return &FullPostInfo,nil
}
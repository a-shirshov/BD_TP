package usecase

import (
	forumRepo "bd_tp/forum/repository"
	"bd_tp/models"
	postRepo "bd_tp/post/repository"
	threadRepo "bd_tp/thread/repository"
	userRepo "bd_tp/user/repository"
	"strconv"
	"strings"
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

func (pU *Usecase) PostDetails(idStr string,related string) (*models.FullPostInfo, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		
		return nil,err
	}

	var FullPostInfo models.FullPostInfo
	post, err := pU.postRepo.GetPostByID(id)
	if err != nil {
		
		return nil,err
	}
	FullPostInfo.Post = post

	if strings.Contains(related, "user") {
		user,err := pU.userRepo.ProfileInfo(post.Author)
		if err != nil {
			
			return nil,err
		}
		FullPostInfo.Author = user
	}

	if strings.Contains(related, "forum") {
		forum,err := pU.forumRepo.ForumDetails(post.Forum)
		if err != nil {
			
			return nil,err
		}
		FullPostInfo.Forum = forum
	}

	if strings.Contains(related, "thread") {
		thread,err := pU.threadRepo.ThreadDetailsByID(post.Thread)
		if err != nil {
			
			return nil,err
		}
		FullPostInfo.Thread = thread
	}

	return &FullPostInfo,nil
}

func (pU *Usecase) UpdatePost(post *models.Post) (*models.Post, error) {
	postFull, err := pU.PostDetails(strconv.Itoa(post.ID),"")
	if err != nil {
		return nil, err
	}
	if post.Message != postFull.Post.Message && post.Message != "" {
		post.Edited = true
	} else {
		return postFull.Post, nil
	}
	updatedPost,err := pU.postRepo.UpdatePost(post)
	if err != nil {
		return nil, err
	}
	return updatedPost, nil
}
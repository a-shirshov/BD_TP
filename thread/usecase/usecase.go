package usecase

import (
	"bd_tp/models"
	postRepo "bd_tp/post/repository"
	threadRepo "bd_tp/thread/repository"
	userRepo "bd_tp/user/repository"
	"errors"
	"fmt"
	"strconv"
)

type Usecase struct {
	threadRepo *threadRepo.Repository
	postRepo *postRepo.Repository
	userRepo *userRepo.Repository
}

func NewThreadUsecase (tR *threadRepo.Repository, pR *postRepo.Repository, uR *userRepo.Repository) *Usecase {
	return &Usecase{
		threadRepo: tR,
		postRepo: pR,
		userRepo: uR,
	}
}

func (tU *Usecase) CreatePosts(postsInput []models.Post, slug_or_id string) ([]models.Post, int, error) {
	id, err := strconv.Atoi(slug_or_id)

	if len(postsInput) == 0 {
		return postsInput,201, nil
	}
	var posts []models.Post

	if err != nil {
		thread, err := tU.threadRepo.ThreadDetailsBySlug(slug_or_id)
		if err != nil {
			return nil,404, err
		}
		fmt.Println(thread.Forum)
		for index := range postsInput {
			postsInput[index].Forum = thread.Forum
			if postsInput[index].Parent != 0 {
				parentPost, err := tU.postRepo.GetPostByID(postsInput[index].Parent)
				if err != nil {
					return nil, 404, err
				}
				threadSlugBase,err := tU.threadRepo.ThreadDetailsBySlug(slug_or_id)
				if err != nil {
					return nil, 404, err
				}
				if parentPost.Thread != threadSlugBase.ID {
					return nil,404, errors.New("thread and post mistake")
				}
			}
		}
		posts, code, err := tU.threadRepo.CreatePostsWithSlug(postsInput,slug_or_id)
		if err != nil {
			fmt.Println(err)
			return nil, code, err
		}
		return posts, code, nil
	}

	thread, err := tU.threadRepo.ThreadDetailsByID(id)
	if err != nil {
		return nil,404, err
	}
	fmt.Println(thread.Forum)

	for index := range postsInput {
		postsInput[index].Forum = thread.Forum
		if postsInput[index].Parent != 0 {
			parentPost, err := tU.postRepo.GetPostByID(postsInput[index].Parent)
			if err != nil {
				fmt.Println("Here err = ",err)
				return nil, 404, err
			}
			if parentPost.Thread != id {
				return nil,404, errors.New("thread and post mistake")
			}
		}
	}
	posts, code, err := tU.threadRepo.CreatePostsWithID(postsInput,id)
	if err != nil {
		fmt.Println(err)
		return nil, code, err
	}
	return posts, code, nil
}

func (tU *Usecase) ThreadDetails(slug_or_id string) (*models.Thread,error) {
	id, err := strconv.Atoi(slug_or_id) 
	if err != nil {
		thread, err := tU.threadRepo.ThreadDetailsBySlug(slug_or_id)
		if err != nil {
			return nil,err
		}
		return thread, nil
	}
	thread, err := tU.threadRepo.ThreadDetailsByID(id)
	if err != nil {
		return nil, err
	}
	return thread,nil
}

func (tU *Usecase) ThreadDetailsUpdate(threadInfo *models.Thread, slug_or_id string) (*models.Thread, error) {
	id, err := strconv.Atoi(slug_or_id) 
	if err != nil {
		thread, err := tU.threadRepo.ThreadDetailsUpdateBySlug(threadInfo,slug_or_id)
		if err != nil {
			return nil, err
		}
		return thread,nil
	}
	thread, err := tU.threadRepo.ThreadDetailsUpdateByID(threadInfo,id)
	if err != nil {
		return nil, err
	}
	return thread,nil
}

func (tU *Usecase) ThreadVote(vote *models.Vote, slug_or_id string) (*models.Thread, error) {
	userId,err := tU.userRepo.GetIdByNickname(vote.Nickname)
	if err != nil {
		fmt.Println(err)
	}
	id, err := strconv.Atoi(slug_or_id) 
	if err != nil {
		fmt.Println("Before ThreadVoteBySlug")
		fmt.Println(slug_or_id)
		/*
		thread, err := tU.threadRepo.ThreadVoteBySlug(vote,slug_or_id,userId)
		if err != nil {
			fmt.Println("err",err)
			return nil, err
		}
		return thread,nil
		*/
		threadInfo, err := tU.threadRepo.ThreadDetailsBySlug(slug_or_id)
		if err != nil {
			return nil,err
		}
		id = threadInfo.ID
	}
	thread, err := tU.threadRepo.ThreadVoteByID(vote,id,userId)
	if err != nil {
		return nil, err
	}
	return thread,nil
}

func (tU *Usecase) ThreadGetPosts(slug_or_id,limit,since,sort,desc string) ([]models.Post,error) {
	id, err := strconv.Atoi(slug_or_id) 
	var thread *models.Thread
	if err != nil {
		thread, err = tU.threadRepo.ThreadDetailsBySlug(slug_or_id)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = tU.threadRepo.ThreadDetailsByID(id)
		if err != nil {
			return nil, err
		}
	}

	if limit == "" {
		limit = "100"
	}
	posts, err := tU.threadRepo.ThreadGetPosts(thread.ID,limit,since,sort,desc)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
package response

import (
	models "bd_tp/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetUserFromRequest(r io.Reader) (*models.User, error) {
	userInput := new(UserResponse)
	//err := json.UnmarshalFromReader(r, userInput)
	err := json.NewDecoder(r).Decode(userInput)
	if err != nil {
		return nil, err
	}
	result := &models.User{
		Nickname: userInput.Nickname,
		Fullname: userInput.Fullname,
		About: userInput.About,
		Email: userInput.Email,
	}
	return result, nil
}

func GetForumFromRequest(r io.Reader) (*models.Forum, error) {
	forumInput := new(ForumResponse)
	//err := json.UnmarshalFromReader(r, forumInput)
	err := json.NewDecoder(r).Decode(forumInput)
	if err != nil {
		return nil, err
	}
	result := &models.Forum{
		Title: forumInput.Title,
		User: forumInput.User,
		Slug: forumInput.Slug,
	}
	return result, nil
}

func GetThreadFromRequest(r io.Reader) (*models.Thread, error) {
	threadInput := new(ThreadResponse)
	//err := json.UnmarshalFromReader(r, forumInput)
	err := json.NewDecoder(r).Decode(threadInput)
	if err != nil {
		return nil, err
	}
	result := &models.Thread{
		Title: threadInput.Title,
		Author: threadInput.Author,
		Message: threadInput.Message,
		Created: threadInput.Created,
	}
	return result, nil
}

func GetThreadsQueryInfo(r io.Reader) (*models.ForumThreadsRequest, error) {
	infoInput := new(ForumThreadsRequest)
	//err := json.UnmarshalFromReader(r, forumInput)
	err := json.NewDecoder(r).Decode(infoInput)
	if err != nil {
		return nil, err
	}
	result := &models.ForumThreadsRequest{
		Limit: infoInput.Limit,
		Since: infoInput.Since,
		Desc: infoInput.Desc,
	}
	return result, nil
}

func GetPostsFromRequest(r io.Reader) ([]models.Post, error) {
	var postsInput PostsRequest
	err := json.NewDecoder(r).Decode(&postsInput)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var posts []models.Post
	fmt.Println(posts)
	fmt.Println("I am posts:",postsInput)
	for _,post := range postsInput.Posts {
		posts = append(posts, models.Post{
			Parent: post.Parent,
			Author: post.Author,
			Message: post.Message,
		})
	}
	return posts,nil
}

func GetThreadUpdateFromRequest(r io.Reader) (*models.Thread,error) {
	var thread ThreadResponse
	err := json.NewDecoder(r).Decode(&thread) 
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := &models.Thread{
		Title: thread.Title,
		Message: thread.Message,
	}
	return result, nil
}

func GetVoteFromRequest(r io.Reader) (*models.Vote, error) {
	var vote VoteRequest
	err := json.NewDecoder(r).Decode(&vote) 
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := &models.Vote{
		Nickname: vote.Nickname,
		Voice: *vote.Voice,
	}
	return result, nil
}

func GetPostRelatedFromRequest(r io.Reader) (*models.PostsRelated, error) {
	var related PostRelated
	err := json.NewDecoder(r).Decode(&related) 
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := &models.PostsRelated{
		Related: related.Related,
	}
	return result, nil
}

func SendResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(response)
	if err != nil {
		return
	}
	_, err = w.Write(b)
	if err != nil {
		return
	}
}
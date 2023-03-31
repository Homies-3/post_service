package service

import (
	"log"
	"post_service/models"
	"post_service/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type postService struct {
	l  *log.Logger
	pR repositories.IPostRepository
}

func NewPostService(l *log.Logger, pR repositories.IPostRepository) postService {
	return postService{
		l:  l,
		pR: pR,
	}
}

func (pS postService) CreatePost(p models.PostRequest, userId string) error {

	uId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		pS.l.Println("Error parsing userid: ", err)
		return err
	}

	gId, err := primitive.ObjectIDFromHex(p.GroupId)
	if err != nil {
		pS.l.Println("Error parsing groupid : ", err)
		return err
	}

	var post models.Post

	post.Content = p.Content
	post.Title = p.Title
	post.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	post.Comments = make([]models.Comment, 0)
	post.UserId = uId
	post.Likes = 0
	post.GroupID = gId

	err = pS.pR.CreatePost(post)
	if err != nil {
		return err
	}

	//TODO send data to fan out service

	return nil
}

func (pS postService) DeletePost() error {

	return nil
}

package services

import (
	"context"
	"log"
	"net/http"
	"post_service/models"
	pb "post_service/pb"
	"post_service/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	l  *log.Logger
	pR repositories.IPostRepository
}

func NewService(l *log.Logger, pR repositories.IPostRepository) Service {
	return Service{
		l:  l,
		pR: pR,
	}
}

func (s Service) CreatePost(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {

	uId, err := primitive.ObjectIDFromHex(req.User.Username)
	if err != nil {
		s.l.Println("Error parsing userid: ", err)
		return &pb.PostResponse{
			Status: http.StatusBadRequest,
			Id:     "",
			Error:  "error parsing userid",
		}, err
	}

	gId, err := primitive.ObjectIDFromHex(req.GroupId)
	if err != nil {
		s.l.Println("Error parsing groupid : ", err)
		return &pb.PostResponse{
			Status: http.StatusBadRequest,
			Id:     "",
			Error:  "error parsing groupid",
		}, err
	}
	var post models.Post

	post.Content = req.Content
	post.Title = req.Title
	post.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	post.Comments = make([]primitive.ObjectID, 0)
	post.UserId = uId
	post.Likes = 0
	post.GroupID = gId

	id, err := s.pR.CreatePost(post)
	if err != nil {
		s.l.Println(err)
		return &pb.PostResponse{
			Status: http.StatusInternalServerError,
			Id:     "",
			Error:  "error creating post",
		}, err
	}

	// TODO send data to fan out service here

	return &pb.PostResponse{
		Status: http.StatusCreated,
		Id:     id.Hex(),
	}, nil
}

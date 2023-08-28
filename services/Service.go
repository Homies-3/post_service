package services

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"post_service/models"
	pb "post_service/pb"
	"post_service/repositories"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	l       *log.Logger
	pR      repositories.IPostRepository
	redis   *redis.Client
	channel *amqp.Channel
	cQ      *amqp.Queue
	uQ      *amqp.Queue

	pb.UnimplementedPostServiceServer
}

func NewService(l *log.Logger, pR repositories.IPostRepository, redis *redis.Client, cQ *amqp.Queue, uQ *amqp.Queue, channel *amqp.Channel) Service {
	return Service{
		l:       l,
		pR:      pR,
		redis:   redis,
		channel: channel,
		cQ:      cQ,
		uQ:      uQ,
	}
}

func (s Service) CreatePost(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	var post models.Post

	post.Content = req.Content
	post.Title = req.Title
	post.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	post.Comments = make([]primitive.ObjectID, 0)
	post.UserName = req.GetUser().GetUsername()
	post.Likes = 0

	id, err := s.pR.CreatePost(post)
	if err != nil {
		s.l.Println(err)
		return &pb.PostResponse{
			Status: http.StatusInternalServerError,
			Id:     "",
			Error:  "error creating post",
		}, err
	}
	var b bytes.Buffer
	if err = json.NewEncoder(&b).Encode(post); err != nil {
		s.l.Println(err)
	}

	s.l.Printf("post created with id %s\n", id)

	err = s.channel.PublishWithContext(ctx, "", s.cQ.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        b.Bytes(),
	})

	if err != nil {
		s.l.Println(err)
	}

	return &pb.PostResponse{
		Status: http.StatusCreated,
		Id:     id.Hex(),
	}, nil
}

package repositories

import (
	"context"
	"log"
	"post_service/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository struct {
	l  *log.Logger
	pC mongo.Collection
}

func NewPostRepository(l *log.Logger, pC mongo.Collection) IPostRepository {
	return PostRepository{
		l:  l,
		pC: pC,
	}
}

func (pR PostRepository) CreatePost(p models.Post) error {
	res, err := pR.pC.InsertOne(context.TODO(), p)
	pR.l.Println(res.InsertedID)

	if err != nil {
		pR.l.Println("Error creating post :", err)
		return err
	}

	return nil
}

func (pR PostRepository) DeletePost() error { return nil }

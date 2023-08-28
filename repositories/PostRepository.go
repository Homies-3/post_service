package repositories

import (
	"context"
	"log"
	"post_service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository struct {
	l  *log.Logger
	db *mongo.Database
}

func NewPostRepository(l *log.Logger, db *mongo.Database) IPostRepository {
	return PostRepository{
		l:  l,
		db: db,
	}
}

func (pR PostRepository) CreatePost(p models.Post) (primitive.ObjectID, error) {
	res, err := pR.db.Collection("posts").InsertOne(context.TODO(), p)
	pR.l.Println(res.InsertedID)

	if err != nil {
		pR.l.Println("Error creating post :", err)
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (pR PostRepository) DeletePost() error { return nil }

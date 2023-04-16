package repositories

import (
	"post_service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPostRepository interface {
	CreatePost(p models.Post) (primitive.ObjectID, error)
	DeletePost() error
}

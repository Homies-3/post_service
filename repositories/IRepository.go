package repositories

import "post_service/models"

type IPostRepository interface {
	CreatePost(p models.Post) error
	DeletePost() error
}

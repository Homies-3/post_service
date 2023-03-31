package service

import "post_service/models"

type IPostService interface {
	CreatePost(p models.PostRequest, userId string) error
	DeletePost() error
}

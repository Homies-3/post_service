package controllers

import "github.com/gin-gonic/gin"

type IController interface {
	CreatePost(c *gin.Context)
	DeletePost(c *gin.Context)
}

package controllers

import (
	"log"
	"net/http"
	"post_service/models"
	service "post_service/services"

	"github.com/gin-gonic/gin"
)

type postController struct {
	l  *log.Logger
	pS service.IPostService
}

func (pC postController) CreatePost(c *gin.Context) {
	var postRequest models.PostRequest

	if c.ShouldBind(&postRequest) == nil {
		userId := c.Query("u_id")
		err := pC.pS.CreatePost(postRequest, userId)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusCreated)
		return
	}

	c.Status(http.StatusBadRequest)
}

func (pC postController) DeletePost(c *gin.Context) {

}

func NewPostController(l *log.Logger, pS service.IPostService) IController {
	return postController{
		l:  l,
		pS: pS,
	}
}

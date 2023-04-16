package models

type PostRequest struct {
	Title   string   `form:"title"`
	Content string   `form:"content"`
	GroupId string   `json:"group_id"`
	Type    PostType `json:"type"`
}

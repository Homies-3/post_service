package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	CreatedAt         primitive.DateTime
	UserId            primitive.ObjectID
	UpdatedAt         primitive.DateTime
	Content           string
	Likes             int
	ReferedToUsername string
}

type Group struct {
	Name         string
	TotalMembers int
}

type PostType int

const (
	TEXT PostType = iota
)

type Post struct {
	CreatedAt primitive.DateTime `bson:"created_at"`
	UserId    primitive.ObjectID
	Likes     int
	Content   string
	Title     string
	GroupID   primitive.ObjectID
	Type      PostType
	Comments  []primitive.ObjectID
}

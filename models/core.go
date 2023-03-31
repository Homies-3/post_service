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

type Post struct {
	CreatedAt primitive.DateTime `bson:"created_at"`
	Comments  []Comment
	UserId    primitive.ObjectID
	Likes     int
	Content   string
	Title     string
	GroupID   primitive.ObjectID
}

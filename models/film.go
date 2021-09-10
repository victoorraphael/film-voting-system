package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Film struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Upvotes   int64              `bson:"upvotes"`
	Downvotes int64              `bson:"downvotes"`
	Score     int64              `bson:"score"`
}

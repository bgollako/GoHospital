package db

import "github.com/mongodb/mongo-go-driver/bson/primitive"

type Patient struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name"`
	Age     int                `json:"age"`
	Disease string             `json:"disease"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactMessage struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Contact   string             `bson:"contact" json:"contact"`
	Message   string             `bson:"message" json:"message"`
	IP        string             `bson:"ip" json:"ip"`
	Handled   bool               `bson:"handled" json:"handled"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebhookEvent struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	EventID   string             `bson:"eventId" json:"eventId"`
	Event     string             `bson:"event" json:"event"`
	RawBody   string             `bson:"rawBody" json:"rawBody"`
	Processed bool               `bson:"processed" json:"processed"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

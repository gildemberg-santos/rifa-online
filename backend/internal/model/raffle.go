package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RaffleStatus string

const (
	RaffleStatusActive    RaffleStatus = "ACTIVE"
	RaffleStatusCancelled RaffleStatus = "CANCELLED"
	RaffleStatusDrawn     RaffleStatus = "DRAWN"
)

type Raffle struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	OrganizerID  primitive.ObjectID `bson:"organizerId" json:"organizerId"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description" json:"description"`
	TicketPrice  int                `bson:"ticketPrice" json:"ticketPrice"`
	MaxNumbers   int                `bson:"maxNumbers" json:"maxNumbers"`
	DrawDate     time.Time          `bson:"drawDate" json:"drawDate"`
	ImageURL     string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Status       RaffleStatus       `bson:"status" json:"status"`
	ExternalID   string             `bson:"externalId,omitempty" json:"externalId,omitempty"`
	WinnerNumber *int               `bson:"winnerNumber,omitempty" json:"winnerNumber,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

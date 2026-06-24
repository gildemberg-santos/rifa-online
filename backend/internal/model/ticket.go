package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketStatus string

const (
	TicketStatusAvailable TicketStatus = "AVAILABLE"
	TicketStatusReserved  TicketStatus = "RESERVED"
	TicketStatusPaid      TicketStatus = "PAID"
)

type Ticket struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	RaffleID   primitive.ObjectID `bson:"raffleId" json:"raffleId"`
	Number     int                `bson:"number" json:"number"`
	Status     TicketStatus       `bson:"status" json:"status"`
	BuyerName  string             `bson:"buyerName,omitempty" json:"buyerName,omitempty"`
	BuyerEmail string             `bson:"buyerEmail,omitempty" json:"buyerEmail,omitempty"`
	PaymentID  string             `bson:"paymentId,omitempty" json:"paymentId,omitempty"`
	ReservedAt *time.Time         `bson:"reservedAt,omitempty" json:"reservedAt,omitempty"`
	PaidAt     *time.Time         `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}

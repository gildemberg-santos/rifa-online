package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
	PaymentStatusExpired  PaymentStatus = "EXPIRED"
)

type PaymentMethod string

const (
	PaymentMethodPIX    PaymentMethod = "PIX"
	PaymentMethodCard   PaymentMethod = "CARD"
	PaymentMethodBoleto PaymentMethod = "BOLETO"
)

type PaymentType string

const (
	PaymentTypeRaffle       PaymentType = "RAFFLE"
	PaymentTypeSubscription PaymentType = "SUBSCRIPTION"
)

type Payment struct {
	ID              primitive.ObjectID   `bson:"_id" json:"id"`
	Type            PaymentType          `bson:"type" json:"type"`
	RaffleID        primitive.ObjectID   `bson:"raffleId,omitempty" json:"raffleId,omitempty"`
	UserID          primitive.ObjectID   `bson:"userId,omitempty" json:"userId,omitempty"`
	TicketIDs       []primitive.ObjectID `bson:"ticketIds,omitempty" json:"ticketIds,omitempty"`
	BuyerName       string               `bson:"buyerName,omitempty" json:"buyerName,omitempty"`
	BuyerEmail      string               `bson:"buyerEmail,omitempty" json:"buyerEmail,omitempty"`
	BuyerPhone      string               `bson:"buyerPhone,omitempty" json:"buyerPhone,omitempty"`
	BuyerPhoneIndex string               `bson:"buyerPhoneIndex,omitempty" json:"-"`
	CheckoutURL     string               `bson:"checkoutUrl,omitempty" json:"checkoutUrl,omitempty"`
	InvoiceSlug     string               `bson:"invoiceSlug,omitempty" json:"invoiceSlug,omitempty"`
	TransactionNSU  string               `bson:"transactionNsu,omitempty" json:"transactionNsu,omitempty"`
	Amount          int                  `bson:"amount" json:"amount"`
	Status          PaymentStatus        `bson:"status" json:"status"`
	PaymentMethod   PaymentMethod        `bson:"paymentMethod,omitempty" json:"paymentMethod,omitempty"`
	PaidAt          *time.Time           `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	CreatedAt       time.Time            `bson:"createdAt" json:"createdAt"`
}

const SubscriptionPrice = 1000 // R$10,00 in cents

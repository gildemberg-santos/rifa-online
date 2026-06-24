package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusPaid      PaymentStatus = "PAID"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
	PaymentStatusExpired   PaymentStatus = "EXPIRED"
)

type PaymentMethod string

const (
	PaymentMethodPIX    PaymentMethod = "PIX"
	PaymentMethodCard   PaymentMethod = "CARD"
	PaymentMethodBoleto PaymentMethod = "BOLETO"
)

type Payment struct {
	ID             primitive.ObjectID   `bson:"_id" json:"id"`
	RaffleID       primitive.ObjectID   `bson:"raffleId" json:"raffleId"`
	TicketIDs      []primitive.ObjectID `bson:"ticketIds" json:"ticketIds"`
	BuyerName      string               `bson:"buyerName" json:"buyerName"`
	BuyerEmail     string               `bson:"buyerEmail" json:"buyerEmail"`
	CheckoutURL    string               `bson:"checkoutUrl,omitempty" json:"checkoutUrl,omitempty"`
	InvoiceSlug    string               `bson:"invoiceSlug,omitempty" json:"invoiceSlug,omitempty"`
	TransactionNSU string               `bson:"transactionNsu,omitempty" json:"transactionNsu,omitempty"`
	Amount         int                  `bson:"amount" json:"amount"`
	Status         PaymentStatus        `bson:"status" json:"status"`
	PaymentMethod  PaymentMethod        `bson:"paymentMethod,omitempty" json:"paymentMethod,omitempty"`
	PaidAt         *time.Time           `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	CreatedAt      time.Time            `bson:"createdAt" json:"createdAt"`
}

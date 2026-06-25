package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusInactive  SubscriptionStatus = "INACTIVE"
	SubscriptionStatusPastDue   SubscriptionStatus = "PAST_DUE"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type User struct {
	ID                   primitive.ObjectID  `bson:"_id" json:"id"`
	Name                 string              `bson:"name" json:"name"`
	Email                string              `bson:"email" json:"email"`
	PasswordHash         string              `bson:"passwordHash" json:"-"`
	Role                 Role                `bson:"role" json:"role"`
	InfinitePayHandle    string              `bson:"infinitePayHandle,omitempty" json:"infinitePayHandle,omitempty"`
	SubscriptionStatus   SubscriptionStatus   `bson:"subscriptionStatus" json:"subscriptionStatus"`
	SubscriptionExpiresAt *time.Time          `bson:"subscriptionExpiresAt,omitempty" json:"subscriptionExpiresAt,omitempty"`
	SubscriptionIsTrial  bool                `bson:"subscriptionIsTrial" json:"subscriptionIsTrial"`
	HasSubscriptionBefore bool               `bson:"hasSubscriptionBefore" json:"hasSubscriptionBefore"`
	CreatedAt            time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time           `bson:"updatedAt" json:"updatedAt"`
}

const TrialDays = 7

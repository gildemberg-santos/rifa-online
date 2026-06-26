package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/user/rifa-online/internal/model"
)

type ContactRepo struct {
	coll *mongo.Collection
}

func NewContactRepo(db *mongo.Database) *ContactRepo {
	return &ContactRepo{coll: db.Collection("contact_messages")}
}

func (r *ContactRepo) Insert(ctx context.Context, msg *model.ContactMessage) error {
	msg.ID = primitive.NewObjectID()
	msg.CreatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, msg)
	return err
}

func (r *ContactRepo) ListRecent(ctx context.Context, limit int64) ([]model.ContactMessage, error) {
	opts := options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit)
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	messages := []model.ContactMessage{}
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/user/rifa-online/internal/model"
)

type WebhookRepo struct {
	coll *mongo.Collection
}

func NewWebhookRepo(db *mongo.Database) *WebhookRepo {
	return &WebhookRepo{coll: db.Collection("webhook_events")}
}

func (r *WebhookRepo) Insert(ctx context.Context, event *model.WebhookEvent) error {
	event.ID = primitive.NewObjectID()
	event.CreatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, event)
	return err
}

func (r *WebhookRepo) FindByEventID(ctx context.Context, eventID string) (*model.WebhookEvent, error) {
	var event model.WebhookEvent
	err := r.coll.FindOne(ctx, bson.M{"eventId": eventID}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *WebhookRepo) MarkAsProcessed(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"processed": true},
	})
	return err
}

func (r *WebhookRepo) ExistsByEventID(ctx context.Context, eventID string) (bool, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{"eventId": eventID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

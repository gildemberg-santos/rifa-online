package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	register(Migration{
		Version:     1,
		Description: "create initial indexes",
		Up:          upCreateIndexes,
	})
}

func upCreateIndexes(ctx context.Context, db *mongo.Database) error {
	type idx struct {
		coll string
		keys bson.D
		opts *options.IndexOptions
	}

	indexes := []idx{
		{coll: "users", keys: bson.D{{Key: "email", Value: 1}}, opts: options.Index().SetUnique(true)},
		{coll: "raffles", keys: bson.D{{Key: "organizerId", Value: 1}, {Key: "status", Value: 1}}},
		{coll: "raffles", keys: bson.D{{Key: "status", Value: 1}, {Key: "drawDate", Value: 1}}},
		{coll: "tickets", keys: bson.D{{Key: "raffleId", Value: 1}, {Key: "number", Value: 1}}, opts: options.Index().SetUnique(true)},
		{coll: "tickets", keys: bson.D{{Key: "raffleId", Value: 1}, {Key: "status", Value: 1}}},
		{coll: "payments", keys: bson.D{{Key: "invoiceSlug", Value: 1}}, opts: options.Index().SetUnique(true).SetSparse(true)},
		{coll: "payments", keys: bson.D{{Key: "buyerPhone", Value: 1}}},
		{coll: "tickets", keys: bson.D{{Key: "buyerPhone", Value: 1}}},
		{coll: "webhook_events", keys: bson.D{{Key: "eventId", Value: 1}}, opts: options.Index().SetUnique(true)},
		{coll: "payments", keys: bson.D{{Key: "type", Value: 1}, {Key: "userId", Value: 1}}},
		{coll: "payments", keys: bson.D{{Key: "userId", Value: 1}}},
	}

	for _, i := range indexes {
		model := mongo.IndexModel{Keys: i.keys}
		if i.opts != nil {
			model.Options = i.opts
		}
		if _, err := db.Collection(i.coll).Indexes().CreateOne(ctx, model); err != nil {
			return err
		}
	}

	return nil
}

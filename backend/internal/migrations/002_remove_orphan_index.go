package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	register(Migration{
		Version:     2,
		Description: "drop orphan AbacatePay index",
		Up:          upDropOrphanIndex,
	})
}

func upDropOrphanIndex(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection("payments")

	cursor, err := coll.Indexes().List(ctx)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		var idx bson.M
		if err := cursor.Decode(&idx); err != nil {
			continue
		}
		name, _ := idx["name"].(string)
		if name == "abacateCheckoutId_1" {
			if _, err := coll.Indexes().DropOne(ctx, name); err != nil {
				return err
			}
		}
	}

	return nil
}

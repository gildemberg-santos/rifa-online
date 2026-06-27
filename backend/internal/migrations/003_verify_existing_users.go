package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	register(Migration{
		Version:     3,
		Description: "mark existing users as email verified",
		Up:          upVerifyExistingUsers,
	})
}

func upVerifyExistingUsers(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("users").UpdateMany(ctx,
		bson.M{"emailVerified": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"emailVerified": true}},
	)
	return err
}

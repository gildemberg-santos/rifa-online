package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/user/rifa-online/internal/crypto"
	"github.com/user/rifa-online/internal/model"
)

type ContactRepo struct {
	coll   *mongo.Collection
	cipher *crypto.Cipher
}

func NewContactRepo(db *mongo.Database, cipher *crypto.Cipher) *ContactRepo {
	return &ContactRepo{coll: db.Collection("contact_messages"), cipher: cipher}
}

func (r *ContactRepo) Insert(ctx context.Context, msg *model.ContactMessage) error {
	msg.ID = primitive.NewObjectID()
	msg.CreatedAt = time.Now()

	// Criptografa os campos com dados pessoais antes de persistir.
	name, err := r.cipher.Encrypt(msg.Name)
	if err != nil {
		return err
	}
	contact, err := r.cipher.Encrypt(msg.Contact)
	if err != nil {
		return err
	}
	message, err := r.cipher.Encrypt(msg.Message)
	if err != nil {
		return err
	}

	doc := *msg
	doc.Name = name
	doc.Contact = contact
	doc.Message = message

	_, err = r.coll.InsertOne(ctx, &doc)
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

	// Descriptografa para exibição no painel admin.
	for i := range messages {
		if messages[i].Name, err = r.cipher.Decrypt(messages[i].Name); err != nil {
			return nil, err
		}
		if messages[i].Contact, err = r.cipher.Decrypt(messages[i].Contact); err != nil {
			return nil, err
		}
		if messages[i].Message, err = r.cipher.Decrypt(messages[i].Message); err != nil {
			return nil, err
		}
	}

	return messages, nil
}

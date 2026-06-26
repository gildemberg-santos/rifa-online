package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/user/rifa-online/internal/crypto"
	"github.com/user/rifa-online/internal/model"
)

type UserRepo struct {
	coll   *mongo.Collection
	cipher *crypto.Cipher
}

func NewUserRepo(db *mongo.Database, cipher *crypto.Cipher) *UserRepo {
	return &UserRepo{coll: db.Collection("users"), cipher: cipher}
}

func (r *UserRepo) encrypt(u *model.User) error {
	var err error
	if u.Name != "" && !crypto.IsEncrypted(u.Name) {
		if u.Name, err = r.cipher.Encrypt(u.Name); err != nil {
			return err
		}
	}
	if u.Phone != "" && !crypto.IsEncrypted(u.Phone) {
		if u.Phone, err = r.cipher.Encrypt(u.Phone); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepo) decrypt(u *model.User) error {
	var err error
	if u.Name != "" && crypto.IsEncrypted(u.Name) {
		if u.Name, err = r.cipher.Decrypt(u.Name); err != nil {
			return err
		}
	}
	if u.Phone != "" && crypto.IsEncrypted(u.Phone) {
		if u.Phone, err = r.cipher.Decrypt(u.Phone); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepo) Insert(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if err := r.encrypt(user); err != nil {
		return err
	}
	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *UserRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	if err := r.decrypt(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	if err := r.decrypt(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	if err := r.encrypt(user); err != nil {
		return err
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (r *UserRepo) UpdateFields(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	updates["updatedAt"] = time.Now()
	if name, ok := updates["name"].(string); ok && name != "" {
		enc, err := r.cipher.Encrypt(name)
		if err != nil {
			return err
		}
		updates["name"] = enc
	}
	if phone, ok := updates["phone"].(string); ok && phone != "" {
		enc, err := r.cipher.Encrypt(phone)
		if err != nil {
			return err
		}
		updates["phone"] = enc
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *UserRepo) FindAll(ctx context.Context) ([]model.User, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	for i := range users {
		if err := r.decrypt(&users[i]); err != nil {
			return nil, err
		}
	}
	return users, nil
}

func (r *UserRepo) CountAll(ctx context.Context) (int, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *UserRepo) CountBySubscription(ctx context.Context, status string) (int, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{"subscriptionStatus": status})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *UserRepo) CountByTrial(ctx context.Context) (int, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{
		"subscriptionStatus":  model.SubscriptionStatusActive,
		"subscriptionIsTrial": true,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

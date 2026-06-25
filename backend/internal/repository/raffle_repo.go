package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/user/rifa-online/internal/model"
)

type RaffleRepo struct {
	coll *mongo.Collection
}

func NewRaffleRepo(db *mongo.Database) *RaffleRepo {
	return &RaffleRepo{coll: db.Collection("raffles")}
}

func (r *RaffleRepo) Insert(ctx context.Context, raffle *model.Raffle) error {
	raffle.ID = primitive.NewObjectID()
	raffle.CreatedAt = time.Now()
	raffle.UpdatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, raffle)
	return err
}

func (r *RaffleRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Raffle, error) {
	var raffle model.Raffle
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&raffle)
	if err != nil {
		return nil, err
	}
	return &raffle, nil
}

func (r *RaffleRepo) FindActive(ctx context.Context) ([]model.Raffle, error) {
	return r.findByStatus(ctx, model.RaffleStatusActive)
}

func (r *RaffleRepo) FindByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]model.Raffle, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"organizerId": organizerID})
	if err != nil {
		return nil, err
	}
	raffles := make([]model.Raffle, 0)
	if err := cursor.All(ctx, &raffles); err != nil {
		return nil, err
	}
	return raffles, nil
}

func (r *RaffleRepo) findByStatus(ctx context.Context, status model.RaffleStatus) ([]model.Raffle, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	raffles := make([]model.Raffle, 0)
	if err := cursor.All(ctx, &raffles); err != nil {
		return nil, err
	}
	return raffles, nil
}

func (r *RaffleRepo) Update(ctx context.Context, raffle *model.Raffle) error {
	raffle.UpdatedAt = time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": raffle.ID}, bson.M{"$set": raffle})
	return err
}

func (r *RaffleRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *RaffleRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.RaffleStatus) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	})
	return err
}

func (r *RaffleRepo) UpdateWinner(ctx context.Context, id primitive.ObjectID, winnerNumber int) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"winnerNumber": winnerNumber,
			"status":       model.RaffleStatusDrawn,
			"updatedAt":    time.Now(),
		},
	})
	return err
}

func (r *RaffleRepo) FindAll(ctx context.Context) ([]model.Raffle, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var raffles []model.Raffle
	if err := cursor.All(ctx, &raffles); err != nil {
		return nil, err
	}
	return raffles, nil
}

func (r *RaffleRepo) CountAll(ctx context.Context) (int, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *RaffleRepo) CountByStatus(ctx context.Context, status string) (int, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{"status": status})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

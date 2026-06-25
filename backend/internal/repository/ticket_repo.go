package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/user/rifa-online/internal/model"
)

type TicketRepo struct {
	coll *mongo.Collection
}

func NewTicketRepo(db *mongo.Database) *TicketRepo {
	return &TicketRepo{coll: db.Collection("tickets")}
}

func (r *TicketRepo) Insert(ctx context.Context, ticket *model.Ticket) error {
	ticket.ID = primitive.NewObjectID()
	ticket.CreatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, ticket)
	return err
}

func (r *TicketRepo) InsertMany(ctx context.Context, tickets []model.Ticket) error {
	now := time.Now()
	docs := make([]interface{}, len(tickets))
	for i := range tickets {
		tickets[i].ID = primitive.NewObjectID()
		tickets[i].CreatedAt = now
		docs[i] = tickets[i]
	}
	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

func (r *TicketRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepo) FindByRaffle(ctx context.Context, raffleID primitive.ObjectID) ([]model.Ticket, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"raffleId": raffleID})
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepo) FindByRaffleAndNumber(ctx context.Context, raffleID primitive.ObjectID, number int) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.coll.FindOne(ctx, bson.M{"raffleId": raffleID, "number": number}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepo) FindByRaffleAndStatus(ctx context.Context, raffleID primitive.ObjectID, status model.TicketStatus) ([]model.Ticket, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"raffleId": raffleID, "status": status})
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepo) Update(ctx context.Context, ticket *model.Ticket) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": ticket.ID}, bson.M{"$set": ticket})
	return err
}

func (r *TicketRepo) MarkAsPaid(ctx context.Context, ids []primitive.ObjectID, buyerName, buyerEmail, paymentID string) error {
	now := time.Now()
	_, err := r.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{
		"$set": bson.M{
			"status":     model.TicketStatusPaid,
			"buyerName":  buyerName,
			"buyerEmail": buyerEmail,
			"paymentId":  paymentID,
			"paidAt":     now,
		},
	})
	return err
}

func (r *TicketRepo) MarkAsReserved(ctx context.Context, ids []primitive.ObjectID) error {
	now := time.Now()
	_, err := r.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{
		"$set": bson.M{
			"status":     model.TicketStatusReserved,
			"reservedAt": now,
		},
	})
	return err
}

func (r *TicketRepo) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Ticket, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepo) FindPaidByRaffle(ctx context.Context, raffleID primitive.ObjectID) ([]model.Ticket, error) {
	return r.FindByRaffleAndStatus(ctx, raffleID, model.TicketStatusPaid)
}

func (r *TicketRepo) FindPaidByEmail(ctx context.Context, email string) ([]model.Ticket, error) {
	cursor, err := r.coll.Find(ctx, bson.M{
		"buyerEmail": email,
		"status":     model.TicketStatusPaid,
	})
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepo) CountByRaffleAndStatus(ctx context.Context, raffleID primitive.ObjectID, status model.TicketStatus) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{"raffleId": raffleID, "status": status})
}

func (r *TicketRepo) FindReservedOlderThan(ctx context.Context, cutoff time.Time) ([]model.Ticket, error) {
	cursor, err := r.coll.Find(ctx, bson.M{
		"status":     model.TicketStatusReserved,
		"reservedAt": bson.M{"$lt": cutoff},
	})
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepo) ReleaseReservations(ctx context.Context, ids []primitive.ObjectID) error {
	_, err := r.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{
		"$set":   bson.M{"status": model.TicketStatusAvailable},
		"$unset": bson.M{"reservedAt": "", "buyerName": "", "buyerEmail": "", "paymentId": ""},
	})
	return err
}

func (r *TicketRepo) DeleteByRaffle(ctx context.Context, raffleID primitive.ObjectID) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"raffleId": raffleID})
	return err
}

func (r *TicketRepo) UpdateManyStatus(ctx context.Context, ids []primitive.ObjectID, status model.TicketStatus) error {
	_, err := r.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{
		"$set": bson.M{"status": status},
	})
	return err
}

func (r *TicketRepo) CountAllPaid(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{"status": model.TicketStatusPaid})
}

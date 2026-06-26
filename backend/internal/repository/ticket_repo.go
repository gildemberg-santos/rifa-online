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

type TicketRepo struct {
	coll   *mongo.Collection
	cipher *crypto.Cipher
}

func NewTicketRepo(db *mongo.Database, cipher *crypto.Cipher) *TicketRepo {
	return &TicketRepo{coll: db.Collection("tickets"), cipher: cipher}
}

// encrypt cifra os campos sensíveis e calcula o índice cego do telefone (a partir
// do texto puro) antes de persistir.
func (r *TicketRepo) encrypt(t *model.Ticket) error {
	var err error
	t.BuyerPhoneIndex = r.cipher.BlindIndex(t.BuyerPhone)
	if t.BuyerName, err = r.cipher.Encrypt(t.BuyerName); err != nil {
		return err
	}
	if t.BuyerPhone, err = r.cipher.Encrypt(t.BuyerPhone); err != nil {
		return err
	}
	return nil
}

// decrypt reverte os campos e remove o índice cego do objeto retornado.
func (r *TicketRepo) decrypt(t *model.Ticket) error {
	var err error
	if t.BuyerName, err = r.cipher.Decrypt(t.BuyerName); err != nil {
		return err
	}
	if t.BuyerPhone, err = r.cipher.Decrypt(t.BuyerPhone); err != nil {
		return err
	}
	t.BuyerPhoneIndex = ""
	return nil
}

func (r *TicketRepo) decryptAll(ts []model.Ticket) error {
	for i := range ts {
		if err := r.decrypt(&ts[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *TicketRepo) Insert(ctx context.Context, ticket *model.Ticket) error {
	ticket.ID = primitive.NewObjectID()
	ticket.CreatedAt = time.Now()
	doc := *ticket
	if err := r.encrypt(&doc); err != nil {
		return err
	}
	_, err := r.coll.InsertOne(ctx, &doc)
	return err
}

func (r *TicketRepo) InsertMany(ctx context.Context, tickets []model.Ticket) error {
	now := time.Now()
	docs := make([]interface{}, len(tickets))
	for i := range tickets {
		tickets[i].ID = primitive.NewObjectID()
		tickets[i].CreatedAt = now
		doc := tickets[i]
		if err := r.encrypt(&doc); err != nil {
			return err
		}
		docs[i] = doc
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
	if err := r.decrypt(&ticket); err != nil {
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
	if err := r.decryptAll(tickets); err != nil {
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
	if err := r.decrypt(&ticket); err != nil {
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
	if err := r.decryptAll(tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepo) Update(ctx context.Context, ticket *model.Ticket) error {
	doc := *ticket
	if err := r.encrypt(&doc); err != nil {
		return err
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
	return err
}

func (r *TicketRepo) MarkAsPaid(ctx context.Context, ids []primitive.ObjectID, buyerName, buyerPhone, paymentID string) error {
	now := time.Now()
	encName, err := r.cipher.Encrypt(buyerName)
	if err != nil {
		return err
	}
	phoneIndex := r.cipher.BlindIndex(buyerPhone)
	encPhone, err := r.cipher.Encrypt(buyerPhone)
	if err != nil {
		return err
	}
	_, err = r.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{
		"$set": bson.M{
			"status":          model.TicketStatusPaid,
			"buyerName":       encName,
			"buyerPhone":      encPhone,
			"buyerPhoneIndex": phoneIndex,
			"paymentId":       paymentID,
			"paidAt":          now,
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
	if err := r.decryptAll(tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepo) FindPaidByRaffle(ctx context.Context, raffleID primitive.ObjectID) ([]model.Ticket, error) {
	return r.FindByRaffleAndStatus(ctx, raffleID, model.TicketStatusPaid)
}

func (r *TicketRepo) FindPaidByPhone(ctx context.Context, phone string) ([]model.Ticket, error) {
	cursor, err := r.coll.Find(ctx, bson.M{
		"buyerPhoneIndex": r.cipher.BlindIndex(phone),
		"status":          model.TicketStatusPaid,
	})
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	if err := r.decryptAll(tickets); err != nil {
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
		"$unset": bson.M{"reservedAt": "", "buyerName": "", "buyerPhone": "", "buyerPhoneIndex": "", "paymentId": ""},
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

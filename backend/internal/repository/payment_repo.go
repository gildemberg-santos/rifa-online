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

type PaymentRepo struct {
	coll   *mongo.Collection
	cipher *crypto.Cipher
}

func NewPaymentRepo(db *mongo.Database, cipher *crypto.Cipher) *PaymentRepo {
	return &PaymentRepo{coll: db.Collection("payments"), cipher: cipher}
}

// encrypt cifra os campos sensíveis e calcula o índice cego do telefone.
func (r *PaymentRepo) encrypt(p *model.Payment) error {
	var err error
	p.BuyerPhoneIndex = r.cipher.BlindIndex(p.BuyerPhone)
	if p.BuyerName, err = r.cipher.Encrypt(p.BuyerName); err != nil {
		return err
	}
	if p.BuyerEmail, err = r.cipher.Encrypt(p.BuyerEmail); err != nil {
		return err
	}
	if p.BuyerPhone, err = r.cipher.Encrypt(p.BuyerPhone); err != nil {
		return err
	}
	return nil
}

func (r *PaymentRepo) decrypt(p *model.Payment) error {
	var err error
	if p.BuyerName, err = r.cipher.Decrypt(p.BuyerName); err != nil {
		return err
	}
	if p.BuyerEmail, err = r.cipher.Decrypt(p.BuyerEmail); err != nil {
		return err
	}
	if p.BuyerPhone, err = r.cipher.Decrypt(p.BuyerPhone); err != nil {
		return err
	}
	p.BuyerPhoneIndex = ""
	return nil
}

func (r *PaymentRepo) decryptAll(ps []model.Payment) error {
	for i := range ps {
		if err := r.decrypt(&ps[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *PaymentRepo) Insert(ctx context.Context, payment *model.Payment) error {
	payment.ID = primitive.NewObjectID()
	payment.CreatedAt = time.Now()
	doc := *payment
	if err := r.encrypt(&doc); err != nil {
		return err
	}
	_, err := r.coll.InsertOne(ctx, &doc)
	return err
}

func (r *PaymentRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Payment, error) {
	var payment model.Payment
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	if err := r.decrypt(&payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepo) FindByInvoiceSlug(ctx context.Context, slug string) (*model.Payment, error) {
	var payment model.Payment
	err := r.coll.FindOne(ctx, bson.M{"invoiceSlug": slug}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	if err := r.decrypt(&payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepo) FindByOrderNSU(ctx context.Context, orderNSU string) (*model.Payment, error) {
	oid, err := primitive.ObjectIDFromHex(orderNSU)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, oid)
}

// FindByEmail busca por buyerEmail (campo agora criptografado, sem índice). Não
// é usado em nenhum fluxo ativo; mantido por compatibilidade.
func (r *PaymentRepo) FindByEmail(ctx context.Context, email string) ([]model.Payment, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"buyerEmail": email})
	if err != nil {
		return nil, err
	}
	payments := make([]model.Payment, 0)
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	if err := r.decryptAll(payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepo) FindByBuyerPhone(ctx context.Context, phone string) ([]model.Payment, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"buyerPhoneIndex": r.cipher.BlindIndex(phone)})
	if err != nil {
		return nil, err
	}
	payments := make([]model.Payment, 0)
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	if err := r.decryptAll(payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepo) FindByRaffle(ctx context.Context, raffleID primitive.ObjectID) ([]model.Payment, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"raffleId": raffleID})
	if err != nil {
		return nil, err
	}
	payments := make([]model.Payment, 0)
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	if err := r.decryptAll(payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepo) UpdateFields(ctx context.Context, id primitive.ObjectID, fields primitive.M) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields})
	return err
}

func (r *PaymentRepo) ExpirePendingOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.coll.UpdateMany(ctx, bson.M{
		"status":    model.PaymentStatusPending,
		"createdAt": bson.M{"$lt": cutoff},
	}, bson.M{"$set": bson.M{"status": model.PaymentStatusExpired}})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (r *PaymentRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.PaymentStatus, paidAt *time.Time) error {
	set := bson.M{
		"status": status,
	}
	if paidAt != nil {
		set["paidAt"] = paidAt
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": set})
	return err
}

func (r *PaymentRepo) Update(ctx context.Context, payment *model.Payment) error {
	doc := *payment
	if err := r.encrypt(&doc); err != nil {
		return err
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
	return err
}

func (r *PaymentRepo) DeleteByRaffle(ctx context.Context, raffleID primitive.ObjectID) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"raffleId": raffleID})
	return err
}

func (r *PaymentRepo) SumPaidByRaffle(ctx context.Context, raffleID primitive.ObjectID) (int64, error) {
	match := bson.M{"$match": bson.M{"raffleId": raffleID, "status": model.PaymentStatusPaid}}
	group := bson.M{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}
	cursor, err := r.coll.Aggregate(ctx, []bson.M{match, group})
	if err != nil {
		return 0, err
	}
	var results []struct {
		Total int64 `bson:"total"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Total, nil
}

func (r *PaymentRepo) SumAllPaid(ctx context.Context) (int64, error) {
	match := bson.M{"$match": bson.M{"status": model.PaymentStatusPaid}}
	group := bson.M{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}
	cursor, err := r.coll.Aggregate(ctx, []bson.M{match, group})
	if err != nil {
		return 0, err
	}
	var results []struct {
		Total int64 `bson:"total"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Total, nil
}

func (r *PaymentRepo) FindPendingSubscriptionByUserID(ctx context.Context, userID primitive.ObjectID) ([]model.Payment, error) {
	cursor, err := r.coll.Find(ctx, bson.M{
		"userId": userID,
		"type":   model.PaymentTypeSubscription,
		"status": model.PaymentStatusPending,
	})
	if err != nil {
		return nil, err
	}
	var payments []model.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	if err := r.decryptAll(payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]model.Payment, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	var payments []model.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	if err := r.decryptAll(payments); err != nil {
		return nil, err
	}
	return payments, nil
}

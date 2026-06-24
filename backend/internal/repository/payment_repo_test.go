package repository

import (
	"context"
	"testing"
	"time"

	"github.com/user/rifa-online/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPaymentRepo(t *testing.T) {
	ctx := context.Background()
	repo := NewPaymentRepo(testDB)

	if err := repo.Init(ctx); err != nil {
		t.Fatalf("Init: %v", err)
	}

	raffleID := primitive.NewObjectID()

	makePayment := func(email, slug string) *model.Payment {
		return &model.Payment{
			RaffleID:    raffleID,
			TicketIDs:   []primitive.ObjectID{primitive.NewObjectID()},
			BuyerName:   "Test Buyer",
			BuyerEmail:  email,
			InvoiceSlug: slug,
			Amount:      5000,
			Status:      model.PaymentStatusPending,
		}
	}

	t.Run("Insert and FindByID", func(t *testing.T) {
		p := makePayment("buyer1@example.com", "bill_001")
		if err := repo.Insert(ctx, p); err != nil {
			t.Fatalf("Insert: %v", err)
		}
		if p.ID == primitive.NilObjectID {
			t.Fatal("expected non-zero ID")
		}

		found, err := repo.FindByID(ctx, p.ID)
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if found.BuyerEmail != "buyer1@example.com" {
			t.Fatalf("expected buyer1@example.com, got %s", found.BuyerEmail)
		}
	})

	t.Run("FindByInvoiceSlug", func(t *testing.T) {
		repo.Insert(ctx, makePayment("buyer2@example.com", "slug_002"))

		found, err := repo.FindByInvoiceSlug(ctx, "slug_002")
		if err != nil {
			t.Fatalf("FindByInvoiceSlug: %v", err)
		}
		if found.InvoiceSlug != "slug_002" {
			t.Fatalf("expected slug_002, got %s", found.InvoiceSlug)
		}
	})

	t.Run("FindByEmail", func(t *testing.T) {
		repo.Insert(ctx, makePayment("buyer3@example.com", "bill_003"))
		repo.Insert(ctx, makePayment("buyer3@example.com", "bill_004"))

		payments, err := repo.FindByEmail(ctx, "buyer3@example.com")
		if err != nil {
			t.Fatalf("FindByEmail: %v", err)
		}
		if len(payments) < 2 {
			t.Fatalf("expected at least 2 payments, got %d", len(payments))
		}
	})

	t.Run("FindByRaffle", func(t *testing.T) {
		payments, err := repo.FindByRaffle(ctx, raffleID)
		if err != nil {
			t.Fatalf("FindByRaffle: %v", err)
		}
		if len(payments) < 4 {
			t.Fatalf("expected at least 4 payments, got %d", len(payments))
		}
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		p := makePayment("status@example.com", "bill_status")
		repo.Insert(ctx, p)

		now := time.Now()
		if err := repo.UpdateStatus(ctx, p.ID, model.PaymentStatusPaid, &now); err != nil {
			t.Fatalf("UpdateStatus: %v", err)
		}

		found, _ := repo.FindByID(ctx, p.ID)
		if found.Status != model.PaymentStatusPaid {
			t.Fatalf("expected PAID, got %s", found.Status)
		}
		if found.PaidAt == nil {
			t.Fatal("expected PaidAt to be set")
		}
	})

	t.Run("SumPaidByRaffle", func(t *testing.T) {
		p := makePayment("sum@example.com", "bill_sum")
		p.Amount = 10000
		p.Status = model.PaymentStatusPaid
		repo.Insert(ctx, p)

		total, err := repo.SumPaidByRaffle(ctx, raffleID)
		if err != nil {
			t.Fatalf("SumPaidByRaffle: %v", err)
		}
		if total < 10000 {
			t.Fatalf("expected total >= 10000, got %d", total)
		}
	})
}

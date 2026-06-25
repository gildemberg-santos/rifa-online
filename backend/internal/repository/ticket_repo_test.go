package repository

import (
	"context"
	"testing"

	"github.com/user/rifa-online/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTicketRepo(t *testing.T) {
	ctx := context.Background()
	repo := NewTicketRepo(testDB)
	raffleID := primitive.NewObjectID()

	makeTicket := func(number int) model.Ticket {
		return model.Ticket{
			RaffleID: raffleID,
			Number:   number,
			Status:   model.TicketStatusAvailable,
		}
	}

	t.Run("InsertMany and FindByRaffle", func(t *testing.T) {
		tickets := []model.Ticket{
			makeTicket(1),
			makeTicket(2),
			makeTicket(3),
		}
		if err := repo.InsertMany(ctx, tickets); err != nil {
			t.Fatalf("InsertMany: %v", err)
		}

		found, err := repo.FindByRaffle(ctx, raffleID)
		if err != nil {
			t.Fatalf("FindByRaffle: %v", err)
		}
		if len(found) != 3 {
			t.Fatalf("expected 3 tickets, got %d", len(found))
		}
	})

	t.Run("FindByRaffleAndNumber", func(t *testing.T) {
		ticket, err := repo.FindByRaffleAndNumber(ctx, raffleID, 1)
		if err != nil {
			t.Fatalf("FindByRaffleAndNumber: %v", err)
		}
		if ticket.Number != 1 {
			t.Fatalf("expected number 1, got %d", ticket.Number)
		}
	})

	t.Run("MarkAsPaid", func(t *testing.T) {
		tickets := []model.Ticket{makeTicket(10), makeTicket(11)}
		repo.InsertMany(ctx, tickets)

		all, _ := repo.FindByRaffle(ctx, raffleID)
		var ids []primitive.ObjectID
		for _, tk := range all {
			if tk.Number == 10 || tk.Number == 11 {
				ids = append(ids, tk.ID)
			}
		}
		if len(ids) != 2 {
			t.Fatalf("expected 2 ticket IDs, got %d", len(ids))
		}

		if err := repo.MarkAsPaid(ctx, ids, "Buyer", "11999999999", "pay_123"); err != nil {
			t.Fatalf("MarkAsPaid: %v", err)
		}

		paid, _ := repo.FindByRaffleAndStatus(ctx, raffleID, model.TicketStatusPaid)
		if len(paid) < 2 {
			t.Fatalf("expected at least 2 paid tickets, got %d", len(paid))
		}
	})

	t.Run("MarkAsReserved", func(t *testing.T) {
		tickets := []model.Ticket{makeTicket(20), makeTicket(21)}
		repo.InsertMany(ctx, tickets)

		all, _ := repo.FindByRaffle(ctx, raffleID)
		var ids []primitive.ObjectID
		for _, tk := range all {
			if tk.Number == 20 || tk.Number == 21 {
				ids = append(ids, tk.ID)
			}
		}

		if err := repo.MarkAsReserved(ctx, ids); err != nil {
			t.Fatalf("MarkAsReserved: %v", err)
		}

		reserved, _ := repo.FindByRaffleAndStatus(ctx, raffleID, model.TicketStatusReserved)
		if len(reserved) < 2 {
			t.Fatalf("expected at least 2 reserved tickets, got %d", len(reserved))
		}
	})

	t.Run("FindPaidByPhone", func(t *testing.T) {
		paid, err := repo.FindPaidByPhone(ctx, "11999999999")
		if err != nil {
			t.Fatalf("FindPaidByPhone: %v", err)
		}
		if len(paid) < 2 {
			t.Fatalf("expected at least 2 paid tickets for buyer, got %d", len(paid))
		}
		for _, tk := range paid {
			if tk.BuyerPhone != "11999999999" {
				t.Fatalf("expected buyerPhone 11999999999, got %s", tk.BuyerPhone)
			}
		}
	})

	t.Run("CountByRaffleAndStatus", func(t *testing.T) {
		count, err := repo.CountByRaffleAndStatus(ctx, raffleID, model.TicketStatusPaid)
		if err != nil {
			t.Fatalf("CountByRaffleAndStatus: %v", err)
		}
		if count < 2 {
			t.Fatalf("expected at least 2 paid tickets, got %d", count)
		}
	})

	t.Run("UpdateManyStatus", func(t *testing.T) {
		all, _ := repo.FindByRaffle(ctx, raffleID)
		if len(all) == 0 {
			t.Fatal("no tickets found")
		}
		var ids []primitive.ObjectID
		for _, tk := range all[:2] {
			ids = append(ids, tk.ID)
		}
		if err := repo.UpdateManyStatus(ctx, ids, model.TicketStatusAvailable); err != nil {
			t.Fatalf("UpdateManyStatus: %v", err)
		}
	})

	t.Run("FindByRaffleAndNumber not found", func(t *testing.T) {
		_, err := repo.FindByRaffleAndNumber(ctx, raffleID, 9999)
		if err == nil {
			t.Fatal("expected error for non-existent ticket")
		}
	})
}

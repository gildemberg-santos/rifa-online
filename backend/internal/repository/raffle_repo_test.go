package repository

import (
	"context"
	"testing"
	"time"

	"github.com/user/rifa-online/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRaffleRepo(t *testing.T) {
	ctx := context.Background()
	repo := NewRaffleRepo(testDB)

	if err := repo.Init(ctx); err != nil {
		t.Fatalf("Init: %v", err)
	}

	makeRaffle := func(title string, status model.RaffleStatus) *model.Raffle {
		return &model.Raffle{
			OrganizerID: primitive.NewObjectID(),
			Title:       title,
			Description: "A test raffle",
			TicketPrice: 500,
			MaxNumbers:  100,
			DrawDate:    time.Now().Add(7 * 24 * time.Hour),
			Status:      status,
		}
	}

	t.Run("Insert and FindByID", func(t *testing.T) {
		raffle := makeRaffle("Test Raffle", model.RaffleStatusActive)
		if err := repo.Insert(ctx, raffle); err != nil {
			t.Fatalf("Insert: %v", err)
		}
		if raffle.ID == primitive.NilObjectID {
			t.Fatal("expected non-zero ID")
		}

		found, err := repo.FindByID(ctx, raffle.ID)
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if found.Title != raffle.Title {
			t.Fatalf("expected '%s', got '%s'", raffle.Title, found.Title)
		}
	})

	t.Run("FindActive", func(t *testing.T) {
		repo.Insert(ctx, makeRaffle("Active 1", model.RaffleStatusActive))
		repo.Insert(ctx, makeRaffle("Active 2", model.RaffleStatusActive))
		repo.Insert(ctx, makeRaffle("Cancelled", model.RaffleStatusCancelled))

		active, err := repo.FindActive(ctx)
		if err != nil {
			t.Fatalf("FindActive: %v", err)
		}
		if len(active) < 2 {
			t.Fatalf("expected at least 2 active raffles, got %d", len(active))
		}
	})

	t.Run("FindByOrganizer", func(t *testing.T) {
		orgID := primitive.NewObjectID()
		repo.Insert(ctx, &model.Raffle{
			OrganizerID: orgID,
			Title:       "Org Raffle",
			TicketPrice: 1000,
			MaxNumbers:  50,
			DrawDate:    time.Now().Add(7 * 24 * time.Hour),
			Status:      model.RaffleStatusActive,
		})
		repo.Insert(ctx, &model.Raffle{
			OrganizerID: orgID,
			Title:       "Org Raffle 2",
			TicketPrice: 1000,
			MaxNumbers:  50,
			DrawDate:    time.Now().Add(7 * 24 * time.Hour),
			Status:      model.RaffleStatusActive,
		})

		raffles, err := repo.FindByOrganizer(ctx, orgID)
		if err != nil {
			t.Fatalf("FindByOrganizer: %v", err)
		}
		if len(raffles) != 2 {
			t.Fatalf("expected 2 raffles, got %d", len(raffles))
		}
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		raffle := makeRaffle("Status Test", model.RaffleStatusActive)
		repo.Insert(ctx, raffle)

		if err := repo.UpdateStatus(ctx, raffle.ID, model.RaffleStatusCancelled); err != nil {
			t.Fatalf("UpdateStatus: %v", err)
		}

		found, _ := repo.FindByID(ctx, raffle.ID)
		if found.Status != model.RaffleStatusCancelled {
			t.Fatalf("expected CANCELLED, got %s", found.Status)
		}
	})

	t.Run("UpdateWinner", func(t *testing.T) {
		raffle := makeRaffle("Winner Test", model.RaffleStatusActive)
		repo.Insert(ctx, raffle)

		if err := repo.UpdateWinner(ctx, raffle.ID, 42); err != nil {
			t.Fatalf("UpdateWinner: %v", err)
		}

		found, _ := repo.FindByID(ctx, raffle.ID)
		if found.Status != model.RaffleStatusDrawn {
			t.Fatalf("expected DRAWN, got %s", found.Status)
		}
		if found.WinnerNumber == nil || *found.WinnerNumber != 42 {
			t.Fatalf("expected winner 42, got %v", found.WinnerNumber)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		raffle := makeRaffle("Delete Me", model.RaffleStatusActive)
		repo.Insert(ctx, raffle)

		if err := repo.Delete(ctx, raffle.ID); err != nil {
			t.Fatalf("Delete: %v", err)
		}

		_, err := repo.FindByID(ctx, raffle.ID)
		if err != mongo.ErrNoDocuments {
			t.Fatalf("expected ErrNoDocuments, got: %v", err)
		}
	})
}

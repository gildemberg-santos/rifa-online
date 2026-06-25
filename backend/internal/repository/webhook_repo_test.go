package repository

import (
	"context"
	"testing"

	"github.com/user/rifa-online/internal/model"
)

func TestWebhookRepo(t *testing.T) {
	ctx := context.Background()
	repo := NewWebhookRepo(testDB)
	t.Run("Insert and FindByEventID", func(t *testing.T) {
		event := &model.WebhookEvent{
			EventID:   "evt_001",
			Event:     "checkout.completed",
			RawBody:   `{"id":"bill_001"}`,
			Processed: false,
		}
		if err := repo.Insert(ctx, event); err != nil {
			t.Fatalf("Insert: %v", err)
		}

		found, err := repo.FindByEventID(ctx, "evt_001")
		if err != nil {
			t.Fatalf("FindByEventID: %v", err)
		}
		if found.Event != "checkout.completed" {
			t.Fatalf("expected checkout.completed, got %s", found.Event)
		}
	})

	t.Run("ExistsByEventID", func(t *testing.T) {
		exists, err := repo.ExistsByEventID(ctx, "evt_001")
		if err != nil {
			t.Fatalf("ExistsByEventID: %v", err)
		}
		if !exists {
			t.Fatal("expected evt_001 to exist")
		}

		exists, err = repo.ExistsByEventID(ctx, "evt_nonexistent")
		if err != nil {
			t.Fatalf("ExistsByEventID: %v", err)
		}
		if exists {
			t.Fatal("expected evt_nonexistent to not exist")
		}
	})

	t.Run("MarkAsProcessed", func(t *testing.T) {
		event := &model.WebhookEvent{
			EventID:   "evt_002",
			Event:     "checkout.completed",
			RawBody:   `{}`,
			Processed: false,
		}
		repo.Insert(ctx, event)

		if err := repo.MarkAsProcessed(ctx, event.ID); err != nil {
			t.Fatalf("MarkAsProcessed: %v", err)
		}

		found, _ := repo.FindByEventID(ctx, "evt_002")
		if !found.Processed {
			t.Fatal("expected processed to be true")
		}
	})

	t.Run("Duplicate eventId returns error", func(t *testing.T) {
		dup := &model.WebhookEvent{
			EventID:   "evt_001",
			Event:     "checkout.completed",
			RawBody:   `{"duplicate":true}`,
			Processed: false,
		}
		if err := repo.Insert(ctx, dup); err == nil {
			t.Fatal("expected error for duplicate eventId")
		}
	})
}

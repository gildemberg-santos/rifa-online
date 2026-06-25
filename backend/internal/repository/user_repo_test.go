package repository

import (
	"context"
	"testing"

	"github.com/user/rifa-online/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUserRepo(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testDB)
	t.Run("Insert and FindByEmail", func(t *testing.T) {
		user := &model.User{
			Name:         "John Doe",
			Email:        "john@example.com",
			PasswordHash: "hashedpassword",
		}
		if err := repo.Insert(ctx, user); err != nil {
			t.Fatalf("Insert: %v", err)
		}
		if user.ID == primitive.NilObjectID {
			t.Fatal("expected non-zero ID")
		}

		found, err := repo.FindByEmail(ctx, "john@example.com")
		if err != nil {
			t.Fatalf("FindByEmail: %v", err)
		}
		if found.Email != user.Email {
			t.Fatalf("expected %s, got %s", user.Email, found.Email)
		}
	})

	t.Run("Duplicate email returns error", func(t *testing.T) {
		user := &model.User{
			Name:         "Jane Doe",
			Email:        "john@example.com",
			PasswordHash: "anotherhash",
		}
		err := repo.Insert(ctx, user)
		if !mongo.IsDuplicateKeyError(err) {
			t.Fatalf("expected duplicate key error, got: %v", err)
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		user := &model.User{
			Name:         "Find Me",
			Email:        "findme@example.com",
			PasswordHash: "hash",
		}
		repo.Insert(ctx, user)

		found, err := repo.FindByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if found.Name != "Find Me" {
			t.Fatalf("expected 'Find Me', got '%s'", found.Name)
		}
	})

	t.Run("FindByEmail not found", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, "nonexistent@example.com")
		if err != mongo.ErrNoDocuments {
			t.Fatalf("expected ErrNoDocuments, got: %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		user := &model.User{
			Name:         "Update Me",
			Email:        "update@example.com",
			PasswordHash: "oldhash",
		}
		repo.Insert(ctx, user)

		user.Name = "Updated Name"
		if err := repo.Update(ctx, user); err != nil {
			t.Fatalf("Update: %v", err)
		}

		found, _ := repo.FindByID(ctx, user.ID)
		if found.Name != "Updated Name" {
			t.Fatalf("expected 'Updated Name', got '%s'", found.Name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		user := &model.User{
			Name:         "Delete Me",
			Email:        "delete@example.com",
			PasswordHash: "hash",
		}
		repo.Insert(ctx, user)

		if err := repo.Delete(ctx, user.ID); err != nil {
			t.Fatalf("Delete: %v", err)
		}

		_, err := repo.FindByID(ctx, user.ID)
		if err != mongo.ErrNoDocuments {
			t.Fatalf("expected ErrNoDocuments after delete, got: %v", err)
		}
	})
}

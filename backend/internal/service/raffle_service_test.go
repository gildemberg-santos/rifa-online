package service

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewRaffleService(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)
	if svc == nil {
		t.Error("NewRaffleService returned nil")
	}
}

func TestRaffleService_CreateValidation(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	tests := []struct {
		name    string
		input   CreateRaffleInput
		wantErr bool
	}{
		{"empty organizer id", CreateRaffleInput{OrganizerID: ""}, true},
		{"invalid organizer id", CreateRaffleInput{OrganizerID: "bad-hex"}, true},
		{"short organizer id", CreateRaffleInput{OrganizerID: "abc123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestRaffleService_UpdateNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Update panics with nil repos as expected")
		}
	}()

	svc.Update(context.Background(), primitive.NewObjectID(), primitive.NewObjectID(), CreateRaffleInput{})
	t.Error("Update should have panicked with nil repo")
}

func TestRaffleService_DeleteNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Delete panics with nil repos as expected")
		}
	}()

	svc.Delete(context.Background(), primitive.NewObjectID(), primitive.NewObjectID())
	t.Error("Delete should have panicked with nil repo")
}

func TestRaffleService_CancelNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Cancel panics with nil repos as expected")
		}
	}()

	svc.Cancel(context.Background(), primitive.NewObjectID(), primitive.NewObjectID())
	t.Error("Cancel should have panicked with nil repo")
}

func TestRaffleService_DrawNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Draw panics with nil repos as expected")
		}
	}()

	svc.Draw(context.Background(), primitive.NewObjectID(), primitive.NewObjectID())
	t.Error("Draw should have panicked with nil repo")
}

package service

import (
	"context"
	"testing"
	"time"

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

func TestRaffleService_GetDetailNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("GetDetail panics with nil repo as expected")
		}
	}()

	svc.GetDetail(context.Background(), primitive.NewObjectID(), "")
	t.Error("GetDetail should have panicked with nil repo")
}

func TestRaffleService_GetMyRafflesNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("GetMyRaffles panics with nil repo as expected")
		}
	}()

	svc.GetMyRaffles(context.Background(), primitive.NewObjectID())
	t.Error("GetMyRaffles should have panicked with nil repo")
}

func TestRaffleService_ListActiveNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("ListActive panics with nil repo as expected")
		}
	}()

	svc.ListActive(context.Background())
	t.Error("ListActive should have panicked with nil repo")
}

func TestRaffleService_GetStatsNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("GetStats panics with nil repo as expected")
		}
	}()

	svc.GetStats(context.Background(), primitive.NewObjectID(), primitive.NewObjectID(), true)
	t.Error("GetStats should have panicked with nil repo")
}

func TestRaffleService_UpdateNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Update panics with nil repos as expected")
		}
	}()

	svc.Update(context.Background(), primitive.NewObjectID(), primitive.NewObjectID(), CreateRaffleInput{}, false)
	t.Error("Update should have panicked with nil repo")
}

func TestRaffleService_DeleteNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Delete panics with nil repos as expected")
		}
	}()

	svc.Delete(context.Background(), primitive.NewObjectID(), primitive.NewObjectID(), false)
	t.Error("Delete should have panicked with nil repo")
}

func TestRaffleService_CancelNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Cancel panics with nil repos as expected")
		}
	}()

	svc.Cancel(context.Background(), primitive.NewObjectID(), primitive.NewObjectID(), false)
	t.Error("Cancel should have panicked with nil repo")
}

func TestRaffleService_DrawNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Draw panics with nil repos as expected")
		}
	}()

	svc.Draw(context.Background(), primitive.NewObjectID(), primitive.NewObjectID(), false)
	t.Error("Draw should have panicked with nil repo")
}

func TestRaffleService_GetDashboardStatsNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("GetDashboardStats panics with nil repo as expected")
		}
	}()

	svc.GetDashboardStats(context.Background(), primitive.NewObjectID())
	t.Error("GetDashboardStats should have panicked with nil repo")
}

func TestRaffleService_CreateNilRepo(t *testing.T) {
	svc := NewRaffleService(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Create panics with nil repo as expected for valid organizer ID")
		}
	}()

	svc.Create(context.Background(), CreateRaffleInput{
		OrganizerID: primitive.NewObjectID().Hex(),
		Title:       "Test",
		TicketPrice: 10,
		MaxNumbers:  100,
		DrawDate:    time.Now().Add(24 * time.Hour),
	})
	t.Error("Create should have panicked with nil repo")
}

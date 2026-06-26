package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/user/rifa-online/internal/crypto"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27018/rifaonline"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("rifaonline")

	cipher, err := crypto.New(
		"22e32ebabf73e7516acb50a35df807285da5f5e740d3ab5635edf25d058912f6",
		"f71f66102fc7a05c224e638aa4f37086a36307d9a4ae90e985406c7601e93583",
	)
	if err != nil {
		log.Fatal(err)
	}

	ticketRepo := repository.NewTicketRepo(db, cipher)
	raffleRepo := repository.NewRaffleRepo(db)

	organizerID, _ := primitive.ObjectIDFromHex("6a3d5669bf0caadb0eecf351")

	raffle := &model.Raffle{
		OrganizerID: organizerID,
		Title:       "Sorteio Teste - 5/10 Vendidos",
		Description: "Rifa com 10 números, apenas 5 vendidos. Testar sorteio apenas entre pagos.",
		TicketPrice: 500,
		MaxNumbers:  10,
		DrawDate:    time.Now().Add(7 * 24 * time.Hour),
		Status:      model.RaffleStatusActive,
	}

	if err := raffleRepo.Insert(ctx, raffle); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Raffle created: %s (%s)\n", raffle.Title, raffle.ID.Hex())

	for i := 1; i <= raffle.MaxNumbers; i++ {
		ticket := &model.Ticket{
			RaffleID:  raffle.ID,
			Number:    i,
			Status:    model.TicketStatusAvailable,
			CreatedAt: time.Now(),
		}
		if err := ticketRepo.Insert(ctx, ticket); err != nil {
			log.Fatal(err)
		}
	}

	mockBuyers := []struct {
		name  string
		phone string
	}{
		{"Ana Beatriz Silva", "11987654321"},
		{"Carlos Eduardo Oliveira", "21976543210"},
		{"Daniela Ferreira Santos", "31965432109"},
		{"Eduardo Gomes Lima", "41954321098"},
		{"Fernanda Costa Rocha", "51943210987"},
	}

	allTickets, err := ticketRepo.FindByRaffle(ctx, raffle.ID)
	if err != nil {
		log.Fatal(err)
	}

	for i, t := range allTickets {
		if i >= 5 {
			break
		}
		buyer := mockBuyers[i]
		now := time.Now()

		encName, err := cipher.Encrypt(buyer.name)
		if err != nil {
			log.Fatal(err)
		}
		encPhone, err := cipher.Encrypt(buyer.phone)
		if err != nil {
			log.Fatal(err)
		}
		phoneIndex := cipher.BlindIndex(buyer.phone)
		paymentID := fmt.Sprintf("mock_pay_%s_%d", raffle.ID.Hex()[:8], t.Number)

		_, err = db.Collection("tickets").UpdateOne(ctx, bson.M{"_id": t.ID}, bson.M{
			"$set": bson.M{
				"status":          model.TicketStatusPaid,
				"buyerName":       encName,
				"buyerPhone":      encPhone,
				"buyerPhoneIndex": phoneIndex,
				"paymentId":       paymentID,
				"paidAt":          now,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Ticket %2d -> %s (%s) [pay: %s]\n", t.Number, buyer.name, buyer.phone, paymentID)
	}

	fmt.Printf("\n  Tickets 6-10 kept as AVAILABLE (not in draw pool)\n")
	fmt.Println("\nDone! Raffle with 5 paid + 5 available tickets created.")
}

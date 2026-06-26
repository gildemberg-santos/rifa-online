package service

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
)

var (
	ErrRaffleNotFound     = errors.New("raffle not found")
	ErrNotRaffleOwner     = errors.New("not the raffle owner")
	ErrRaffleAlreadyDrawn = errors.New("raffle already drawn")
	ErrRaffleCancelled    = errors.New("raffle is cancelled")
)

type RaffleService struct {
	raffleRepo  *repository.RaffleRepo
	ticketRepo  *repository.TicketRepo
	paymentRepo *repository.PaymentRepo
	userRepo    *repository.UserRepo
}

func NewRaffleService(raffleRepo *repository.RaffleRepo, ticketRepo *repository.TicketRepo, paymentRepo *repository.PaymentRepo, userRepo *repository.UserRepo) *RaffleService {
	return &RaffleService{
		raffleRepo:  raffleRepo,
		ticketRepo:  ticketRepo,
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
	}
}

type CreateRaffleInput struct {
	OrganizerID string
	Title       string
	Description string
	TicketPrice int
	MaxNumbers  int
	DrawDate    time.Time
	ImageURL    string
}

func (s *RaffleService) Create(ctx context.Context, input CreateRaffleInput) (*model.Raffle, error) {
	oid, err := primitive.ObjectIDFromHex(input.OrganizerID)
	if err != nil {
		return nil, errors.New("invalid organizer id")
	}

	user, err := s.userRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.SubscriptionStatus != model.SubscriptionStatusActive {
		return nil, ErrSubscriptionNotActive
	}

	if user.SubscriptionExpiresAt != nil && time.Now().After(*user.SubscriptionExpiresAt) {
		s.userRepo.UpdateFields(ctx, oid, primitive.M{
			"subscriptionStatus": model.SubscriptionStatusPastDue,
		})
		return nil, ErrSubscriptionNotActive
	}

	if input.Title == "" || len(input.Title) > 100 {
		return nil, errors.New("title must be between 1 and 100 characters")
	}
	if len(input.Description) > 500 {
		return nil, errors.New("description must be at most 500 characters")
	}
	if input.TicketPrice <= 0 {
		return nil, errors.New("ticket price must be positive")
	}
	if input.MaxNumbers < 1 || input.MaxNumbers > 1000 {
		return nil, errors.New("number of tickets must be between 1 and 1000")
	}
	if input.DrawDate.Before(time.Now()) {
		return nil, errors.New("draw date must be in the future")
	}

	raffle := &model.Raffle{
		OrganizerID: oid,
		Title:       input.Title,
		Description: input.Description,
		TicketPrice: input.TicketPrice,
		MaxNumbers:  input.MaxNumbers,
		DrawDate:    input.DrawDate,
		ImageURL:    input.ImageURL,
		Status:      model.RaffleStatusActive,
	}

	if err := s.raffleRepo.Insert(ctx, raffle); err != nil {
		return nil, err
	}

	tickets := make([]model.Ticket, input.MaxNumbers)
	for i := 0; i < input.MaxNumbers; i++ {
		tickets[i] = model.Ticket{
			RaffleID: raffle.ID,
			Number:   i + 1,
			Status:   model.TicketStatusAvailable,
		}
	}

	if err := s.ticketRepo.InsertMany(ctx, tickets); err != nil {
		return nil, err
	}

	return raffle, nil
}

func (s *RaffleService) ListActive(ctx context.Context) ([]model.Raffle, error) {
	return s.raffleRepo.FindActive(ctx)
}

type RaffleDetail struct {
	Raffle  *model.Raffle  `json:"raffle"`
	Tickets []model.Ticket `json:"tickets"`
}

const reservationTTL = 10 * time.Minute

func (s *RaffleService) GetDetail(ctx context.Context, raffleID primitive.ObjectID, userID string) (*RaffleDetail, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return nil, ErrRaffleNotFound
	}

	tickets, err := s.ticketRepo.FindByRaffle(ctx, raffleID)
	if err != nil {
		return nil, err
	}

	isOwner := userID != "" && raffle.OrganizerID.Hex() == userID

	now := time.Now()
	for i := range tickets {
		if tickets[i].Status == model.TicketStatusReserved && tickets[i].ReservedAt != nil {
			elapsed := now.Sub(*tickets[i].ReservedAt)
			remaining := reservationTTL - elapsed
			if remaining > 0 {
				secs := int(remaining.Seconds())
				tickets[i].ReservationExpiresIn = &secs
			} else {
				zero := 0
				tickets[i].ReservationExpiresIn = &zero
			}
		}

		if !isOwner {
			tickets[i].BuyerName = ""
			tickets[i].BuyerPhone = ""
			tickets[i].PaymentID = ""
		}
	}

	return &RaffleDetail{Raffle: raffle, Tickets: tickets}, nil
}

func (s *RaffleService) GetMyRaffles(ctx context.Context, organizerID primitive.ObjectID) ([]model.Raffle, error) {
	return s.raffleRepo.FindByOrganizer(ctx, organizerID)
}

func (s *RaffleService) Update(ctx context.Context, raffleID primitive.ObjectID, organizerID primitive.ObjectID, input CreateRaffleInput, isAdmin bool) (*model.Raffle, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return nil, ErrRaffleNotFound
	}
	if raffle.OrganizerID != organizerID && !isAdmin {
		return nil, ErrNotRaffleOwner
	}
	if raffle.Status != model.RaffleStatusActive {
		return nil, errors.New("can only edit active raffles")
	}

	if input.Title == "" || len(input.Title) > 100 {
		return nil, errors.New("title must be between 1 and 100 characters")
	}
	if len(input.Description) > 500 {
		return nil, errors.New("description must be at most 500 characters")
	}
	if input.TicketPrice <= 0 {
		return nil, errors.New("ticket price must be positive")
	}
	if input.DrawDate.Before(time.Now()) {
		return nil, errors.New("draw date must be in the future")
	}

	raffle.Title = input.Title
	raffle.Description = input.Description
	raffle.TicketPrice = input.TicketPrice
	raffle.DrawDate = input.DrawDate
	raffle.ImageURL = input.ImageURL

	if err := s.raffleRepo.Update(ctx, raffle); err != nil {
		return nil, err
	}
	return raffle, nil
}

func (s *RaffleService) Delete(ctx context.Context, raffleID primitive.ObjectID, organizerID primitive.ObjectID, isAdmin bool) error {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return ErrRaffleNotFound
	}
	if raffle.OrganizerID != organizerID && !isAdmin {
		return ErrNotRaffleOwner
	}

	if err := s.ticketRepo.DeleteByRaffle(ctx, raffleID); err != nil {
		return err
	}
	if err := s.paymentRepo.DeleteByRaffle(ctx, raffleID); err != nil {
		return err
	}
	return s.raffleRepo.Delete(ctx, raffleID)
}

func (s *RaffleService) Cancel(ctx context.Context, raffleID primitive.ObjectID, organizerID primitive.ObjectID, isAdmin bool) error {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return ErrRaffleNotFound
	}
	if raffle.OrganizerID != organizerID && !isAdmin {
		return ErrNotRaffleOwner
	}
	if raffle.Status == model.RaffleStatusDrawn {
		return ErrRaffleAlreadyDrawn
	}

	return s.raffleRepo.UpdateStatus(ctx, raffleID, model.RaffleStatusCancelled)
}

type RaffleStats struct {
	TotalSold       int     `json:"totalSold"`
	TotalRevenue    int64   `json:"totalRevenue"`
	PercentageSold  float64 `json:"percentageSold"`
	TicketPrice     int     `json:"ticketPrice"`
	MaxNumbers      int     `json:"maxNumbers"`
}

func (s *RaffleService) GetStats(ctx context.Context, raffleID primitive.ObjectID, organizerID primitive.ObjectID, isAdmin bool) (*RaffleStats, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return nil, ErrRaffleNotFound
	}
	if raffle.OrganizerID != organizerID && !isAdmin {
		return nil, ErrNotRaffleOwner
	}

	soldCount, err := s.ticketRepo.CountByRaffleAndStatus(ctx, raffleID, model.TicketStatusPaid)
	if err != nil {
		return nil, err
	}

	totalRevenue, err := s.paymentRepo.SumPaidByRaffle(ctx, raffleID)
	if err != nil {
		return nil, err
	}

	percentage := 0.0
	if raffle.MaxNumbers > 0 {
		percentage = math.Round(float64(soldCount)/float64(raffle.MaxNumbers)*10000) / 100
	}

	return &RaffleStats{
		TotalSold:      int(soldCount),
		TotalRevenue:   totalRevenue,
		PercentageSold: percentage,
		TicketPrice:    raffle.TicketPrice,
		MaxNumbers:     raffle.MaxNumbers,
	}, nil
}

type DashboardStats struct {
	TotalRaffles         int     `json:"totalRaffles"`
	ActiveRaffles        int     `json:"activeRaffles"`
	DrawnRaffles         int     `json:"drawnRaffles"`
	CancelledRaffles     int     `json:"cancelledRaffles"`
	TotalSoldTickets     int     `json:"totalSoldTickets"`
	TotalRevenue         int64   `json:"totalRevenue"`
	TotalReservedTickets int     `json:"totalReservedTickets"`
	TotalMaxNumbers      int     `json:"totalMaxNumbers"`
	TotalAvailableTickets int    `json:"totalAvailableTickets"`
}

func (s *RaffleService) GetDashboardStats(ctx context.Context, organizerID primitive.ObjectID) (*DashboardStats, error) {
	raffles, err := s.raffleRepo.FindByOrganizer(ctx, organizerID)
	if err != nil {
		return nil, err
	}

	stats := &DashboardStats{
		TotalRaffles: len(raffles),
	}

	for _, r := range raffles {
		switch r.Status {
		case model.RaffleStatusActive:
			stats.ActiveRaffles++
		case model.RaffleStatusDrawn:
			stats.DrawnRaffles++
		case model.RaffleStatusCancelled:
			stats.CancelledRaffles++
		}

		stats.TotalMaxNumbers += r.MaxNumbers

		soldCount, err := s.ticketRepo.CountByRaffleAndStatus(ctx, r.ID, model.TicketStatusPaid)
		if err != nil {
			return nil, err
		}
		stats.TotalSoldTickets += int(soldCount)

		reservedCount, err := s.ticketRepo.CountByRaffleAndStatus(ctx, r.ID, model.TicketStatusReserved)
		if err != nil {
			return nil, err
		}
		stats.TotalReservedTickets += int(reservedCount)

		revenue, err := s.paymentRepo.SumPaidByRaffle(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		stats.TotalRevenue += revenue
	}

	stats.TotalAvailableTickets = stats.TotalMaxNumbers - stats.TotalSoldTickets - stats.TotalReservedTickets
	if stats.TotalAvailableTickets < 0 {
		stats.TotalAvailableTickets = 0
	}

	return stats, nil
}

type DrawResult struct {
	WinnerNumber int            `json:"winnerNumber"`
	Raffle       *model.Raffle `json:"raffle"`
}

func (s *RaffleService) Draw(ctx context.Context, raffleID primitive.ObjectID, organizerID primitive.ObjectID, isAdmin bool) (*DrawResult, error) {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return nil, ErrRaffleNotFound
	}
	if raffle.OrganizerID != organizerID && !isAdmin {
		return nil, ErrNotRaffleOwner
	}
	if raffle.Status != model.RaffleStatusActive {
		return nil, errors.New("raffle is not active")
	}

	paidTickets, err := s.ticketRepo.FindPaidByRaffle(ctx, raffleID)
	if err != nil {
		return nil, err
	}
	if len(paidTickets) == 0 {
		return nil, errors.New("no paid tickets to draw from")
	}

	winner := paidTickets[time.Now().UnixNano()%int64(len(paidTickets))]

	if err := s.raffleRepo.UpdateWinner(ctx, raffleID, winner.Number); err != nil {
		return nil, err
	}

	raffle.Status = model.RaffleStatusDrawn
	raffle.WinnerNumber = &winner.Number

	return &DrawResult{
		WinnerNumber: winner.Number,
		Raffle:       raffle,
	}, nil
}

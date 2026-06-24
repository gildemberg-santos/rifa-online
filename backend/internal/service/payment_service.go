package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/pkg/infinitepay"
)

var (
	ErrNumbersUnavailable = errors.New("one or more numbers are unavailable")
	ErrRaffleNotActive    = errors.New("raffle is not active")
)

type PaymentService struct {
	raffleRepo      *repository.RaffleRepo
	ticketRepo      *repository.TicketRepo
	paymentRepo     *repository.PaymentRepo
	infiniteClient  *infinitepay.Client
	redisClient     *redis.Client
	cfg             *config.Config
}

func NewPaymentService(
	raffleRepo *repository.RaffleRepo,
	ticketRepo *repository.TicketRepo,
	paymentRepo *repository.PaymentRepo,
	infiniteClient *infinitepay.Client,
	redisClient *redis.Client,
	cfg *config.Config,
) *PaymentService {
	return &PaymentService{
		raffleRepo:      raffleRepo,
		ticketRepo:      ticketRepo,
		paymentRepo:     paymentRepo,
		infiniteClient:  infiniteClient,
		redisClient:     redisClient,
		cfg:             cfg,
	}
}

type CheckoutInput struct {
	RaffleID   string
	Numbers    []int
	BuyerName  string
	BuyerEmail string
	BuyerPhone string
}

type CheckoutResult struct {
	CheckoutURL string `json:"checkoutUrl"`
	PaymentID   string `json:"paymentId"`
}

func (s *PaymentService) CreateCheckout(ctx context.Context, input CheckoutInput) (*CheckoutResult, error) {
	raffleID, err := primitive.ObjectIDFromHex(input.RaffleID)
	if err != nil {
		return nil, errors.New("invalid raffle id")
	}

	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return nil, ErrRaffleNotFound
	}
	if raffle.Status != model.RaffleStatusActive {
		return nil, ErrRaffleNotActive
	}

	lockKeys := make([]string, len(input.Numbers))
	for i, num := range input.Numbers {
		lockKey := fmt.Sprintf("raffle:%s:reserve:%d", input.RaffleID, num)
		lockKeys[i] = lockKey

		ok, err := s.redisClient.SetNX(ctx, lockKey, input.BuyerEmail, 15*time.Minute).Result()
		if err != nil {
			return nil, fmt.Errorf("redis error: %w", err)
		}
		if !ok {
			return nil, ErrNumbersUnavailable
		}
	}

	defer func() {
		if len(lockKeys) > 0 {
			s.redisClient.Del(ctx, lockKeys...)
		}
	}()

	var ticketIDs []primitive.ObjectID
	for _, num := range input.Numbers {
		ticket, err := s.ticketRepo.FindByRaffleAndNumber(ctx, raffleID, num)
		if err != nil {
			return nil, ErrNumbersUnavailable
		}
		if ticket.Status != model.TicketStatusAvailable {
			return nil, ErrNumbersUnavailable
		}
		ticketIDs = append(ticketIDs, ticket.ID)
	}

	if err := s.ticketRepo.MarkAsReserved(ctx, ticketIDs); err != nil {
		return nil, err
	}

	totalAmount := raffle.TicketPrice * len(input.Numbers)

	payment := &model.Payment{
		RaffleID:   raffleID,
		TicketIDs:  ticketIDs,
		BuyerName:  input.BuyerName,
		BuyerEmail: input.BuyerEmail,
		BuyerPhone: input.BuyerPhone,
		Amount:     totalAmount,
		Status:     model.PaymentStatusPending,
	}

	if err := s.paymentRepo.Insert(ctx, payment); err != nil {
		return nil, err
	}

	checkout, err := s.infiniteClient.CreateCheckout(infinitepay.CreateCheckoutRequest{
		Items: []infinitepay.CheckoutItem{
			{
				Quantity:    len(input.Numbers),
				Price:       totalAmount,
				Description: fmt.Sprintf("%d números - %s", len(input.Numbers), raffle.Title),
			},
		},
		OrderNSU:   payment.ID.Hex(),
		Customer: &infinitepay.Customer{
			Name:        input.BuyerName,
			Email:       input.BuyerEmail,
			PhoneNumber: input.BuyerPhone,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout: %w", err)
	}

	payment.CheckoutURL = checkout.URL
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return &CheckoutResult{
		CheckoutURL: checkout.URL,
		PaymentID:   payment.ID.Hex(),
	}, nil
}

func (s *PaymentService) GetMyPayments(ctx context.Context, email string) ([]model.Payment, error) {
	return s.paymentRepo.FindByEmail(ctx, email)
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, paymentID string) (*model.Payment, error) {
	oid, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return nil, errors.New("invalid payment id")
	}
	return s.paymentRepo.FindByID(ctx, oid)
}

func (s *PaymentService) GetMyTickets(ctx context.Context, email string) ([]model.Ticket, error) {
	return s.ticketRepo.FindPaidByEmail(ctx, email)
}



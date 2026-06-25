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
	ErrNumbersUnavailable   = errors.New("one or more numbers are unavailable")
	ErrRaffleNotActive      = errors.New("raffle is not active")
	ErrOrganizerNotFound    = errors.New("organizer not found")
	ErrInvalidPaymentType   = errors.New("invalid payment type")
)

type PaymentService struct {
	raffleRepo      *repository.RaffleRepo
	ticketRepo      *repository.TicketRepo
	paymentRepo     *repository.PaymentRepo
	userRepo        *repository.UserRepo
	infiniteClient  *infinitepay.Client
	redisClient     *redis.Client
	cfg             *config.Config
}

func NewPaymentService(
	raffleRepo *repository.RaffleRepo,
	ticketRepo *repository.TicketRepo,
	paymentRepo *repository.PaymentRepo,
	userRepo *repository.UserRepo,
	infiniteClient *infinitepay.Client,
	redisClient *redis.Client,
	cfg *config.Config,
) *PaymentService {
	return &PaymentService{
		raffleRepo:      raffleRepo,
		ticketRepo:      ticketRepo,
		paymentRepo:     paymentRepo,
		userRepo:        userRepo,
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

	buyerEmail := input.BuyerEmail
	if buyerEmail == "" {
		buyerEmail = input.BuyerPhone + "@c.rifa"
	}

	organizer, err := s.userRepo.FindByID(ctx, raffle.OrganizerID)
	if err != nil {
		return nil, ErrOrganizerNotFound
	}

	infinitePayHandle := organizer.InfinitePayHandle
	if infinitePayHandle == "" {
		infinitePayHandle = s.cfg.InfinitePayHandle
	}

	payment := &model.Payment{
		Type:       model.PaymentTypeRaffle,
		RaffleID:   raffleID,
		TicketIDs:  ticketIDs,
		BuyerName:  input.BuyerName,
		BuyerEmail: buyerEmail,
		BuyerPhone: input.BuyerPhone,
		Amount:     totalAmount,
		Status:     model.PaymentStatusPending,
	}

	if err := s.paymentRepo.Insert(ctx, payment); err != nil {
		return nil, err
	}

	webhookURL := s.cfg.FrontendURL + "/api/v1/webhooks/infinitepay"

	checkout, err := s.infiniteClient.CreateCheckout(infinitepay.CreateCheckoutRequest{
		Handle: infinitePayHandle,
		Items: []infinitepay.CheckoutItem{
			{
				Quantity:    len(input.Numbers),
				Price:       totalAmount,
				Description: fmt.Sprintf("%d números - %s", len(input.Numbers), raffle.Title),
			},
		},
		OrderNSU:   payment.ID.Hex(),
		WebhookURL: webhookURL,
		Customer: &infinitepay.Customer{
			Name:        input.BuyerName,
			Email:       buyerEmail,
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

type MyPurchaseItem struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Amount        int       `json:"amount"`
	Status        string    `json:"status"`
	BuyerName     string    `json:"buyerName"`
	CreatedAt     time.Time `json:"createdAt"`
	PaidAt        *time.Time `json:"paidAt,omitempty"`
	RaffleID      string    `json:"raffleId,omitempty"`
	RaffleTitle   string    `json:"raffleTitle,omitempty"`
	RaffleStatus  string    `json:"raffleStatus,omitempty"`
	TicketNumbers []int     `json:"ticketNumbers,omitempty"`
}

func (s *PaymentService) GetMyPurchases(ctx context.Context, userID primitive.ObjectID, email string) ([]MyPurchaseItem, error) {
	paymentsByEmail, _ := s.paymentRepo.FindByEmail(ctx, email)
	paymentsByUser, _ := s.paymentRepo.FindByUserID(ctx, userID)

	seen := make(map[string]bool)
	merged := make([]model.Payment, 0)

	for i := range paymentsByEmail {
		id := paymentsByEmail[i].ID.Hex()
		if !seen[id] {
			seen[id] = true
			merged = append(merged, paymentsByEmail[i])
		}
	}
	for i := range paymentsByUser {
		id := paymentsByUser[i].ID.Hex()
		if !seen[id] {
			seen[id] = true
			merged = append(merged, paymentsByUser[i])
		}
	}

	result := make([]MyPurchaseItem, 0, len(merged))
	for _, p := range merged {
		item := MyPurchaseItem{
			ID:        p.ID.Hex(),
			Type:      string(p.Type),
			Amount:    p.Amount,
			Status:    string(p.Status),
			BuyerName: p.BuyerName,
			CreatedAt: p.CreatedAt,
			PaidAt:    p.PaidAt,
		}

		if p.Type == model.PaymentTypeRaffle && p.RaffleID != primitive.NilObjectID {
			item.RaffleID = p.RaffleID.Hex()
			if raffle, err := s.raffleRepo.FindByID(ctx, p.RaffleID); err == nil {
				item.RaffleTitle = raffle.Title
				item.RaffleStatus = string(raffle.Status)
			}
			if len(p.TicketIDs) > 0 {
				tickets, err := s.ticketRepo.FindByIDs(ctx, p.TicketIDs)
				if err == nil {
					for _, t := range tickets {
						item.TicketNumbers = append(item.TicketNumbers, t.Number)
					}
				}
			}
		}

		result = append(result, item)
	}

	return result, nil
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



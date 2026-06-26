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
	ErrNumbersUnavailable    = errors.New("one or more numbers are unavailable")
	ErrRaffleNotActive       = errors.New("raffle is not active")
	ErrOrganizerNotFound     = errors.New("organizer not found")
	ErrInvalidPaymentType    = errors.New("invalid payment type")
	ErrInvalidPaymentID      = errors.New("invalid payment id")
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentNotPending     = errors.New("payment is not pending")
	ErrPaymentNotConfirmed   = errors.New("payment not confirmed by provider")
	ErrPaymentAmountMismatch = errors.New("paid amount is less than expected")
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

type DevCheckoutResult struct {
	PaymentID   string `json:"paymentId"`
	TicketCount int    `json:"ticketCount"`
}

func (s *PaymentService) CreateDevCheckout(ctx context.Context, input CheckoutInput) (*DevCheckoutResult, error) {
	if input.BuyerName == "" || len(input.BuyerName) > 150 {
		return nil, errors.New("buyer name must be between 1 and 150 characters")
	}
	if len(input.BuyerPhone) != 10 && len(input.BuyerPhone) != 11 {
		return nil, errors.New("buyer phone must have 10 or 11 digits")
	}

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
		ok, err := s.redisClient.SetNX(ctx, lockKey, input.BuyerEmail, 5*time.Minute).Result()
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

	totalAmount := raffle.TicketPrice * len(input.Numbers)

	payment := &model.Payment{
		Type:        model.PaymentTypeRaffle,
		RaffleID:    raffleID,
		TicketIDs:   ticketIDs,
		BuyerName:   input.BuyerName,
		BuyerPhone:  input.BuyerPhone,
		Amount:      totalAmount,
		Status:      model.PaymentStatusPaid,
		PaymentMethod: model.PaymentMethodPIX,
	}

	if err := s.paymentRepo.Insert(ctx, payment); err != nil {
		return nil, err
	}

	if _, err := s.ticketRepo.MarkAsPaid(ctx, ticketIDs, input.BuyerName, input.BuyerPhone, payment.ID.Hex()); err != nil {
		return nil, err
	}

	return &DevCheckoutResult{
		PaymentID:   payment.ID.Hex(),
		TicketCount: len(ticketIDs),
	}, nil
}

func (s *PaymentService) CreateCheckout(ctx context.Context, input CheckoutInput) (*CheckoutResult, error) {
	if input.BuyerName == "" || len(input.BuyerName) > 150 {
		return nil, errors.New("buyer name must be between 1 and 150 characters")
	}
	if len(input.BuyerPhone) != 10 && len(input.BuyerPhone) != 11 {
		return nil, errors.New("buyer phone must have 10 or 11 digits")
	}

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

		ok, err := s.redisClient.SetNX(ctx, lockKey, input.BuyerEmail, 5*time.Minute).Result()
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
		BuyerPhone: input.BuyerPhone,
		Amount:     totalAmount,
		Status:     model.PaymentStatusPending,
	}

	if err := s.paymentRepo.Insert(ctx, payment); err != nil {
		return nil, err
	}

	webhookURL := s.cfg.FrontendURL + "/api/v1/webhooks/infinitepay"
	redirectURL := s.cfg.FrontendURL + "/payment/success?paymentId=" + payment.ID.Hex()

	checkout, err := s.infiniteClient.CreateCheckout(infinitepay.CreateCheckoutRequest{
		RedirectURL: redirectURL,
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

func (s *PaymentService) GetMyPurchases(ctx context.Context, userID primitive.ObjectID, phone string) ([]MyPurchaseItem, error) {
	paymentsByPhone, _ := s.paymentRepo.FindByBuyerPhone(ctx, phone)
	paymentsByUser, _ := s.paymentRepo.FindByUserID(ctx, userID)

	seen := make(map[string]bool)
	merged := make([]model.Payment, 0)

	for i := range paymentsByPhone {
		id := paymentsByPhone[i].ID.Hex()
		if !seen[id] {
			seen[id] = true
			merged = append(merged, paymentsByPhone[i])
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

func (s *PaymentService) GetPaymentByID(ctx context.Context, paymentID string) (*model.Payment, error) {
	oid, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return nil, errors.New("invalid payment id")
	}
	return s.paymentRepo.FindByID(ctx, oid)
}

func (s *PaymentService) FindPendingRafflePayments(ctx context.Context) ([]model.Payment, error) {
	return s.paymentRepo.FindPendingRafflePayments(ctx)
}

// resolveRaffleHandle returns the InfinitePay handle used to settle a raffle
// payment: the organizer's own handle, falling back to the platform handle.
func (s *PaymentService) resolveRaffleHandle(ctx context.Context, raffleID primitive.ObjectID) string {
	raffle, err := s.raffleRepo.FindByID(ctx, raffleID)
	if err != nil {
		return s.cfg.InfinitePayHandle
	}
	organizer, err := s.userRepo.FindByID(ctx, raffle.OrganizerID)
	if err != nil || organizer.InfinitePayHandle == "" {
		return s.cfg.InfinitePayHandle
	}
	return organizer.InfinitePayHandle
}

// ConfirmRafflePayment verifies a pending raffle payment directly with
// InfinitePay (order_nsu == payment id) and, only if genuinely paid for at
// least the expected amount, marks the payment and its tickets as paid.
// Idempotent: an already-paid payment returns success without side effects.
// transactionNSU and slug are optional params from the InfinitePay redirect.
func (s *PaymentService) ConfirmRafflePayment(ctx context.Context, paymentID, transactionNSU, slug string) (*model.Payment, error) {
	oid, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return nil, ErrInvalidPaymentID
	}

	payment, err := s.paymentRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	if payment.Type != model.PaymentTypeRaffle {
		return nil, ErrInvalidPaymentType
	}
	if payment.Status == model.PaymentStatusPaid {
		return payment, nil
	}
	if payment.Status != model.PaymentStatusPending {
		return nil, ErrPaymentNotPending
	}

	handle := s.resolveRaffleHandle(ctx, payment.RaffleID)
	check, err := s.infiniteClient.CheckPayment(infinitepay.PaymentCheckRequest{
		Handle:         handle,
		OrderNSU:       payment.ID.Hex(),
		TransactionNSU: transactionNSU,
		Slug:           slug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify payment: %w", err)
	}
	if !check.Success || !check.Paid {
		return nil, ErrPaymentNotConfirmed
	}
	if check.PaidAmount < payment.Amount {
		return nil, ErrPaymentAmountMismatch
	}

	now := time.Now()
	if err := s.paymentRepo.UpdateStatus(ctx, payment.ID, model.PaymentStatusPaid, &now); err != nil {
		return nil, err
	}

	method := model.PaymentMethodPIX
	if check.CaptureMethod == "credit_card" {
		method = model.PaymentMethodCard
	}
	if err := s.paymentRepo.UpdateFields(ctx, payment.ID, primitive.M{"paymentMethod": method}); err != nil {
		return nil, err
	}

	updated, err := s.ticketRepo.MarkAsPaid(ctx, payment.TicketIDs, payment.BuyerName, payment.BuyerPhone, payment.ID.Hex())
	if err != nil {
		return nil, err
	}
	if updated == 0 && len(payment.TicketIDs) > 0 {
		return nil, errors.New("tickets no longer available (reservation expired)")
	}

	payment.Status = model.PaymentStatusPaid
	payment.PaidAt = &now
	payment.PaymentMethod = method
	return payment, nil
}



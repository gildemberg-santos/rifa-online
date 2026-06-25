package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/pkg/infinitepay"
)

var (
	ErrSubscriptionNotActive    = errors.New("subscription is not active")
	ErrNoInfinitePayHandle      = errors.New("infinitePay handle not configured")
	ErrSubscriptionAlreadyActive = errors.New("você já possui uma assinatura ativa")
)

type SubscriptionService struct {
	userRepo       *repository.UserRepo
	paymentRepo    *repository.PaymentRepo
	infiniteClient *infinitepay.Client
	cfg            *config.Config
}

func NewSubscriptionService(
	userRepo *repository.UserRepo,
	paymentRepo *repository.PaymentRepo,
	infiniteClient *infinitepay.Client,
	cfg *config.Config,
) *SubscriptionService {
	return &SubscriptionService{
		userRepo:      userRepo,
		paymentRepo:   paymentRepo,
		infiniteClient: infiniteClient,
		cfg:           cfg,
	}
}

type SubscriptionCheckoutResult struct {
	CheckoutURL string `json:"checkoutUrl,omitempty"`
	IsTrial     bool   `json:"isTrial"`
}

func (s *SubscriptionService) CreateSubscriptionCheckout(ctx context.Context, userID string) (*SubscriptionCheckoutResult, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.userRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.SubscriptionStatus == model.SubscriptionStatusActive {
		if user.SubscriptionExpiresAt == nil || time.Now().Before(*user.SubscriptionExpiresAt) {
			return nil, ErrSubscriptionAlreadyActive
		}
	}

	pending, err := s.paymentRepo.FindPendingSubscriptionByUserID(ctx, oid)
	if err != nil {
		return nil, err
	}
	for _, p := range pending {
		now := time.Now()
		s.paymentRepo.UpdateStatus(ctx, p.ID, model.PaymentStatusExpired, &now)
	}

	if !user.HasSubscriptionBefore && user.SubscriptionStatus != model.SubscriptionStatusActive {
		expiresAt := time.Now().AddDate(0, 0, model.TrialDays)

		if err := s.userRepo.UpdateFields(ctx, oid, primitive.M{
			"subscriptionStatus":      model.SubscriptionStatusActive,
			"subscriptionExpiresAt":   expiresAt,
			"subscriptionIsTrial":     true,
			"hasSubscriptionBefore":   true,
		}); err != nil {
			return nil, err
		}

		return &SubscriptionCheckoutResult{
			IsTrial: true,
		}, nil
	}

	payment := &model.Payment{
		Type:       model.PaymentTypeSubscription,
		UserID:     oid,
		BuyerEmail: user.Email,
		Amount:     model.SubscriptionPrice,
		Status:     model.PaymentStatusPending,
	}

	if err := s.paymentRepo.Insert(ctx, payment); err != nil {
		return nil, err
	}

	orderNSU := "sub_" + payment.ID.Hex()

	checkout, err := s.infiniteClient.CreateCheckout(infinitepay.CreateCheckoutRequest{
		Items: []infinitepay.CheckoutItem{
			{
				Quantity:    1,
				Price:       model.SubscriptionPrice,
				Description: "Assinatura mensal - Rifa Online",
			},
		},
		OrderNSU: orderNSU,
		Customer: &infinitepay.Customer{
			Name:  user.Name,
			Email: user.Email,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription checkout: %w", err)
	}

	payment.CheckoutURL = checkout.URL
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return &SubscriptionCheckoutResult{
		CheckoutURL: checkout.URL,
		IsTrial:     false,
	}, nil
}

func (s *SubscriptionService) ActivateSubscription(ctx context.Context, payment *model.Payment) error {
	now := time.Now()
	expiresAt := now.AddDate(0, 1, 0)

	if err := s.userRepo.UpdateFields(ctx, payment.UserID, primitive.M{
		"subscriptionStatus":     model.SubscriptionStatusActive,
		"subscriptionExpiresAt":  expiresAt,
		"subscriptionIsTrial":    false,
		"hasSubscriptionBefore":  true,
	}); err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) CheckSubscription(ctx context.Context, userID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrUserNotFound
	}

	user, err := s.userRepo.FindByID(ctx, oid)
	if err != nil {
		return ErrUserNotFound
	}

	if user.SubscriptionStatus != model.SubscriptionStatusActive {
		return ErrSubscriptionNotActive
	}

	if user.SubscriptionExpiresAt != nil && time.Now().After(*user.SubscriptionExpiresAt) {
		s.userRepo.UpdateFields(ctx, oid, primitive.M{
			"subscriptionStatus": model.SubscriptionStatusPastDue,
		})
		return ErrSubscriptionNotActive
	}

	return nil
}

func (s *SubscriptionService) GetStatus(ctx context.Context, userID string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.SubscriptionExpiresAt != nil && time.Now().After(*user.SubscriptionExpiresAt) {
		if user.SubscriptionStatus == model.SubscriptionStatusActive {
			s.userRepo.UpdateFields(ctx, oid, primitive.M{
				"subscriptionStatus": model.SubscriptionStatusPastDue,
			})
			user.SubscriptionStatus = model.SubscriptionStatusPastDue
		}
	}

	return user, nil
}

func (s *SubscriptionService) CreateDevSubscriptionCheckout(ctx context.Context, userID string) (*SubscriptionCheckoutResult, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.userRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.SubscriptionStatus == model.SubscriptionStatusActive {
		if user.SubscriptionExpiresAt == nil || time.Now().Before(*user.SubscriptionExpiresAt) {
			return nil, ErrSubscriptionAlreadyActive
		}
	}

	pending, err := s.paymentRepo.FindPendingSubscriptionByUserID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if len(pending) > 0 {
		// Expire pending payments
		for _, p := range pending {
			now := time.Now()
			s.paymentRepo.UpdateStatus(ctx, p.ID, model.PaymentStatusExpired, &now)
		}
	}

	payment := &model.Payment{
		Type:       model.PaymentTypeSubscription,
		UserID:     oid,
		BuyerEmail: user.Email,
		Amount:     model.SubscriptionPrice,
		Status:     model.PaymentStatusPaid,
		PaymentMethod: model.PaymentMethodPIX,
	}

	if err := s.paymentRepo.Insert(ctx, payment); err != nil {
		return nil, err
	}

	now := time.Now()
	if err := s.paymentRepo.UpdateStatus(ctx, payment.ID, model.PaymentStatusPaid, &now); err != nil {
		return nil, err
	}

	if err := s.ActivateSubscription(ctx, payment); err != nil {
		return nil, err
	}

	return &SubscriptionCheckoutResult{
		IsTrial: false,
	}, nil
}

func (s *SubscriptionService) UpdateInfinitePayHandle(ctx context.Context, userID string, handle string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := s.userRepo.UpdateFields(ctx, oid, primitive.M{
		"infinitePayHandle": handle,
	}); err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, oid)
}

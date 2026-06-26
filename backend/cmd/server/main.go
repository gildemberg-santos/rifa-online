package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/crypto"
	"github.com/user/rifa-online/internal/handler"
	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/migrations"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/internal/service"
	"github.com/user/rifa-online/pkg/infinitepay"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	cfg := config.Load()

	var logHandler slog.Handler
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	if cfg.LogFormat == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	logger.Info("starting server", "port", cfg.Port, "log_format", cfg.LogFormat, "log_level", cfg.LogLevel)

	// Segurança: não permite subir fora de desenvolvimento com um JWT_SECRET fraco/padrão,
	// o que permitiria forjar tokens e burlar a autenticação.
	if cfg.AppEnv != "development" && (cfg.JWTSecret == "" || cfg.JWTSecret == "change-me" || len(cfg.JWTSecret) < 32) {
		logger.Error("insecure JWT_SECRET: defina um segredo forte (>= 32 caracteres) em produção")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Error("failed to connect to mongodb", "error", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.Background())

	optsRedis, err := redis.ParseURL(cfg.RedisURI)
	if err != nil {
		logger.Error("failed to parse redis url", "error", err)
		os.Exit(1)
	}
	redisClient := redis.NewClient(optsRedis)
	defer redisClient.Close()

	db := mongoClient.Database(cfg.MongoDBName)

	// Criptografia de campos sensíveis em repouso. Em produção, a ausência de
	// chave deixa os dados em texto puro — alertamos para que seja configurada.
	dataCipher, err := crypto.New(cfg.DataEncryptionKey, cfg.BlindIndexKey)
	if err != nil {
		logger.Error("failed to init data cipher", "error", err)
		os.Exit(1)
	}
	if cfg.AppEnv != "development" && !dataCipher.Enabled() {
		logger.Warn("DATA_ENCRYPTION_KEY ausente: dados sensíveis serão gravados sem criptografia")
	}

	userRepo := repository.NewUserRepo(db, dataCipher)
	raffleRepo := repository.NewRaffleRepo(db)
	ticketRepo := repository.NewTicketRepo(db, dataCipher)
	paymentRepo := repository.NewPaymentRepo(db, dataCipher)
	webhookRepo := repository.NewWebhookRepo(db)
	contactRepo := repository.NewContactRepo(db, dataCipher)

	if err := migrations.Run(ctx, db); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	reservationTTL := 10 * time.Minute
	cleanupInterval := 5 * time.Minute
	go cleanupExpiredReservations(logger, ticketRepo, paymentRepo, reservationTTL, cleanupInterval)

	seedDefaultUser(userRepo)

	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)

	webhookURL := cfg.FrontendURL + "/api/v1/webhooks/infinitepay"
	infiniteClient := infinitepay.NewClient(cfg.InfinitePayHandle, webhookURL, cfg.FrontendURL, cfg.InfinitePayBaseURL)

	raffleService := service.NewRaffleService(raffleRepo, ticketRepo, paymentRepo, userRepo)
	paymentService := service.NewPaymentService(raffleRepo, ticketRepo, paymentRepo, userRepo, infiniteClient, redisClient, cfg)
	subscriptionService := service.NewSubscriptionService(userRepo, paymentRepo, infiniteClient, cfg)

	raffleHandler := handler.NewRaffleHandler(raffleService, userRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService, paymentRepo, ticketRepo, userRepo)
	webhookHandler := handler.NewWebhookHandler(paymentService, subscriptionService, webhookRepo, logger)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	adminHandler := handler.NewAdminHandler(userRepo, raffleRepo, ticketRepo, paymentRepo)
	contactHandler := handler.NewContactHandler(contactRepo)

	authMw := middleware.Auth(cfg)

	r := chi.NewRouter()

	r.Use(middleware.StructuredLogger(logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			// Limita tentativas de cadastro/login por IP (anti brute force / spam).
			r.Group(func(r chi.Router) {
				r.Use(middleware.RateLimit(10, time.Minute))
				r.Post("/register", authHandler.Register)
				r.Post("/login", authHandler.Login)
			})
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/logout", authHandler.Logout)
		})

		subMw := middleware.RequiresSubscription(userRepo)

		r.Route("/raffles", func(r chi.Router) {
			r.Get("/", raffleHandler.List)
			r.Get("/{id}", raffleHandler.GetDetail)

			r.Group(func(r chi.Router) {
				r.Use(authMw)
				r.Use(subMw)
				r.Post("/{id}/checkout", paymentHandler.Checkout)
				if cfg.AppEnv == "development" {
					r.Post("/{id}/dev-checkout", paymentHandler.DevCheckout)
				}
			})

			r.Group(func(r chi.Router) {
				r.Use(authMw)
				r.Post("/", raffleHandler.Create)
				r.Put("/{id}", raffleHandler.Update)
				r.Patch("/{id}/cancel", raffleHandler.Cancel)
				r.Delete("/{id}", raffleHandler.Delete)
				r.Post("/{id}/draw", raffleHandler.Draw)
				r.Get("/my", raffleHandler.MyRaffles)
				r.Get("/{id}/stats", raffleHandler.Stats)
			})
		})

		r.Route("/payments", func(r chi.Router) {
			r.Get("/my", paymentHandler.MyPayments)
			r.Get("/{id}", paymentHandler.GetPayment)
			r.Get("/my/tickets", paymentHandler.MyTickets)

			r.Group(func(r chi.Router) {
				r.Use(authMw)
				r.Post("/{id}/confirm", paymentHandler.ConfirmPayment)
			})
		})

		r.Post("/webhooks/infinitepay", webhookHandler.HandleInfinitePay)

		r.With(middleware.RateLimit(5, time.Minute)).Post("/contact", contactHandler.Create)

		r.Route("/dashboard", func(r chi.Router) {
			r.Use(authMw)
			r.Get("/stats", raffleHandler.DashboardStats)
		})

		r.Route("/me", func(r chi.Router) {
			r.Use(authMw)
			r.Get("/", authHandler.GetProfile)
			r.Put("/", authHandler.UpdateProfile)
			r.Put("/infinite-pay-handle", subscriptionHandler.UpdateInfinitePayHandle)
			r.Get("/purchases", paymentHandler.MyPurchases)
		})

		r.Route("/subscription", func(r chi.Router) {
			r.Use(authMw)
			r.Post("/checkout", subscriptionHandler.Checkout)
			if cfg.AppEnv == "development" {
				r.Post("/dev-checkout", subscriptionHandler.DevCheckout)
			}
			r.Get("/status", subscriptionHandler.Status)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(authMw)
			r.Use(middleware.Admin(userRepo))
			r.Get("/users", adminHandler.Users)
			r.Get("/users/{id}", adminHandler.UserDetails)
			r.Put("/users/{id}/subscription", adminHandler.UpdateUserSubscription)
			r.Get("/raffles", adminHandler.Raffles)
			r.Get("/stats", adminHandler.Stats)
			r.Get("/contact-messages", contactHandler.List)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		checks := map[string]string{
			"server": "ok",
		}

		if err := mongoClient.Ping(ctx, nil); err != nil {
			checks["mongodb"] = fmt.Sprintf("error: %v", err)
		} else {
			checks["mongodb"] = "ok"
		}

		if err := redisClient.Ping(ctx).Err(); err != nil {
			checks["redis"] = fmt.Sprintf("error: %v", err)
		} else {
			checks["redis"] = "ok"
		}

		status := http.StatusOK
		for _, v := range checks {
			if v != "ok" {
				status = http.StatusServiceUnavailable
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(checks)
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server starting", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced shutdown", "error", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}

func seedDefaultUser(userRepo *repository.UserRepo) {
	ctx := context.Background()
	email := "admin@email.com"

	_, err := userRepo.FindByEmail(ctx, email)
	if err == nil {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), 12)
	if err != nil {
		slog.Warn("failed to hash default user password", "error", err)
		return
	}

	now := time.Now()
	expiresAt := now.AddDate(1, 0, 0)
	user := &model.User{
		Name:                  "Administrador",
		Email:                 email,
		PasswordHash:          string(hash),
		Role:                  model.RoleAdmin,
		SubscriptionStatus:    model.SubscriptionStatusActive,
		SubscriptionExpiresAt: &expiresAt,
		HasSubscriptionBefore: true,
	}

	if err := userRepo.Insert(ctx, user); err != nil {
		slog.Warn("failed to seed default user", "error", err)
		return
	}

	slog.Info("default admin user created", "email", email)
}

func cleanupExpiredReservations(logger *slog.Logger, ticketRepo *repository.TicketRepo, paymentRepo *repository.PaymentRepo, ttl, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		cutoff := time.Now().Add(-ttl)

		releasedCount := 0
		tickets, err := ticketRepo.FindReservedOlderThan(ctx, cutoff)
		if err != nil {
			logger.Error("failed to find expired reservations", "error", err)
			cancel()
			continue
		}

		if len(tickets) > 0 {
			ids := make([]primitive.ObjectID, len(tickets))
			for i, t := range tickets {
				ids[i] = t.ID
			}
			if err := ticketRepo.ReleaseReservations(ctx, ids); err != nil {
				logger.Error("failed to release expired reservations", "error", err)
			} else {
				releasedCount = len(ids)
			}
		}

		expiredCount, err := paymentRepo.ExpirePendingOlderThan(ctx, cutoff)
		if err != nil {
			logger.Error("failed to expire old pending payments", "error", err)
		}

		if releasedCount > 0 || expiredCount > 0 {
			logger.Info("reservation cleanup complete", "released", releasedCount, "expired_payments", expiredCount)
		}

		cancel()
	}
}

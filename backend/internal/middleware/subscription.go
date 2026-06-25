package middleware

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
)

func RequiresSubscription(userRepo *repository.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := UserIDFromContext(r.Context())
			if userID == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			oid, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				http.Error(w, `{"error":"invalid user"}`, http.StatusUnauthorized)
				return
			}

			user, err := userRepo.FindByID(r.Context(), oid)
			if err != nil {
				http.Error(w, `{"error":"user not found"}`, http.StatusUnauthorized)
				return
			}

			if user.Role == model.RoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			if user.SubscriptionStatus != model.SubscriptionStatusActive {
				http.Error(w, `{"error":"subscription is not active"}`, http.StatusForbidden)
				return
			}

			if user.SubscriptionExpiresAt != nil && time.Now().After(*user.SubscriptionExpiresAt) {
				userRepo.UpdateFields(r.Context(), oid, bson.M{
					"subscriptionStatus": model.SubscriptionStatusPastDue,
				})
				http.Error(w, `{"error":"subscription is not active"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

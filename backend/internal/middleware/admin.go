package middleware

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
)

func Admin(userRepo *repository.UserRepo) func(http.Handler) http.Handler {
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
			if err != nil || user.Role != model.RoleAdmin {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

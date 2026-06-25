package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/user/rifa-online/internal/auth"
	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserNotFound           = errors.New("user not found")
)

type AuthService struct {
	userRepo *repository.UserRepo
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepo, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type AuthResult struct {
	User         *model.User
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if len(input.Name) < 2 || len(input.Name) > 100 {
		return nil, errors.New("name must be between 2 and 100 characters")
	}
	if len(input.Email) > 255 {
		return nil, errors.New("email must be at most 255 characters")
	}
	if len(input.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}
	if len(input.Password) > 72 {
		return nil, errors.New("password must be at most 72 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         model.RoleUser,
	}

	if err := s.userRepo.Insert(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return s.userRepo.FindByID(ctx, oid)
}

type UpdateProfileInput struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	updates := bson.M{}

	if input.Name != "" {
		if len(input.Name) < 2 || len(input.Name) > 100 {
			return nil, errors.New("name must be between 2 and 100 characters")
		}
		updates["name"] = strings.TrimSpace(input.Name)
	}
	if input.Email != "" {
		if len(input.Email) > 255 {
			return nil, errors.New("email must be at most 255 characters")
		}
		email := strings.TrimSpace(strings.ToLower(input.Email))
		existing, err := s.userRepo.FindByEmail(ctx, email)
		if err == nil && existing.ID != oid {
			return nil, ErrEmailAlreadyRegistered
		}
		updates["email"] = email
	}
	if input.Phone != "" {
		if len(input.Phone) < 10 || len(input.Phone) > 11 {
			return nil, errors.New("phone must have 10 or 11 digits")
		}
		updates["phone"] = input.Phone
	}
	if input.Password != "" {
		if len(input.Password) < 6 {
			return nil, errors.New("password must be at least 6 characters")
		}
		if len(input.Password) > 72 {
			return nil, errors.New("password must be at most 72 characters")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		if err != nil {
			return nil, err
		}
		updates["passwordHash"] = string(hash)
	}

	if len(updates) == 0 {
		return s.GetProfile(ctx, userID)
	}

	if err := s.userRepo.UpdateFields(ctx, oid, updates); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, err
	}

	return s.userRepo.FindByID(ctx, oid)
}

func (s *AuthService) RefreshToken(ctx context.Context, tokenStr string) (*AuthResult, error) {
	claims, err := auth.ValidateToken(tokenStr, s.cfg.JWTSecret)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	oid, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

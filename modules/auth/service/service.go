package service

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zharf/payment-gateway-simulator/internal/platform/config"
	"github.com/zharf/payment-gateway-simulator/internal/platform/security"
	audit "github.com/zharf/payment-gateway-simulator/modules/auditlog/service"
	"github.com/zharf/payment-gateway-simulator/modules/auth/dto"
	"github.com/zharf/payment-gateway-simulator/modules/auth/entity"
	"github.com/zharf/payment-gateway-simulator/modules/auth/repository"
	"github.com/zharf/payment-gateway-simulator/modules/auth/validator"
	"gorm.io/gorm"
)

type Service struct {
	cfg   config.Config
	repo  *repository.Repository
	redis *redis.Client
	audit *audit.Service
}

func New(cfg config.Config, repo *repository.Repository, redis *redis.Client, audit *audit.Service) *Service {
	return &Service{cfg: cfg, repo: repo, redis: redis, audit: audit}
}

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) (*dto.TokenResponse, error) {
	if err := validator.Register(req); err != nil {
		return nil, err
	}
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &entity.User{Name: req.Name, Email: req.Email, PasswordHash: hash}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	s.audit.Record(ctx, user.ID, "", "auth.register", "users", map[string]string{"email": user.Email})
	return s.issueTokens(ctx, user)
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error) {
	if err := validator.Login(req); err != nil {
		return nil, err
	}
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil || !security.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}
	s.audit.Record(ctx, user.ID, "", "auth.login", "users", map[string]string{"email": user.Email})
	return s.issueTokens(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := security.ParseJWT(refreshToken, s.cfg.RefreshSecret)
	if err != nil || claims.Type != "refresh" {
		return nil, errors.New("invalid refresh token")
	}
	if s.redis != nil {
		key := "refresh:" + claims.UserID + ":" + claims.ID
		exists, _ := s.redis.Exists(ctx, key).Result()
		if exists == 0 {
			return nil, errors.New("refresh token revoked")
		}
		_ = s.redis.Del(ctx, key).Err()
	}
	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, user)
}

func (s *Service) Logout(ctx context.Context, accessJTI, userID, refreshToken string) error {
	if s.redis == nil {
		return nil
	}
	_ = s.redis.Set(ctx, "jwt:blacklist:"+accessJTI, "1", s.cfg.AccessTTL).Err()
	if refreshToken != "" {
		if claims, err := security.ParseJWT(refreshToken, s.cfg.RefreshSecret); err == nil {
			_ = s.redis.Del(ctx, "refresh:"+claims.UserID+":"+claims.ID).Err()
		}
	}
	s.audit.Record(ctx, userID, "", "auth.logout", "users", nil)
	return nil
}

func (s *Service) issueTokens(ctx context.Context, user *entity.User) (*dto.TokenResponse, error) {
	access, _, err := security.GenerateJWT(user.ID, user.Email, "access", s.cfg.AccessSecret, s.cfg.AccessTTL)
	if err != nil {
		return nil, err
	}
	refresh, refreshJTI, err := security.GenerateJWT(user.ID, user.Email, "refresh", s.cfg.RefreshSecret, s.cfg.RefreshTTL)
	if err != nil {
		return nil, err
	}
	if s.redis != nil {
		_ = s.redis.Set(ctx, "refresh:"+user.ID+":"+refreshJTI, "1", s.cfg.RefreshTTL).Err()
	}
	return &dto.TokenResponse{AccessToken: access, RefreshToken: refresh, TokenType: "Bearer", ExpiresIn: int64(s.cfg.AccessTTL / time.Second)}, nil
}

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/noctisnyxx/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type ITokenUsecase interface {
	GenerateAccessToken(ctx context.Context, claims domain.TokenClaims) (string, error)
	GenerateRefreshToken(ctx context.Context, claims domain.TokenClaims) (string, error)
	Revoke(ctx context.Context, token string) error
}

type TokenUsecase struct {
	logger *slog.Logger
	repo   domain.ITokenRepository
}

func NewTokenUsecase(logger *slog.Logger, repo domain.ITokenRepository) *TokenUsecase {
	return &TokenUsecase{
		logger: logger,
		repo:   repo,
	}
}

func (uc *TokenUsecase) GenerateAccessToken(ctx context.Context, claims domain.TokenClaims) (string, error) {
	uc.logger.Debug("Call Generate Access Token")
	aToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id":   claims.UserId,
		"client_id": claims.ClientId,
		"exp":       claims.ExpireDuration,
	})
	aTokenString, err := aToken.SignedString("")
	if err != nil {
		uc.logger.Error("failed to sign access token", "error", err)
		return "", err
	}
	return aTokenString, nil
}

func (uc *TokenUsecase) GenerateRefreshToken(ctx context.Context, claims domain.TokenClaims) (string, error) {
	uc.logger.Debug("Call Generate Refresh Token")
	rToken := make([]byte, 32)
	_, err := rand.Read(rToken)
	if err != nil {
		return "", err
	}
	rTokenString := base64.URLEncoding.EncodeToString(rToken)
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(rTokenString), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("failed to hash password", "error", err)
		return "", err
	}
	hashedRTokenString := string(hashedRefreshToken)
	RTokenDomain := domain.RefreshToken{
		Token:    hashedRTokenString,
		UserId:   claims.UserId,
		ClientId: claims.ClientId,
	}
	if err := uc.repo.InsertRefreshToken(ctx, RTokenDomain); err != nil {
		return "", err
	}
	return rTokenString, nil
}

func (uc *TokenUsecase) RevokeToken(ctx context.Context, token string) error {
	uc.logger.Debug("Call Revoke Token")
	if err := uc.repo.RevokeRefreshToken(ctx, token); err != nil {
		return err
	}
	return nil
}

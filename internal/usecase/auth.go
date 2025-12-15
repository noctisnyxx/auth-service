package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"time"

	"github.com/noctisnyxx/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUsecase interface {
	BasicLogin(ctx context.Context, username, password, clientId string) (string, error)
}

type AuthUsecase struct {
	logger   *slog.Logger
	userUC   IUserUsecase
	clientUC IClientUsecase
	tokenUC  ITokenUsecase
	repo     domain.IAuthRepository
}

func NewAuthUsecase(logger *slog.Logger, userUC IUserUsecase, clientUC IClientUsecase, tokenUC ITokenUsecase, repo domain.IAuthRepository) *AuthUsecase {
	return &AuthUsecase{
		logger:   logger,
		userUC:   userUC,
		clientUC: clientUC,
		tokenUC:  tokenUC,
		repo:     repo,
	}
}

func (uc *AuthUsecase) BasicLogin(ctx context.Context, username, password, clientId string) (string, error) {
	uc.logger.Debug("Call Basic Login")
	users, err := uc.userUC.Find(ctx, domain.FindUserQuery{Username: &username})
	if err != nil {
		uc.logger.Error("failed to find user by username", "error", err)
		return "", err
	}
	if len(users) == 0 {
		return "", domain.ErrInvalidCredentials
	}
	user := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		uc.logger.Error("invalid password", "error", err)
		return "", domain.ErrInvalidCredentials
	}
	authCode, err := uc.GenerateAuthCode(ctx)
	if err != nil {
		return "", err
	}
	domainAuthCode := domain.AuthCode{
		Code:      authCode,
		UserID:    user.ID,
		ClientID:  clientId,
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}
	if err := uc.repo.InsertCode(ctx, domainAuthCode); err != nil {
		return "", err
	}
	return authCode, nil
}

func (uc *AuthUsecase) GenerateAuthCode(ctx context.Context) (string, error) {
	authCode := make([]byte, 32)
	_, err := rand.Read(authCode)
	if err != nil {
		return "", err
	}
	authCodeString := base64.URLEncoding.EncodeToString(authCode)
	return authCodeString, nil
}

func (uc *AuthUsecase) ExchangeCode(ctx context.Context, code, clientID, clientSecret string) (domain.Token, error) {
	clients, err := uc.clientUC.Find(ctx, domain.FindClientQuery{
		ID:     &clientID,
		Secret: &clientSecret,
	})
	if err != nil {
		return domain.Token{}, err
	}
	if len(clients) == 0 {
		return domain.Token{}, domain.ErrInvalidCredentials
	}

	authCodes, err := uc.repo.FindAuthCode(ctx, code)
	if err != nil {
		return domain.Token{}, err
	}
	if len(authCodes) == 0 {
		return domain.Token{}, domain.ErrInvalidCredentials
	}
	authCode := authCodes[0]

	if authCode.ClientID != clientID {
		return domain.Token{}, domain.ErrInvalidCredentials
	}

	if authCode.ExpiresAt.Before(time.Now()) {
		return domain.Token{}, domain.ErrInvalidCredentials
	}

	if err := uc.repo.DeleteAuthCode(ctx, code); err != nil {
		return domain.Token{}, err
	}

	accessTokenClaims := domain.TokenClaims{
		UserId:         authCode.UserID,
		ClientId:       authCode.ClientID,
		ExpireDuration: 15 * time.Minute,
	}
	accessToken, err := uc.tokenUC.GenerateAccessToken(ctx, accessTokenClaims)
	if err != nil {
		return domain.Token{}, err
	}
	refreshTokenClaims := domain.TokenClaims{
		UserId:         authCode.UserID,
		ClientId:       authCode.ClientID,
		ExpireDuration: 7 * 24 * time.Hour,
	}
	refreshToken, err := uc.tokenUC.GenerateRefreshToken(ctx, refreshTokenClaims)
	if err != nil {
		return domain.Token{}, err
	}
	return domain.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/noctisnyxx/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	RegisterUser(context.Context, *domain.User) error
	FindUser(context.Context, domain.FindUserQuery) ([]domain.User, error)
}

type UserUsecase struct {
	logger *slog.Logger
	repo   domain.IUserRepository
}

func NewUserUsecase(logger *slog.Logger, repo domain.IUserRepository) *UserUsecase {
	return &UserUsecase{
		logger: logger,
		repo:   repo,
	}
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, u *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("failed to hash password", "error", err)
		return err
	}
	u.Password = string(hashedPassword)
	u.ID = uuid.New().String()
	return uc.repo.Insert(ctx, u)
}

func (uc *UserUsecase) FindUser(ctx context.Context, q domain.FindUserQuery) ([]domain.User, error) {
	return uc.repo.FindUsers(ctx, q)
}

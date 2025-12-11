package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/noctisnyxx/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Register(context.Context, *domain.User) error
	Find(context.Context, domain.FindUserQuery) ([]domain.User, error)
	Update(ctx context.Context, id string, u domain.UserUpdate) error
	Delete(ctx context.Context, id string) error
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

func (uc *UserUsecase) Register(ctx context.Context, u *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("failed to hash password", "error", err)
		return err
	}
	u.Password = string(hashedPassword)
	u.ID = uuid.New().String()
	return uc.repo.Insert(ctx, u)
}

func (uc *UserUsecase) Find(ctx context.Context, q domain.FindUserQuery) ([]domain.User, error) {
	return uc.repo.Find(ctx, q)
}

func (uc *UserUsecase) Update(ctx context.Context, id string, u domain.UserUpdate) error {
	f := false
	if u.Username != nil {
		users, err := uc.Find(ctx, domain.FindUserQuery{Username: u.Username})
		if err != nil {
			uc.logger.Error("failed to find user by username", "error", err)
			return err
		}
		if len(users) > 0 && id != users[0].ID {
			return domain.ErrUserAlreadyExists
		}
	}
	if u.Email != nil {
		users, err := uc.Find(ctx, domain.FindUserQuery{Email: u.Email})
		if err != nil {
			uc.logger.Error("failed to find user by email", "error", err)
			return err
		}
		if len(users) > 0 && id != users[0].ID {
			return domain.ErrUserAlreadyExists
		}
		u.EmailVerified = &f
	}
	if u.Phone != nil {
		u.PhoneVerified = &f
	}
	if u.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
		if err != nil {
			uc.logger.Error("failed to hash password", "error", err)
			return err
		}
		hashedPasswordStr := string(hashedPassword)
		u.Password = &hashedPasswordStr
	}
	return uc.repo.Update(ctx, id, u)
}

func (uc *UserUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

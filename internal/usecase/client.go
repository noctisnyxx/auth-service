package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"strings"

	"github.com/noctisnyxx/auth-service/internal/domain"
)

type IClientUsecase interface {
	Register(ctx context.Context, c *domain.Client) error
	Find(ctx context.Context, q domain.FindClientQuery) ([]domain.Client, error)
	Update(ctx context.Context, id string, u domain.ClientUpdate) error
	Delete(ctx context.Context, id string) error
}

type ClientUsecase struct {
	logger *slog.Logger
	repo   domain.IClientRepository
}

func NewClientUsecase(logger *slog.Logger, repo domain.IClientRepository) *ClientUsecase {
	return &ClientUsecase{
		logger: logger,
		repo:   repo,
	}
}

func (uc *ClientUsecase) Register(ctx context.Context, c *domain.Client) error {
	id, err := uc.generateId(c.Name)
	if err != nil {
		return err
	}
	secret, err := uc.generateSecret()
	if err != nil {
		return err
	}
	c.ID = id
	c.Secret = secret
	return uc.repo.Insert(ctx, c)
}

func (uc *ClientUsecase) Find(ctx context.Context, q domain.FindClientQuery) ([]domain.Client, error) {
	return uc.repo.Find(ctx, q)
}

func (uc *ClientUsecase) Update(ctx context.Context, id string, u domain.ClientUpdate) error {
	return uc.repo.Update(ctx, id, u)
}

func (uc *ClientUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *ClientUsecase) generateSecret() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (uc *ClientUsecase) generateId(name string) (string, error) {
	name = strings.Join(strings.Fields(name), " ")
	var id string
	for _, c := range name {
		if c == ' ' {
			id += "-"
		} else {
			id += string(c)
		}
	}
	clients, err := uc.Find(context.Background(), domain.FindClientQuery{
		ID: &id,
	})
	if err != nil {
		uc.logger.Error("failed to find client by id", "error", err)
		return "", err
	}
	if len(clients) > 0 {
		return "", domain.ErrDuplicateClientName
	}
	return id, nil
}

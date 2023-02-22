package auth

import (
	"context"
	"errors"
	"models/ent"
	"models/ent/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService struct {
	client *ent.Client
	logger *zap.Logger

	secret string
}

func NewAuthService(client *ent.Client, secret string, logger *zap.Logger) *AuthService {
	return &AuthService{client: client, secret: secret, logger: logger}
}

func (s *AuthService) GetUserFromId(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	usr, err := s.client.User.Query().Where(user.ID(id)).Only(ctx)
	if err != nil {
		s.logger.Error("user not found", zap.Error(err))
		return nil, errors.New("user not found")
	}

	return usr, nil
}

func (s *AuthService) GetSecret() string {
	return s.secret
}

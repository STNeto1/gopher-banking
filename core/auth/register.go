package auth

import (
	"context"
	"errors"
	"models/ent"
	"models/ent/user"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserPayload struct {
	Name     string
	Email    string
	Password string
}

func (s *AuthService) RegisterUser(ctx context.Context, payload RegisterUserPayload) (*ent.User, error) {
	_, err := s.client.User.Query().Where(user.Email(payload.Email)).Only(ctx)
	if err == nil {
		s.logger.Error("user already exists", zap.Error(err))
		return nil, errors.New("user already exists")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("error while hashing password", zap.Error(err))
		return nil, errors.New("error while hashing password")
	}

	usr, err := s.client.User.
		Create().
		SetName(payload.Name).
		SetEmail(payload.Email).
		SetPassword(string(pwdHash)).
		SetBalance(0.0).
		Save(ctx)

	if err != nil {
		s.logger.Error("error while creating user", zap.Error(err))
		return nil, errors.New("error while creating user")
	}

	return usr, nil
}

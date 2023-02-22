package auth

import (
	"context"
	"errors"
	"models/ent"
	"models/ent/user"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*ent.User, error) {
	usr, err := s.client.User.Query().Where(user.Email(email)).Only(ctx)
	if err != nil {
		s.logger.Error("user not found", zap.Error(err))
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		s.logger.Error("invalid password", zap.Error(err))
		return nil, errors.New("invalid password")
	}

	return usr, nil
}

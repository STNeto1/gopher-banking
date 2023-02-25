package auth

import (
	"context"
	"errors"
	"fmt"
	"models/ent"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *AuthService) GenerateToken(user *ent.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		s.logger.Error("error while signing token", zap.Error(err))
		return "", err
	}

	return t, nil
}

func (s *AuthService) ValidateToken(token string) (*jwt.StandardClaims, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		s.logger.Error("error while parsing token", zap.Error(err))
		return nil, err
	}

	if !t.Valid {
		s.logger.Error("invalid token")
		return nil, err
	}

	return t.Claims.(*jwt.StandardClaims), nil
}

func (s *AuthService) ExtractUser(raw interface{}) (*ent.User, error) {
	valStr := fmt.Sprint(raw)
	if valStr == "" {
		return nil, errors.New("")
	}

	claims, err := s.ValidateToken(valStr)
	if err != nil {
		return nil, errors.New("")
	}

	parsedUserId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("")
	}

	usr, err := s.GetUserFromId(context.Background(), parsedUserId)
	if err != nil {
		return nil, errors.New("")
	}

	return usr, nil
}

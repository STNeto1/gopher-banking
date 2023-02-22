package auth_test

import (
	"context"
	"core/auth"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCreatingUser(t *testing.T) {
	s, client, l := CreateAuthService(t)
	defer client.Close()
	defer l.Sync()

	_, err := s.RegisterUser(context.Background(), auth.RegisterUserPayload{
		Name:     "John Doe",
		Email:    "some-mail@mail.com",
		Password: "some-password",
	})

	assert.Nil(t, err)
}

func TestCreatingUserWithExistingEmail(t *testing.T) {
	s, client, l := CreateAuthService(t)
	defer client.Close()
	defer l.Sync()

	existingMail := "existing@mail.com"

	_, err := client.User.
		Create().
		SetName("John Doe").
		SetEmail(existingMail).
		SetPassword("some-password").
		SetBalance(0).
		Save(context.Background())
	assert.Nil(t, err)

	_, err = s.RegisterUser(context.Background(), auth.RegisterUserPayload{
		Name:     "John Doe",
		Email:    existingMail,
		Password: "some-password",
	})

	assert.NotNil(t, err)
}

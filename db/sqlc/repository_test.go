package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	username := "test"
	password := "test"
	role := "user"
	email := "test@example.com"
	ctx := context.Background()
	repo := NewRepository(testConn)
	err := repo.CreateUserWithRole(ctx, CreateUserParams{
		Username: username,
		Password: password,
		Email:    email,
	}, role)
	require.Nil(t, err)

	user, err := repo.GetUser(ctx, username)
	require.Nil(t, err)
	require.Equal(t, username, user.Username)
	require.Equal(t, password, user.Password)
	require.Equal(t, email, user.Email)
	require.True(t, user.RoleDetails.Valid)
	require.Equal(t, role, user.RoleDetails.String)

	err = repo.RemoveUser(ctx, username)
	require.Nil(t, err)

	_, err = repo.GetUser(ctx, username)
	require.NotNil(t, err)
}

func TestRepositoryFailedAdd(t *testing.T) {
	username := "test2"
	password := "test2"
	role := "dummy"
	email := "test2@example.com"
	ctx := context.Background()

	repo := NewRepository(testConn)
	err := repo.CreateUserWithRole(ctx, CreateUserParams{
		Username: username,
		Password: password,
		Email:    email,
	}, role)
	require.NotNil(t, err)

	_, err = repo.GetUser(ctx, username)
	require.NotNil(t, err)

	err = repo.RemoveUser(ctx, username)
	require.NotNil(t, err)
}

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUserUser(t *testing.T) {
	userName := "test1"
	password := "pass1"
	email := "test1@example.com"
	user := CreateUserParams{
		Username: userName,
		Password: password,
		Email:    email,
	}
	ctx := context.Background()
	result, err := testQueries.CreateUser(ctx, user)
	require.Nil(t, err)
	rowAffected, err := result.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, rowAffected, int64(1))

	_, err = testQueries.CreateUser(ctx, user)
	require.NotNil(t, err)

	result, err = testQueries.DeleteUser(ctx, userName)
	require.Nil(t, err)
	rowAffected, err = result.RowsAffected()
	require.Nil(t, err)
	require.Equal(t, rowAffected, int64(1))

}

package auth

import (
	"authentication/models"
	"authentication/service/token"
	"context"
)

type VerifyAuth interface {
	Verify(context.Context, *models.UserDetail) (token.ClaimMap, error)
}

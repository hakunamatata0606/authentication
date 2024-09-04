package auth

import (
	db "authentication/db/sqlc"
	"authentication/models"
	"authentication/service/password"
	"authentication/service/token"
	"context"
	"fmt"
	"strings"
)

type MysqlAuth struct {
	repo      db.RepositoryInterface
	pwdVerify password.PasswordManager
}

func NewMysqlAuth(repo db.RepositoryInterface, pwdVerify password.PasswordManager) VerifyAuth {
	return &MysqlAuth{
		repo:      repo,
		pwdVerify: pwdVerify,
	}
}

func (auth *MysqlAuth) Verify(ctx context.Context, userDetail *models.UserDetail) (token.ClaimMap, error) {
	user, err := auth.repo.GetUser(ctx, userDetail.Username)
	if err != nil {
		return nil, err
	}
	if !auth.pwdVerify.VerifyPassword(userDetail.Password, user.Password) {
		return nil, fmt.Errorf("password not macthed")
	}
	var roles []string
	if user.RoleDetails.Valid {
		roles = strings.Split(user.RoleDetails.String, ",")
	}
	claims := token.ClaimMap{
		"user": user.Username,
		"role": roles,
	}
	return claims, err
}

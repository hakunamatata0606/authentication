package auth

import (
	db "authentication/db/sqlc"
	"authentication/models"
	"authentication/service/password"
	"authentication/service/token"
	"context"
	"database/sql"
	"fmt"
)

type MysqlAuth struct {
	repo      db.RepositoryInterface
	pwdVerify password.PasswordManager
}

func NewMysqlAuth(dbP *sql.DB, pwdVerify password.PasswordManager) VerifyAuth {
	return &MysqlAuth{
		repo:      db.NewRepository(dbP),
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
	claims := token.ClaimMap{
		"user": user.Username,
		"role": user.RoleDetail,
	}
	return claims, err
}

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type RepositoryInterface interface {
	CreateUserWithRole(context.Context, CreateUserParams, string) error
	RemoveUser(context.Context, string) error
	GetUser(context.Context, string) (GetUserRoleRow, error)
}

type Repository struct {
	Queries
	db *sql.DB
}

func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{
		Queries: *New(db),
		db:      db,
	}
}

func (repo *Repository) execWithTx(ctx context.Context, f func(*Queries) error) error {
	fmt.Println(repo.db)
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := repo.WithTx(tx)
	err = f(query)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (repo *Repository) CreateUserWithRole(ctx context.Context, user CreateUserParams, role string) error {
	err := repo.execWithTx(ctx, func(q *Queries) error {
		_, err := q.CreateUser(ctx, user)
		if err != nil {
			return err
		}
		createdUser, err := q.GetUserForUpdate(ctx, user.Username)
		if err != nil {
			return err
		}
		roleDetail, err := q.GetRoleIdByDetail(ctx, role)
		if err != nil {
			return err
		}
		_, err = q.AddUserRole(ctx, AddUserRoleParams{
			UserID: createdUser.ID.Int32,
			RoleID: roleDetail.ID.Int32,
		})
		return err
	})
	return err
}

func (repo *Repository) RemoveUser(ctx context.Context, username string) error {
	err := repo.execWithTx(ctx, func(q *Queries) error {
		user, err := repo.GetUserForUpdate(ctx, username)
		if err != nil {
			return err
		}
		_, err = repo.DeleteRole(ctx, user.ID.Int32)
		if err != nil {
			return err
		}
		_, err = q.DeleteUser(ctx, username)
		return err
	})
	return err
}

func (repo *Repository) GetUser(ctx context.Context, username string) (GetUserRoleRow, error) {
	return repo.GetUserRole(ctx, username)
}

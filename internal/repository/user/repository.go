package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Kosfedev/auth/internal/model"
	"github.com/Kosfedev/auth/internal/repository"
	"github.com/Kosfedev/auth/internal/repository/user/converter"
	modelRepo "github.com/Kosfedev/auth/internal/repository/user/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	tableAuth       = "auth"
	columnID        = "id"
	columnName      = "name"
	columnEmail     = "email"
	columnRole      = "role"
	columnPassword  = "password"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"
)

type repo struct {
	con *pgx.Conn
}

func NewRepository(con *pgx.Conn) repository.UserRepository {
	return &repo{
		con: con,
	}
}

func (r *repo) Create(ctx context.Context, userData *model.NewUserData) (int64, error) {
	passHash, err := hashPassword(userData.Password)
	if err != nil {
		return 0, err
	}

	builderInsert := sq.Insert(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Columns(columnName, columnEmail, columnRole, columnPassword).
		Values(userData.Name, userData.Email, userData.Role, passHash).
		Suffix(fmt.Sprintf("RETURNING %v", columnID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	err = r.con.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.UserData, error) {
	builderSelect := sq.Select(columnID, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt).
		From(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{columnID: id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var user = &modelRepo.UserData{}
	err = r.con.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.UserDataFromRepo(user), nil
}

func (r *repo) Patch(ctx context.Context, userData *model.UpdatedUserData, id int64) (*model.UserData, error) {
	builderUpdate := sq.Update(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Set(columnUpdatedAt, time.Now()).
		Where(sq.Eq{columnID: id}).
		Suffix(fmt.Sprintf("RETURNING %v, %v, %v, %v, %v, %v", columnID, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt))

	if userData.Name != nil {
		builderUpdate = builderUpdate.Set(columnName, *userData.Name)
	}
	if userData.Email != nil {
		builderUpdate = builderUpdate.Set(columnEmail, *userData.Email)
	}
	if userData.Role != nil {
		builderUpdate = builderUpdate.Set(columnRole, *userData.Role)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	var user = &modelRepo.UserData{}
	err = r.con.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.UserDataFromRepo(user), nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{columnID: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	_, err = r.con.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// TODO: relocate
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

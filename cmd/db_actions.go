package main

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func createUser(ctx context.Context, user *NewUserData) (int64, error) {
	passHash, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	builderInsert := sq.Insert(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Columns(columnName, columnEmail, columnRole, columnPassword).
		Values(user.Name, user.Email, user.Role, passHash).
		Suffix(fmt.Sprintf("RETURNING %v", columnID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	err = con.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func getUser(ctx context.Context, id int64) (*UserData, error) {
	builderSelect := sq.Select(columnID, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt).
		From(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{columnID: id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var user = &UserData{}
	err = con.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func updateUser(ctx context.Context, user *UpdateUserData, id int64) (*UserData, error) {
	builderUpdate := sq.Update(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Set(columnUpdatedAt, time.Now()).
		Where(sq.Eq{columnID: id}).
		Suffix(fmt.Sprintf("RETURNING %v, %v, %v, %v, %v, %v", columnID, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt))

	if user.Name != nil {
		builderUpdate = builderUpdate.Set(columnName, *user.Name)
	}
	if user.Email != nil {
		builderUpdate = builderUpdate.Set(columnEmail, *user.Email)
	}
	if user.Role != nil {
		builderUpdate = builderUpdate.Set(columnRole, *user.Role)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	var updatedUser = &UserData{}
	err = con.QueryRow(ctx, query, args...).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Role, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func deleteUser(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{columnID: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	_, err = con.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

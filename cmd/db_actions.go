package main

import (
	"fmt"
	"reflect"
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

func createUser(user *NewUserData) (int, error) {
	builderInsert := sq.Insert(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Columns(columnName, columnEmail, columnRole, columnPassword).
		Values(user.Name, user.Email, user.Role, user.Password).
		Suffix(fmt.Sprintf("RETURNING %v", columnID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int
	err = con.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func getUser(id int) (*UserData, error) {
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

func updateUser(user *UpdateUserData, id int) (*UserData, error) {
	builderUpdate := sq.Update(tableAuth).
		PlaceholderFormat(sq.Dollar).
		Set(columnUpdatedAt, time.Now()).
		Where(sq.Eq{columnID: id}).
		Suffix(fmt.Sprintf("RETURNING %v, %v, %v, %v, %v, %v", columnID, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt))

	values := reflect.ValueOf(*user)
	types := values.Type()

	for i := 0; i < values.NumField(); i++ {
		if value := values.Field(i); !value.IsNil() {
			builderUpdate = builderUpdate.Set(types.Field(i).Name, value.Interface())
		}
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

func deleteUser(id int) error {
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

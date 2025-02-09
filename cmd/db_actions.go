package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func createUser(user *NewUserData) (int, error) {
	pgDSN, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return 0, errors.New("PG_DSN environment variable not set")
	}

	ctx := context.Background()
	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return 0, err
	}
	defer closeConOrLog(ctx, con)

	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role", "password").
		Values(user.Name, user.Email, user.Role, user.Password).
		Suffix("RETURNING id")

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
	pgDSN, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return nil, errors.New("PG_DSN environment variable not set")
	}

	ctx := context.Background()
	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return nil, err
	}
	defer closeConOrLog(ctx, con)

	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

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
	pgDSN, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return nil, errors.New("PG_DSN environment variable not set")
	}

	ctx := context.Background()
	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return nil, err
	}
	defer closeConOrLog(ctx, con)

	builderUpdate := sq.Update("auth").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id})

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

	_, err = con.Query(ctx, query, args...)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}

	updatedUser, err := getUser(id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func deleteUser(id int) error {
	pgDSN, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return errors.New("PG_DSN environment variable not set")
	}

	ctx := context.Background()
	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return err
	}
	defer closeConOrLog(ctx, con)

	builderDelete := sq.Delete("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	_, err = con.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	return nil
}

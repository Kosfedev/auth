package main

import (
	"context"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func closeConOrLog(ctx context.Context, con *pgx.Conn) {
	if err := con.Close(ctx); err != nil {
		log.Printf("Error while closing connection: %+v\n", err)
	}
}

func parseID(idStr string) (int, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

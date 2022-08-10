package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func Open(user string, password string, host string, port string, dbname string) *pgx.Conn {
	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}

	return conn
}

func Close(conn *pgx.Conn) {
	conn.Close(context.Background())
}

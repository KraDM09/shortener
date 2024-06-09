package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/KraDM09/shortener/internal/app/config"
	"github.com/KraDM09/shortener/internal/app/util"
	"github.com/jackc/pgx/v5"
)

type Database struct{}

type Record struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short"`
	OriginalURL string `json:"original_url" db:"original"`
}

func (s Database) Save(hash string, url string) {
	conn, err := pgx.Connect(context.Background(), config.FlagDatabaseDsn)
	if err != nil {
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(),
		"INSERT INTO shortener.urls (uuid, original, short)"+
			"VALUES ($1, $2, $3)",
		util.CreateUUID(),
		url,
		hash,
	)
	if err != nil {
		return
	}
}

func (s Database) Get(hash string) string {
	conn, err := pgx.Connect(context.Background(), config.FlagDatabaseDsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(),
		"SELECT original FROM shortener.urls WHERE short = $1",
		hash,
	)

	defer rows.Close()

	records := make([]Record, 0, 1)

	for rows.Next() {
		var r Record
		err = rows.Scan(&r.OriginalURL)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, r)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return records[0].OriginalURL
}

func (s Database) Migrate() error {
	conn, err := pgx.Connect(context.Background(), config.FlagDatabaseDsn)
	if err != nil {
		return fmt.Errorf("cannot connect to db. %w", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(
		context.Background(),
		"CREATE schema IF NOT EXISTS shortener;"+
			"CREATE TABLE IF NOT EXISTS shortener.urls ("+
			"uuid UUID PRIMARY KEY,"+
			"original TEXT NOT NULL,"+
			"short TEXT NOT NULL);")
	if err != nil {
		return nil
	}

	return err
}

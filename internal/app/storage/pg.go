package storage

import (
	"context"

	"github.com/KraDM09/shortener/internal/app/util"
	"github.com/jackc/pgx/v5"
)

type PG struct {
	conn *pgx.Conn
}

// NewStore возвращает новый экземпляр PostgreSQL-хранилища
func (pg PG) NewStore(conn *pgx.Conn) *PG {
	return &PG{conn: conn}
}

func (pg PG) Save(hash string, url string, userID string) (string, error) {
	row, err := pg.conn.Exec(context.Background(),
		"INSERT INTO shortener.urls (uuid, original, short, user_id)"+
			"VALUES ($1, $2, $3, $4) ON CONFLICT (original) DO NOTHING RETURNING *",
		util.CreateUUID(),
		url,
		hash,
		userID,
	)
	if err != nil {
		return "", err
	}

	if row.RowsAffected() == 0 {
		short, err := pg.GetHashByOriginal(url)
		if err != nil {
			return "", err
		}

		return short, ErrConflict
	}

	return hash, nil
}

func (pg PG) SaveBatch(batch []URL, userID string) error {
	_, err := pg.conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"shortener", "urls"},
		[]string{"uuid", "original", "short", "user_id"},
		pgx.CopyFromSlice(len(batch), func(i int) ([]any, error) {
			return []any{
				util.CreateUUID(),
				batch[i].Original,
				batch[i].Short,
				userID,
			}, nil
		}),
	)
	if err != nil {
		return err
	}

	return nil
}

func (pg PG) GetHashByOriginal(original string) (string, error) {
	rows, err := pg.conn.Query(context.Background(),
		"SELECT short FROM shortener.urls WHERE original = $1",
		original,
	)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	records := make([]URL, 0, 1)

	for rows.Next() {
		var r URL
		err = rows.Scan(&r.Short)
		if err != nil {
			return "", err
		}
		records = append(records, r)
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}

	return records[0].Short, nil
}

func (pg PG) Get(hash string) (string, error) {
	rows, err := pg.conn.Query(context.Background(),
		"SELECT original FROM shortener.urls WHERE short = $1",
		hash,
	)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	records := make([]URL, 0, 1)

	for rows.Next() {
		var r URL
		err = rows.Scan(&r.Original)
		if err != nil {
			return "", err
		}
		records = append(records, r)
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}

	return records[0].Original, nil
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (pg PG) Bootstrap(ctx context.Context) error {
	tx, err := pg.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
        CREATE schema IF NOT EXISTS shortener
    `)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS shortener.urls (
            uuid UUID PRIMARY KEY,
            original TEXT NOT NULL UNIQUE,
            short TEXT NOT NULL,
            user_id UUID NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (pg PG) GetUrlsByUserID(userID string) (*[]URL, error) {
	rows, err := pg.conn.Query(context.Background(),
		`SELECT short, original
			 FROM shortener.urls
			 WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	urls := make([]URL, 0, 1)

	for rows.Next() {
		var r URL
		err = rows.Scan(&r.Short, &r.Original)
		if err != nil {
			return nil, err
		}
		urls = append(urls, r)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &urls, nil
}

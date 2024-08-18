package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KraDM09/shortener/internal/app/util"
	"github.com/jackc/pgx/v5"
)

type PG struct {
	pool *pgxpool.Pool
}

// NewStore возвращает новый экземпляр PostgreSQL-хранилища
func (pg PG) NewStore(pool *pgxpool.Pool) *PG {
	return &PG{pool: pool}
}

func (pg PG) Save(hash string, url string, userID string) (string, error) {
	row, err := pg.pool.Exec(context.Background(),
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
	_, err := pg.pool.CopyFrom(
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
	rows, err := pg.pool.Query(context.Background(),
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

func (pg PG) Get(hash string) (*URL, error) {
	rows, err := pg.pool.Query(context.Background(),
		"SELECT original, is_deleted FROM shortener.urls WHERE short = $1",
		hash,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	records := make([]URL, 0, 1)

	for rows.Next() {
		var r URL
		err = rows.Scan(&r.Original, &r.IsDeleted)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &records[0], nil
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (pg PG) Bootstrap(ctx context.Context) error {
	tx, err := pg.pool.Begin(ctx)
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
            user_id UUID NOT NULL,
            is_deleted BOOL DEFAULT FALSE
        )
    `)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (pg PG) GetUrlsByUserID(userID string) (*[]URL, error) {
	rows, err := pg.pool.Query(context.Background(),
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

func (pg PG) DeleteUrls(ctx context.Context, deleteHashes ...DeleteHash) error {
	// соберём данные для создания запроса с групповым обновлением
	var values []string
	for _, hash := range deleteHashes {
		// в нашем запросе по 2 параметра на каждое сообщение
		params := fmt.Sprintf("('%s', '%s')", hash.UserID, hash.Short)
		values = append(values, params)
	}

	// составляем строку запроса
	query := `UPDATE shortener.urls
	SET is_deleted = TRUE
	FROM (VALUES ` + strings.Join(values, ",") + `) AS new_values (user_id, short)
	WHERE urls.user_id = new_values.user_id::UUID
	  AND urls.short = new_values.short
	  AND urls.is_deleted = FALSE;`

	// добавляем новые сообщения в БД
	_, err := pg.pool.Exec(ctx, query)

	return err
}

func (pg PG) GetQuantityUserShortUrls(
	userID string,
	shortUrls *[]string,
) (int, error) {
	var values []string
	for _, short := range *shortUrls {
		params := fmt.Sprintf("'%s'", short)
		values = append(values, params)
	}

	rows, err := pg.pool.Query(context.Background(),
		`SELECT count(short) AS quantity
			FROM shortener.urls
			WHERE user_id = $1::UUID
			  AND is_deleted = FALSE
			  AND short IN (`+strings.Join(values, ",")+`);`,
		userID,
	)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	quantity := 0

	for rows.Next() {
		err = rows.Scan(&quantity)
		if err != nil {
			return 0, err
		}
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

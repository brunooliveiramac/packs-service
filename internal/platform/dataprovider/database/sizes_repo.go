package database

import (
	"context"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SizeRow struct {
	ID     int64
	Size   int
	Active bool
}

type SizesRepository struct{ pool *pgxpool.Pool }

func NewSizesRepository(pool *pgxpool.Pool) *SizesRepository { return &SizesRepository{pool: pool} }

func (r *SizesRepository) List(ctx context.Context, includeInactive bool) ([]SizeRow, error) {
	q := `SELECT id, size, active FROM pack_sizes`
	if !includeInactive {
		q += ` WHERE active = true`
	}
	q += ` ORDER BY size ASC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SizeRow
	for rows.Next() {
		var sr SizeRow
		if err := rows.Scan(&sr.ID, &sr.Size, &sr.Active); err != nil {
			return nil, err
		}
		out = append(out, sr)
	}
	return out, rows.Err()
}

func (r *SizesRepository) Upsert(ctx context.Context, size int, active bool) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO pack_sizes (size, active, created_at, updated_at)
		VALUES ($1, $2, now(), now())
		ON CONFLICT (size) DO UPDATE SET active = EXCLUDED.active, updated_at = now()
	`, size, active)
	return err
}

func (r *SizesRepository) SetActive(ctx context.Context, size int, active bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE pack_sizes SET active = $2, updated_at = now() WHERE size = $1`, size, active)
	return err
}

func (r *SizesRepository) Delete(ctx context.Context, size int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM pack_sizes WHERE size = $1`, size)
	return err
}

// Helper to return active sizes as sorted ints.
func (r *SizesRepository) ActiveSizes(ctx context.Context) ([]int, error) {
	rows, err := r.List(ctx, false)
	if err != nil {
		return nil, err
	}
	s := make([]int, 0, len(rows))
	for _, r := range rows {
		s = append(s, r.Size)
	}
	sort.Ints(s)
	return s, nil
}

var _ = time.Now



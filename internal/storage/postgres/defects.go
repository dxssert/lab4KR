package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/dxssert/lab4KR/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DefectRepository struct {
	pool *pgxpool.Pool
}

func NewDefectRepository(pool *pgxpool.Pool) *DefectRepository {
	return &DefectRepository{pool: pool}
}

func (repo *DefectRepository) List(ctx context.Context) ([]model.Defect, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := repo.pool.Query(ctx, "SELECT * FROM public.get_defects()")
	if err != nil {
		return nil, fmt.Errorf("list defects: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.Defect])
}

func (repo *DefectRepository) Get(ctx context.Context, id int64) (*model.Defect, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := repo.pool.Query(ctx, "SELECT * FROM public.get_defect($1)", id)
	if err != nil {
		return nil, fmt.Errorf("get defect %d: %w", id, err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.Defect])
	if err != nil {
		return nil, fmt.Errorf("get defect %d: %w", id, err)
	}
	return &item, nil
}

func (repo *DefectRepository) Create(ctx context.Context,
	productionOutputID int64, reasonID int64, defectGrams int64, createdBy int64,
	comment string) (*model.Defect, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := repo.pool.Query(ctx, "SELECT * FROM public.create_defect($1,$2,$3,$4,$5)",
		productionOutputID, reasonID, defectGrams, createdBy, comment)

	if err != nil {
		return nil, fmt.Errorf("create defect %d: %w", productionOutputID, err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.Defect])
	if err != nil {
		return nil, fmt.Errorf("create defect %d: %w", productionOutputID, err)
	}

	return &item, nil
}
func (repo *DefectRepository) Update(
	ctx context.Context,
	id int64,
	productionOutputID int64,
	reasonID int64,
	defectGrams int64,
	comment string,
) (*model.Defect, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := repo.pool.Query(
		ctx,
		"SELECT * FROM public.update_defect($1, $2, $3, $4, $5)",
		id, productionOutputID, reasonID, defectGrams, comment,
	)
	if err != nil {
		return nil, fmt.Errorf("update defect %d: %w", id, err)
	}

	item, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByPos[model.Defect],
	)
	if err != nil {
		return nil, fmt.Errorf("update defect %d: %w", id, err)
	}

	return &item, nil
}

func (repo *DefectRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var deletedID int64

	err := repo.pool.QueryRow(
		ctx,
		"SELECT * FROM public.delete_defect($1)",
		id,
	).Scan(&deletedID)
	if err != nil {
		return fmt.Errorf("delete defect %d: %w", id, err)
	}

	return nil
}

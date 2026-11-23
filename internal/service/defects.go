package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dxssert/lab4KR/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalid  = errors.New("некорректные данные")
	ErrNotFound = errors.New("запись брака не найдена")
)

type DefectData struct {
	ProductionOutputID int64
	ReasonID           int64
	DefectGrams        int64
	Comment            string
}

type DefectRepository interface {
	List(context.Context) ([]model.Defect, error)
	Get(context.Context, int64) (*model.Defect, error)
	Create(context.Context, int64, int64, int64, int64, string) (*model.Defect, error)
	Update(context.Context, int64, int64, int64, int64, string) (*model.Defect, error)
	Delete(context.Context, int64) error
}

type Defects struct {
	repo DefectRepository
}

func NewDefects(repo DefectRepository) *Defects {
	return &Defects{repo: repo}
}

func (s *Defects) List(ctx context.Context) ([]model.Defect, error) {
	return s.repo.List(ctx)
}

func (s *Defects) Get(ctx context.Context, id int64) (*model.Defect, error) {
	item, err := s.repo.Get(ctx, id)
	return item, classify(err)
}

func (s *Defects) Create(ctx context.Context, data DefectData, authorID int64) (*model.Defect, error) {
	if err := data.validate(); err != nil {
		return nil, err
	}
	if authorID <= 0 {
		return nil, fmt.Errorf("%w: created_by должен быть положительным", ErrInvalid)
	}
	item, err := s.repo.Create(ctx, data.ProductionOutputID, data.ReasonID, data.DefectGrams, authorID, strings.TrimSpace(data.Comment))
	return item, classify(err)
}

func (s *Defects) Update(ctx context.Context, id int64, data DefectData) (*model.Defect, error) {
	if err := data.validate(); err != nil {
		return nil, err
	}
	item, err := s.repo.Update(ctx, id, data.ProductionOutputID, data.ReasonID, data.DefectGrams, strings.TrimSpace(data.Comment))
	return item, classify(err)
}

func (s *Defects) Delete(ctx context.Context, id int64) error {
	return classify(s.repo.Delete(ctx, id))
}

func (d DefectData) validate() error {
	if d.ProductionOutputID <= 0 || d.ReasonID <= 0 || d.DefectGrams <= 0 {
		return fmt.Errorf("%w: production_output_id, reason_id и defect_grams должны быть положительными", ErrInvalid)
	}
	return nil
}

func classify(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var dbErr *pgconn.PgError
	if errors.As(err, &dbErr) {
		switch dbErr.Code {
		case "23503":
			return fmt.Errorf("%w: выпуск, причина брака или автор не найдены", ErrInvalid)
		case "23514", "23502":
			return fmt.Errorf("%w: нарушены ограничения записи", ErrInvalid)
		}
	}
	return err
}

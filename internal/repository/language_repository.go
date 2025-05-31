package repository

import (
	"context"
	"errors"
	"shadowify/internal/model"

	"github.com/jmoiron/sqlx"
)

type LanguageRepository struct {
	db *sqlx.DB
}

func NewLanguageRepository(db *sqlx.DB) *LanguageRepository {
	return &LanguageRepository{
		db: db,
	}
}

func (r *LanguageRepository) Create(ctx context.Context, language *model.Language) error {
	query := `
		INSERT INTO languages (code, name, flag_url)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		language.Code,
		language.Name,
		language.FlagURL,
	).Scan(&language.Id, &language.CreatedAt, &language.UpdatedAt)
}

func (r *LanguageRepository) GetByID(ctx context.Context, id string) (*model.Language, error) {
	var language model.Language
	query := `SELECT * FROM languages WHERE id = $1`
	err := r.db.GetContext(ctx, &language, query, id)
	if err != nil {
		return nil, err
	}
	return &language, nil
}

func (r *LanguageRepository) GetAll(ctx context.Context) ([]*model.Language, error) {
	var languages []*model.Language
	query := `SELECT * FROM languages ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &languages, query)
	if err != nil {
		return nil, err
	}
	return languages, nil
}

func (r *LanguageRepository) Update(ctx context.Context, language *model.Language) error {
	query := `
		UPDATE languages 
		SET code = $1, name = $2, flag_url = $3
		WHERE id = $4 RETURNING updated_at
	`
	result := r.db.QueryRowContext(
		ctx,
		query,
		language.Code,
		language.Name,
		language.FlagURL,
		language.Id,
	)

	err := result.Scan(&language.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *LanguageRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM languages WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("language not found")
	}

	return nil
}

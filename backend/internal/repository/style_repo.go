package repository

import (
	"web-gallery/internal/model"

	"github.com/jmoiron/sqlx"
)

type StyleRepository struct {
	db *sqlx.DB
}

func NewStyleRepository(db *sqlx.DB) *StyleRepository {
	return &StyleRepository{db: db}
}

func (r *StyleRepository) GetAll() ([]model.Style, error) {
	styles := []model.Style{}
	err := r.db.Select(&styles, `SELECT id, name FROM styles ORDER BY name`)
	return styles, err
}

package repository

import (
	"web-gallery/internal/model"

	"github.com/jmoiron/sqlx"
)

type MaterialRepository struct {
	db *sqlx.DB
}

func NewMaterialRepository(db *sqlx.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) GetAll() ([]model.Material, error) {
	materials := []model.Material{}

	err := r.db.Select(&materials, `SELECT id, name FROM materials`)
	return materials, err
}

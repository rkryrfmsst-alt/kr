package repository

import (
	"web-gallery/internal/model"

	"github.com/jmoiron/sqlx"
)

type PlotRepository struct {
	db *sqlx.DB
}

func NewPlotRepository(db *sqlx.DB) *PlotRepository {
	return &PlotRepository{db: db}
}

func (r *PlotRepository) GetAll() ([]model.Plot, error) {
	plots := []model.Plot{}
	err := r.db.Select(&plots, `SELECT id, name FROM plots ORDER BY name`)
	return plots, err
}

package repository

import (
	"database/sql"
	"errors"
	"web-gallery/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PaintingRepository struct {
	db *sqlx.DB
}

func NewPaintingRepository(db *sqlx.DB) *PaintingRepository {
	return &PaintingRepository{db: db}
}

func (r *PaintingRepository) GetAll() ([]model.Painting, error) {
	var paintings []model.Painting
	if err := r.db.Select(&paintings, `SELECT id, title, description, year, image_path FROM paintings`); err != nil {
		return nil, err
	}
	if len(paintings) == 0 {
		return []model.Painting{}, nil
	}

	ids := make([]int, len(paintings))
	idx := make(map[int]*model.Painting, len(paintings))
	for i := range paintings {
		paintings[i].Authors = []model.Author{}
		paintings[i].Styles = []model.Style{}
		paintings[i].Plots = []model.Plot{}
		ids[i] = paintings[i].ID
		idx[paintings[i].ID] = &paintings[i]
	}

	type matRow struct {
		PaintingID int `db:"painting_id"`
		model.Material
	}
	var matRows []matRow
	if err := r.db.Select(&matRows, `
		SELECT p.id AS painting_id, m.id, m.name
		FROM materials m
		JOIN paintings p ON p.material_id = m.id
		WHERE p.id = ANY($1)
	`, pq.Array(ids)); err != nil {
		return nil, err
	}
	for _, row := range matRows {
		m := row.Material
		idx[row.PaintingID].Material = &m
	}

	type authRow struct {
		PaintingID int `db:"painting_id"`
		model.Author
	}
	var authRows []authRow
	if err := r.db.Select(&authRows, `
		SELECT pa.painting_id, a.id, a.first_name, a.last_name, a.middle_name, a.description
		FROM authors a
		JOIN paintings_authors pa ON pa.author_id = a.id
		WHERE pa.painting_id = ANY($1)
	`, pq.Array(ids)); err != nil {
		return nil, err
	}
	for _, row := range authRows {
		idx[row.PaintingID].Authors = append(idx[row.PaintingID].Authors, row.Author)
	}

	type styleRow struct {
		PaintingID int `db:"painting_id"`
		model.Style
	}
	var styleRows []styleRow
	if err := r.db.Select(&styleRows, `
		SELECT ps.painting_id, s.id, s.name
		FROM styles s
		JOIN paintings_styles ps ON ps.style_id = s.id
		WHERE ps.painting_id = ANY($1)
	`, pq.Array(ids)); err != nil {
		return nil, err
	}
	for _, row := range styleRows {
		idx[row.PaintingID].Styles = append(idx[row.PaintingID].Styles, row.Style)
	}

	type plotRow struct {
		PaintingID int `db:"painting_id"`
		model.Plot
	}
	var plotRows []plotRow
	if err := r.db.Select(&plotRows, `
		SELECT pp.painting_id, pl.id, pl.name
		FROM plots pl
		JOIN paintings_plots pp ON pp.plot_id = pl.id
		WHERE pp.painting_id = ANY($1)
	`, pq.Array(ids)); err != nil {
		return nil, err
	}
	for _, row := range plotRows {
		idx[row.PaintingID].Plots = append(idx[row.PaintingID].Plots, row.Plot)
	}

	return paintings, nil
}

func (r *PaintingRepository) GetByID(id int) (*model.Painting, error) {
	p := &model.Painting{}
	if err := r.db.Get(p, `SELECT id, title, description, year, image_path FROM paintings WHERE id = $1`, id); err != nil {
		return nil, err
	}
	return p, r.fillRelations(p)
}

func (r *PaintingRepository) fillRelations(p *model.Painting) error {
	mat := model.Material{}
	err := r.db.Get(&mat, `
		SELECT m.id, m.name
		FROM materials m
		JOIN paintings ON paintings.material_id = m.id
		WHERE paintings.id = $1
	`, p.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err == nil {
		p.Material = &mat
	}

	if err := r.db.Select(&p.Authors, `
		SELECT a.id, a.first_name, a.last_name, a.middle_name, a.description
		FROM authors a
		JOIN paintings_authors pa ON pa.author_id = a.id
		WHERE pa.painting_id = $1
	`, p.ID); err != nil {
		return err
	}

	if err := r.db.Select(&p.Styles, `
		SELECT s.id, s.name
		FROM styles s
		JOIN paintings_styles ps ON ps.style_id = s.id
		WHERE ps.painting_id = $1
	`, p.ID); err != nil {
		return err
	}

	if err := r.db.Select(&p.Plots, `
		SELECT pl.id, pl.name
		FROM plots pl
		JOIN paintings_plots pp ON pp.plot_id = pl.id
		WHERE pp.painting_id = $1
	`, p.ID); err != nil {
		return err
	}

	return nil
}

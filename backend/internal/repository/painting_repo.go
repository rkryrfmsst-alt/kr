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

func (r *PaintingRepository) Create(title string, description *string, year *int, materialID *int, imagePath *string, authorIDs, styleIDs, plotIDs []int) (*model.Painting, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO paintings (title, description, year, material_id, image_path)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, description, year, materialID, imagePath,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	for _, aID := range authorIDs {
		if _, err := r.db.Exec(`INSERT INTO paintings_authors (painting_id, author_id) VALUES ($1, $2)`, id, aID); err != nil {
			return nil, err
		}
	}
	for _, sID := range styleIDs {
		if _, err := r.db.Exec(`INSERT INTO paintings_styles (painting_id, style_id) VALUES ($1, $2)`, id, sID); err != nil {
			return nil, err
		}
	}
	for _, pID := range plotIDs {
		if _, err := r.db.Exec(`INSERT INTO paintings_plots (painting_id, plot_id) VALUES ($1, $2)`, id, pID); err != nil {
			return nil, err
		}
	}
	return r.GetByID(id)
}

func (r *PaintingRepository) Update(id int, title string, description *string, year *int, materialID *int, imagePath *string, authorIDs, styleIDs, plotIDs []int) (*model.Painting, error) {
	if imagePath != nil {
		if _, err := r.db.Exec(
			`UPDATE paintings SET title=$1, description=$2, year=$3, material_id=$4, image_path=$5 WHERE id=$6`,
			title, description, year, materialID, imagePath, id,
		); err != nil {
			return nil, err
		}
	} else {
		if _, err := r.db.Exec(
			`UPDATE paintings SET title=$1, description=$2, year=$3, material_id=$4 WHERE id=$5`,
			title, description, year, materialID, id,
		); err != nil {
			return nil, err
		}
	}

	r.db.Exec(`DELETE FROM paintings_authors WHERE painting_id=$1`, id)
	r.db.Exec(`DELETE FROM paintings_styles  WHERE painting_id=$1`, id)
	r.db.Exec(`DELETE FROM paintings_plots   WHERE painting_id=$1`, id)

	for _, aID := range authorIDs {
		if _, err := r.db.Exec(`INSERT INTO paintings_authors (painting_id, author_id) VALUES ($1,$2)`, id, aID); err != nil {
			return nil, err
		}
	}
	for _, sID := range styleIDs {
		if _, err := r.db.Exec(`INSERT INTO paintings_styles (painting_id, style_id) VALUES ($1,$2)`, id, sID); err != nil {
			return nil, err
		}
	}
	for _, pID := range plotIDs {
		if _, err := r.db.Exec(`INSERT INTO paintings_plots (painting_id, plot_id) VALUES ($1,$2)`, id, pID); err != nil {
			return nil, err
		}
	}
	return r.GetByID(id)
}

func (r *PaintingRepository) Delete(id int) error {
	r.db.Exec(`DELETE FROM paintings_authors WHERE painting_id=$1`, id)
	r.db.Exec(`DELETE FROM paintings_styles  WHERE painting_id=$1`, id)
	r.db.Exec(`DELETE FROM paintings_plots   WHERE painting_id=$1`, id)
	_, err := r.db.Exec(`DELETE FROM paintings WHERE id=$1`, id)
	return err
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

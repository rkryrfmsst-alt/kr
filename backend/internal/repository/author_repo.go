package repository

import (
	"web-gallery/internal/model"

	"github.com/jmoiron/sqlx"
)

type AuthorRepository struct {
	db *sqlx.DB
}

func NewAuthorRepository(db *sqlx.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) GetAll() ([]model.Author, error) {
	authors := []model.Author{}
	err := r.db.Select(&authors, `SELECT id, first_name, last_name, middle_name, description FROM authors`)
	return authors, err
}

func (r *AuthorRepository) GetByID(id int) (*model.Author, error) {
	author := &model.Author{}
	err := r.db.Get(author, `SELECT id, first_name, last_name, middle_name, description FROM authors WHERE id = $1`, id)
	return author, err
}

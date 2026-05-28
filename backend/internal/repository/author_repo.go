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

func (r *AuthorRepository) Create(firstName, lastName string, middleName, description *string) (*model.Author, error) {
	author := &model.Author{}
	err := r.db.QueryRowx(
		`INSERT INTO authors (first_name, last_name, middle_name, description)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, first_name, last_name, middle_name, description`,
		firstName, lastName, middleName, description,
	).StructScan(author)
	return author, err
}

func (r *AuthorRepository) Update(id int, firstName, lastName string, middleName, description *string) (*model.Author, error) {
	author := &model.Author{}
	err := r.db.QueryRowx(
		`UPDATE authors SET first_name=$1, last_name=$2, middle_name=$3, description=$4
		 WHERE id=$5
		 RETURNING id, first_name, last_name, middle_name, description`,
		firstName, lastName, middleName, description, id,
	).StructScan(author)
	return author, err
}

func (r *AuthorRepository) Delete(id int) error {
	r.db.Exec(`DELETE FROM paintings_authors WHERE author_id=$1`, id)
	_, err := r.db.Exec(`DELETE FROM authors WHERE id=$1`, id)
	return err
}

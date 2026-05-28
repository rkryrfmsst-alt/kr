package repository

import (
	"web-gallery/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.Get(u, `SELECT id, email, username, password, role FROM users WHERE email = $1`, email)
	return u, err
}

func (r *UserRepository) Create(email, username, passwordHash string) (*model.User, error) {
	u := &model.User{}
	err := r.db.Get(u,
		`INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
		 RETURNING id, email, username, role`,
		email, username, passwordHash,
	)
	return u, err
}

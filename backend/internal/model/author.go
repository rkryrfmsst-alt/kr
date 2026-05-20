package model

type Author struct {
	ID          int     `db:"id" json:"id"`
	FirstName   string  `db:"first_name" json:"first_name" binding:"required"`
	LastName    string  `db:"last_name" json:"last_name" binding:"required"`
	MiddleName  *string `db:"middle_name" json:"middle_name"`
	Description *string `db:"description" json:"description"`
}

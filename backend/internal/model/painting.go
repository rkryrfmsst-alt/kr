package model

type Painting struct {
	ID          int       `db:"id"          json:"id"`
	Title       string    `db:"title"        json:"title"`
	Description *string   `db:"description"  json:"description"`
	Year        *int      `db:"year"         json:"year"`
	ImagePath   *string   `db:"image_path"   json:"image_path"`
	Material    *Material `db:"-"            json:"material"`
	Authors     []Author  `db:"-"            json:"authors"`
	Styles      []Style   `db:"-"            json:"styles"`
	Plots       []Plot    `db:"-"            json:"plots"`
}

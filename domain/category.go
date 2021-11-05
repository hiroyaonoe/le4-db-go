package domain

type Category struct {
	CategoryID int    `db:"category_id"`
	Name       string `db:"name"`
}

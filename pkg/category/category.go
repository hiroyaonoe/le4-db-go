package category

type Category struct {
	CategoryID int    `db:"category_id"`
	Name       string `db:"name"`
}

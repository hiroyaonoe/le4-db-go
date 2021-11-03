package domain

type Tag struct {
	TagID    int    `db:"tag_id"`
	Name     string `db:"name"`
	ThreadID int    `db:"thread_id"`
}

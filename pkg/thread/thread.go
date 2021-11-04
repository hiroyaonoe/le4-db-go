package thread

import (
	"github.com/hiroyaonoe/le4-db-go/pkg/datetime"
)

type Thread struct {
	ThreadID     int               `db:"thread_id"`
	Title        string            `db:"title"`
	UserID       int               `db:"user_id"`
	UserName     string            `db:"user_name"`
	CategoryID   int               `db:"category_id"`
	CategoryName string            `db:"category_name"`
	CreatedAt    datetime.DateTime `db:"created_at"`
}

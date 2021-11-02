package comment

import "github.com/hiroyaonoe/le4-db-go/pkg/datetime"

type Comment struct {
	CommentID int               `db:"comment_id"`
	ThreadID  int               `db:"thread_id"`
	Content   string            `db:"content"`
	UserID    int               `db:"user_id"`
	UserName  string            `db:"user_name"`
	CreatedAt datetime.DateTime `db:"created_at"`
}

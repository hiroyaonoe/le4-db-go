package domain

type Comment struct {
	CommentID   int      `db:"comment_id"`
	ThreadID    int      `db:"thread_id"`
	ThreadTitle string   `db:"thread_title"`
	Content     string   `db:"content"`
	UserID      int      `db:"user_id"`
	UserName    string   `db:"user_name"`
	CreatedAt   DateTime `db:"created_at"`
}

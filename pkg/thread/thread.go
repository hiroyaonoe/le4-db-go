package thread

import "time"

type Thread struct {
	ThreadID  int       `db:"thread_id"`
	Title     string    `db:"title"`
	UserID    int       `db:"user_id"`
	UserName  string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
}

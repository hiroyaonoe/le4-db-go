package thread

import "time"

type PostThread struct {
	ThreadID  int       `db:"thread_id"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

package thread

type Thread struct {
	ThreadID int    `db:"thread_id"`
	Title    string `db:"title"`
}

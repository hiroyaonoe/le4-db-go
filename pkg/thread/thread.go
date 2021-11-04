package thread

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Thread struct {
	ThreadID  int       `db:"thread_id"`
	Title     string    `db:"title"`
	UserID    int       `db:"user_id"`
	UserName  string    `db:"user_name"`
	CreatedAt DateTime  `db:"created_at"`
}

type DateTime struct {
	t time.Time
}

func NewDateTime(t time.Time) DateTime {
	return DateTime{t}
}

func (t DateTime) String() string {
	return t.t.Local().Format("2006/01/02 15:04:05")
}

// Scan はデータベースの値をPasswordにマッピングする
func (t *DateTime) Scan(value interface{}) error {
	time, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("Invalid value:%s", value)
	}
	t.t = time
	return nil
}

// Value はPasswordのフィールドのうちデータベースに保存するものを指定する
func (t DateTime) Value() (driver.Value, error) {
	return t.t, nil
}

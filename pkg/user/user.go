package user

type User struct {
	UserID   int      `db:"user_id"`
	Role     string   `db:"role"`
	Name     string   `db:"name"`
	Password Password `db:"password"`
}

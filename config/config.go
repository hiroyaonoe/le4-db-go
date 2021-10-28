package config

import "os"

func DSN() string {
	return "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"
}

func Port() string {
	return ":" + os.Getenv("APP_PORT")
}

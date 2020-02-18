package sqltest

import (
	"os"
)

var TestPostgresConnection = LoadEnv("TEST_POSTGRES_CONNECTION", "postgres://localhost:5432/postgres?user=postgres&sslmode=disable")
var TestMysqlConnection = LoadEnv("TEST_MYSQL_CONNECTION", "root@/mysql")

func LoadEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallback
}

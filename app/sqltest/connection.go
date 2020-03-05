package sqltest

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

var TestPostgresConnection = LoadEnv("TEST_POSTGRES_CONNECTION", "postgres://localhost:5432/postgres?user=postgres&password=test_password&sslmode=disable")
var TestMysqlConnection = LoadEnv("TEST_MYSQL_CONNECTION", "root@/mysql")

func LoadEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallback
}

func GetTestPostgres() *sql.DB {
	s, err := sql.Open("postgres", TestPostgresConnection)
	if err != nil {
		panic(err)
	}
	WaitForDB(s)
	return s
}

func GetTestMysql() *sql.DB {
	s, err := sql.Open("mysql", TestMysqlConnection)
	if err != nil {
		panic(err)
	}
	WaitForDB(s)
	return s
}

// When running tests in docker compose database may not be immediately available
// Sometime it requires few seconds for initialization
func WaitForDB(db *sql.DB) {
	startTime := time.Now()
	for {
		err := db.Ping()
		if err == nil {
			return
		}
		if time.Since(startTime) > 10*time.Second {
			panic(fmt.Sprintf("timeout, couldn't connect to DB: %s", err))
		}
	}
}

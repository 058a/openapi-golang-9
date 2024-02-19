package env

import "os"

func GetDbDriver() string {
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}
	return dbDriver
}

func GetDbDataSourceName() string {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		dbDsn = "host=localhost port=5432 user=user password=password dbname=openapi sslmode=disable"
	}
	return dbDsn
}

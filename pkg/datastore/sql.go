package datastore

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	logger "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

func InitMySQL(conn string) *sqlx.DB {
	connStr := fmt.Sprintf("%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", conn)
	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		logger.Fatal("sqlx.Connect Error: ", err)

	}
	return db
}
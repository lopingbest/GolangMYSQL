package Golang_MYSQL

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestEmpty(t *testing.T) {

}

func TestTestOpenConnection(t *testing.T) {
	db, err := sql.Open("mysql", "root:123@tcp(localhost:3306)/belajar_golang_database")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//gunakan DB
}

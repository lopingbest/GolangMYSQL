package Golang_MYSQL

import (
	"database/sql"
	"time"
)

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:123@tcp(localhost:3306)/belajar_golang_database")
	if err != nil {
		panic(err)
	}

	//value didalam adalah minimal
	db.SetMaxIdleConns(10)
	//value didalam adalah jumlah nilai
	db.SetMaxOpenConns(100)
	//kalo 5 menit enggak ada pergerakan, maka akan di close
	db.SetConnMaxLifetime(5 * time.Minute)
	// kalo udah 60menit, akan dibuatkan koneksi baru
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

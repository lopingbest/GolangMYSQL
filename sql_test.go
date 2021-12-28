package Golang_MYSQL

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES('joko','Joko')"
	//exec tidak akan mengembalikan result
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//selama masih ada data di rows, maka akan terus diambil
	for rows.Next() {
		var id, name string
		//didalam kurung dimasukan parameterdari apa yang akan  kita ambil
		//Pointer dipakai karena kita akan ngeset data dari parameter. Kalo enggak pointer, maka data tidak akan dipakai
		err = rows.Scan(&id, &name)
		//kalo data sudah tidak ada makan akan muncul panic
		if err != nil {
			panic(err)
		}
		//output
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	//kalo data sudah tidak ada makan akan muncul panic
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//selama masih ada data di rows, maka akan terus diambil
	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		//data/timestamp tetep menggunakan time.Time kalau database cuma sampai tanggal, maka nanti jam menit detiknya akan kosong semua
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool

		//didalam kurung dimasukan parameterdari apa yang akan  kita ambil
		//Pointer dipakai karena kita akan ngeset data dari parameter. Kalo enggak pointer, maka data tidak akan dipakai
		err = rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
		if email.Valid {
			fmt.Println("Email:", email.String)
		}
		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		if birthDate.Valid {
			fmt.Println("Birth Date:", birthDate.Time)
		}
		fmt.Println("Married:", married)
		fmt.Println("Created At:", createdAt)
	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	//sql dari user
	username := "admin'; #"
	password := "salah"

	script := "SELECT username FROM user WHERE username='" + username +
		"' AND password='" + password + "'LIMIT 1"
	fmt.Println(script)
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//hanya satu data
	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Success login", username)
	} else {
		fmt.Println("Gagal login")
	}
}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	//sql dari user
	username := "admin"
	password := "admin"

	//sql dengan parameter
	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"
	fmt.Println(script)
	rows, err := db.QueryContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//hanya satu data
	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Success login", username)
	} else {
		fmt.Println("Gagal login")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "galih; DROP TABLE user;#"
	password := "gaih"

	script := "INSERT INTO user(username, password) VALUES(?,?)"
	//exec tidak akan mengembalikan result
	_, err := db.ExecContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "galih@gmail.com"
	comment := "test"

	script := "INSERT INTO comments(email, comment) VALUES(?,?)"
	//exec tidak akan mengembalikan result
	result, err := db.ExecContext(ctx, script, username, comment)
	if err != nil {
		panic(err)
	}
	//membalikan type data lastId
	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id", insertId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO comments(email,comment) VALUES(?,?)"
	//sql script sudah dibinding disini
	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	//udah ada koneksi ke database,tidak usah menenyai pool kembali
	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "galih" + strconv.Itoa(i) + "gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)
		//sql script sudah di binding di atas
		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment Id", id)
	}

}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	//memulai koneksi
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "INSERT INTO comments(email,comment) VALUES(?,?)"
	//do transaction
	for i := 0; i < 10; i++ {
		email := "galih" + strconv.Itoa(i) + "gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)
		//sql script sudah di binding di atas
		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment Id", id)
	}

	//untuk meneruskan transaksi
	//err = tx.Rollback()
	//untuk mengembalikan transaksi
	err = tx.Rollback()
	if err != nil {
		panic(err)
	}
}

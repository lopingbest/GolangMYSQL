package Golang_MYSQL

import (
	"context"
	"fmt"
	"testing"
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

package repository

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	Golang_MYSQL "golang_mysql"
	"golang_mysql/entity"
	"testing"
)

func TestCommenInsert(t *testing.T) {
	commentRepository := NewCommentRepository(Golang_MYSQL.GetConnection())
	ctx := context.Background()
	comment := entity.Comment{
		Email:   "repository@test.com",
		Comment: "Test Repository",
	}
	commentRepository.Insert(ctx, comment)

	result, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestFindById(t *testing.T) {
	commentRepository := NewCommentRepository(Golang_MYSQL.GetConnection())

	comment, err := commentRepository.FindById(context.Background(), 90)
	if err != nil {
		panic(err)
	}
	fmt.Println(comment)
}

func TestFindALl(t *testing.T) {
	commentRepository := NewCommentRepository(Golang_MYSQL.GetConnection())
	comments, err := commentRepository.FindAll(context.Background())
	if err != nil {
		panic(err)
	}
	for _, comment := range comments {
		fmt.Println(comment)
	}
}

package main

import (
	"log"

	"github.com/go-comm/sqlxmodel"
)

func main() {
	m := sqlxmodel.NewSqlxModel("db")

	err := m.WriteToFile(&User{}, "examples/main/user_model.go")
	if err != nil {
		log.Println(err)
	}
}

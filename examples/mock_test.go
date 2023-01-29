package examples

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/go-comm/sqlxmodel"
)

type MockDB struct{}

func (MockDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if strings.Contains(query, "t_user") {
		return json.Unmarshal([]byte(`{"id":"1000","name":"tom","role_id":1}`), dest)
	}
	if strings.Contains(query, "t_role") {
		return json.Unmarshal([]byte(`{"id":1,"name":"role1"}`), dest)
	}
	return sql.ErrNoRows
}

func (MockDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if strings.Contains(query, "t_user") {
		return json.Unmarshal([]byte(`[{"id":"1000","name":"tom"},{"id":"1001","name":"john"}]`), dest)
	}
	if strings.Contains(query, "t_role") {
		return json.Unmarshal([]byte(`[{"id":1,"name":"role_1"},{"id":2,"name":"role_2"}]`), dest)
	}
	return sql.ErrNoRows
}

var mockdb = new(MockDB)

func init() {
	sqlxmodel.SetShowSQL(true)
}

func TestCreateModel(t *testing.T) {
	m := sqlxmodel.NewSqlxModel("db")

	m.WriteToFile("model.go", &User{}, &Role{})
}

func ExampleGetUser() {
	var u User
	if err := UserModel.QueryFirstByPrimaryKey(context.TODO(), mockdb, &u, "", 1); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(u.ID, u.Name)

	// Output:
	// 1000 tom
}

func ExampleQueryUsers() {
	var users []User
	if err := UserModel.QueryList(context.TODO(), mockdb, &users, "", "id>10"); err != nil {
		log.Fatalln(err)
	}
	for _, u := range users {
		fmt.Println(u.ID, u.Name)
	}

	// Output:
	// 1000 tom
	// 1001 john
}

func ExampleQueryUsersByRef() {
	var users []*User
	users = append(users, &User{Base: Base{ID: "1000"}, Name: "tom", RoleID: 1})
	users = append(users, &User{Base: Base{ID: "1001"}, Name: "join", RoleID: 1})

	if err := sqlxmodel.RelatedWithRef(context.TODO(), mockdb, &users, "Role", "RoleID"); err != nil {
		log.Fatalln(err)
	}

	for _, u := range users {
		fmt.Println(u.ID, u.Name, u.Role.Name)
	}

	// Output:
	// 1000 tom role1
	// 1001 join role1
}

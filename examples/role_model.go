// !!!Don't Edit it!!!
package examples

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-comm/sqlxmodel"
	"github.com/jmoiron/sqlx"
)

var (
	_ context.Context
	_ sql.DB
	_ sqlx.DB
	_ sqlxmodel.SqlxModel

	_ = strings.Join
	_ = fmt.Println
	_ = reflect.ValueOf
)

// RoleModel model of Role
//
// !!!Don't Edit it!!!
var RoleModel = new(Role)

// QueryFirstByPrimaryKey query one record by primary key
//
// var records []*Role
//
// QueryFirstByPrimaryKey(ctx, db, &records, "", 100)
//
// SQL: select `id`,`name` from t_role where `id`=?
//
// !!!Don't Edit it!!!
func (model Role) QueryFirstByPrimaryKey(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, pk interface{}) error {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select `id`,`name`")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_role where `id`=?")
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), pk)
}

// QueryFirst query one record
//
// var record Role
//
// QueryFirst(ctx, db, &record, "", "where `id`=?", 100)
//
// SQL: select `id`,`name` from t_role where `id`=?
//
// !!!Don't Edit it!!!
func (model Role) QueryFirst(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, whereAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select `id`,`name`")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_role")
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if !sqlxmodel.HasPrefixToken(where, "where") {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn("", where, args...)
			sqlBuilder.WriteString(where)
		} else {
			return fmt.Errorf("expect string, but type %T", whereAndArgs[0])
		}
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), args...)
}

// QueryList query all records
//
// var records []*Role
//
// QueryList(ctx, db, &records, "", "where id>? order by id desc", 100)
//
// SQL: select `id`,`name` from t_role where id>? order by id desc
//
// !!!Don't Edit it!!!
func (model Role) QueryList(ctx context.Context, db sqlxmodel.SelectContext, dest interface{}, selection string, whereAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select `id`,`name`")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_role")
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if !sqlxmodel.HasPrefixToken(where, "where") {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn("", where, args...)
			sqlBuilder.WriteString(where)
		} else {
			return fmt.Errorf("expect string, but type %T", whereAndArgs[0])
		}
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.SelectContext(ctx, dest, sqlBuilder.String(), args...)
}

// Update update a record
//
// Update(ctx, db, "`name`=?", "where id=?", 100)
//
// SQL: update t_role set `name`=? where id=?
//
// !!!Don't Edit it!!!
func (model Role) Update(ctx context.Context, db sqlxmodel.ExecContext, section string, whereAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("update t_role")
	if !sqlxmodel.HasPrefixToken(section, "set") {
		sqlBuilder.WriteString(" set ")
	} else {
		sqlBuilder.WriteString(" ")
	}
	sqlBuilder.WriteString(section)
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if !sqlxmodel.HasPrefixToken(where, "where") {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn(section, where, args...)
			sqlBuilder.WriteString(where)
		} else {
			return nil, fmt.Errorf("expect string, but type %T", whereAndArgs[0])
		}
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.ExecContext(ctx, sqlBuilder.String(), args...)
}

// NamedUpdate update a record
//
// NamedUpdate(ctx, db, "", "", &record)
//
// SQL: update t_role set `name`=:name where `id`=?
//
// !!!Don't Edit it!!!
func (model Role) NamedUpdate(ctx context.Context, db sqlxmodel.NamedExecContext, section string, where string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update t_role")
	if section == "" {
		sqlBuilder.WriteString(" `name`=:name")
	} else {
		if !sqlxmodel.HasPrefixToken(section, "select") {
			sqlBuilder.WriteString(" set ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(section)
	}
	if where == "" {
		sqlBuilder.WriteString(" where `id`=:id")
	} else {
		if !sqlxmodel.HasPrefixToken(where, "where") {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(where)
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	if e, ok := values.(interface {
		BeforeUpdate()
	}); ok {
		e.BeforeUpdate()
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
}

// NamedUpdateColumns update a record
//
// NamedUpdateColumns(ctx, db, nil, "", &record)
//
// SQL: update t_role set `name`=:name where `id`=?
//
// columns: []string{"id","version=version+1"} is also supported.
//
// !!!Don't Edit it!!!
func (model Role) NamedUpdateColumns(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, where string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update t_role set")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" `name`=:name")
	} else {
		var formatColumn = func(s string) string {
			var p int = -1
			for i := 0; i < len(s); i++ {
				if s[i] == '=' {
					p = i
					break
				}
			}
			if p < 0 {
				return "`" + s + "`=:" + s
			}
			if p >= len(s)-1 || p <= 0 {
				return ""
			}
			return s
		}
		sqlBuilder.WriteString(" ")
		sqlBuilder.WriteString(formatColumn(columns[0]))
		for i := 1; i < len(columns); i++ {
			sqlBuilder.WriteString(",")
			sqlBuilder.WriteString(formatColumn(columns[i]))
		}
	}
	if where == "" {
		sqlBuilder.WriteString(" where `id`=:id")
	} else {
		if !sqlxmodel.HasPrefixToken(where, "where") {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(where)
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	if e, ok := values.(interface {
		BeforeUpdate()
	}); ok {
		e.BeforeUpdate()
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
}

// Insert insert a record
//
// Insert(ctx, db, &record)
//
// SQL: insert into t_role(`id`,`name`)values(:id,:name)
//
// !!!Don't Edit it!!!
func (model Role) Insert(ctx context.Context, db sqlxmodel.NamedExecContext, values interface{}) (sql.Result, error) {
	s := "insert into t_role(`id`,`name`)values(:id,:name)"
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(s)
	}
	if e, ok := values.(interface {
		BeforeInsert()
	}); ok {
		e.BeforeInsert()
	}
	return db.NamedExecContext(ctx, s, values)
}

// DeleteByPrimaryKey delete one record by primary key
//
// DeleteByPrimaryKey(ctx, db, 100)
//
// SQL: delete from t_role where `id`=?
//
// !!!Don't Edit it!!!
func (model Role) DeleteByPrimaryKey(ctx context.Context, db sqlxmodel.ExecContext, pk interface{}) (sql.Result, error) {
	s := "delete from t_role where `id`=?"
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(s)
	}
	return db.ExecContext(ctx, s, pk)
}

// Delete query records
//
// Delete(ctx, db, "where `id`=?", 100)
//
// SQL: delete from t_role where `id`=?
//
// !!!Don't Edit it!!!
func (model Role) Delete(ctx context.Context, db sqlxmodel.ExecContext, whereAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.WriteString("delete from t_role")
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if !sqlxmodel.HasPrefixToken(where, "where") {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn("", where, args...)
			sqlBuilder.WriteString(where)
		} else {
			return nil, fmt.Errorf("expect string, but type %T", whereAndArgs[0])
		}
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.ExecContext(ctx, sqlBuilder.String(), args...)
}

// Count count
//
// Count(ctx, db, "")
//
// SQL: select count(1) as c from t_role
//
// !!!Don't Edit it!!!
func (model Role) Count(ctx context.Context, db sqlxmodel.QueryRowContext, whereAndArgs ...interface{}) (int64, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("select count(1) as c")
	sqlBuilder.WriteString(" from t_role")
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if !sqlxmodel.HasPrefixToken(where, "where") {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn("", where, args...)
			sqlBuilder.WriteString(where)
		} else {
			return 0, fmt.Errorf("expect string, but type %T", whereAndArgs[0])
		}
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	row := db.QueryRowContext(ctx, sqlBuilder.String(), args...)
	var c int64
	err := row.Scan(&c)
	return c, err
}

// RelatedWith
//
// RelatedWith(ctx, db, "Creater", 1)
//
// !!!Don't Edit it!!!
func (model *Role) RelatedWith(ctx context.Context, db sqlxmodel.GetContext, field string, pk interface{}) error {
	return sqlxmodel.RelatedWith(ctx, db, model, field, pk)
}

// RelatedWithRef
//
// RelatedWithRef(ctx, db, "Creater", "CreaterID")
//
// !!!Don't Edit it!!!
func (model *Role) RelatedWithRef(ctx context.Context, db sqlxmodel.GetContext, field string, ref ...string) error {
	return sqlxmodel.RelatedWithRef(ctx, db, model, field, ref...)
}


// !!!Don't Edit it!!!
package examples

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	_ = log.Println
	_ = reflect.ValueOf
)

// UserModel model of User
//
// !!!Don't Edit it!!!
var UserModel = new(User)

// QueryFirstByPrimaryKey query one record by primary key
//
// var records []*User
//
// QueryFirstByPrimaryKey(ctx, db, &records, "", 100)
//
// SQL: select A.name,A.email,A.role_id,A.id,A.createtime,A.creater,A.modifytime,A.modifier,A.version,A.defunct,A.deleted from t_user where A.id=?
//
// !!!Don't Edit it!!!
func (model User) QueryFirstByPrimaryKey(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, pk interface{}) error {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select A.name,A.email,A.role_id,A.id,A.createtime,A.creater,A.modifytime,A.modifier,A.version,A.defunct,A.deleted")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_user A where A.id=?")
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), pk)
}

// QueryFirst query one record
//
// var record User
//
// QueryFirst(ctx, db, &record, "", "where `id`=?", 100)
//
// SQL: select A.name,A.email,A.role_id,A.id,A.createtime,A.creater,A.modifytime,A.modifier,A.version,A.defunct,A.deleted from t_user A where A.id=?
//
// !!!Don't Edit it!!!
func (model User) QueryFirst(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, clauseAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select A.name,A.email,A.role_id,A.id,A.createtime,A.creater,A.modifytime,A.modifier,A.version,A.defunct,A.deleted")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_user A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
		}
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), args...)
}

// QueryList query all records
//
// var records []*User
//
// QueryList(ctx, db, &records, "", "where id>? order by id desc", 100)
//
// SQL: select A.name,A.email,A.role_id,A.id,A.createtime,A.creater,A.modifytime,A.modifier,A.version,A.defunct,A.deleted from t_user A where id>? order by A.id desc
//
// !!!Don't Edit it!!!
func (model User) QueryList(ctx context.Context, db sqlxmodel.SelectContext, dest interface{}, selection string, clauseAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select A.name,A.email,A.role_id,A.id,A.createtime,A.creater,A.modifytime,A.modifier,A.version,A.defunct,A.deleted")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_user A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
		}
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.SelectContext(ctx, dest, sqlBuilder.String(), args...)
}

// Update update a record
//
// Update(ctx, db, "`name`=?,`email`=?,`role_id`=?,`createtime`=?,`creater`=?,`modifytime`=?,`modifier`=?,`version`=?,`defunct`=?,`deleted`=?", "where id=?", 100)
//
// SQL: update t_user set `name`=?,`email`=?,`role_id`=?,`createtime`=?,`creater`=?,`modifytime`=?,`modifier`=?,`version`=?,`defunct`=?,`deleted`=? where id=?
//
// !!!Don't Edit it!!!
func (model User) Update(ctx context.Context, db sqlxmodel.ExecContext, set string, clauseAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("update t_user A")
	if !sqlxmodel.HasPrefixToken(set, "set") {
		sqlBuilder.WriteString(" set ")
	} else {
		sqlBuilder.WriteString(" ")
	}
	sqlBuilder.WriteString(set)
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, strings.Count(set, "?"))
			sqlBuilder.WriteString(clause)
		} else {
			return nil, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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
// SQL: update t_user set `name`=:name,`email`=:email,`role_id`=:role_id,`createtime`=:createtime,`creater`=:creater,`modifytime`=:modifytime,`modifier`=:modifier,`version`=:version,`defunct`=:defunct,`deleted`=:deleted where `id`=?
//
// !!!Don't Edit it!!!
func (model User) NamedUpdate(ctx context.Context, db sqlxmodel.NamedExecContext, set string, clause string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update t_user")
	if set == "" {
		sqlBuilder.WriteString(" set `name`=:name,`email`=:email,`role_id`=:role_id,`createtime`=:createtime,`creater`=:creater,`modifytime`=:modifytime,`modifier`=:modifier,`version`=:version,`defunct`=:defunct,`deleted`=:deleted")
	} else {
		if !sqlxmodel.HasPrefixToken(set, "set") {
			sqlBuilder.WriteString(" set ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(set)
	}
	if clause == "" {
		sqlBuilder.WriteString(" where `id`=:id")
	} else {
		if sqlxmodel.IfClauseAppendWhere(clause) {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(clause)
	}
	if err := sqlxmodel.BeforeUpdate(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
}

// NamedUpdateColumns update a record
//
// NamedUpdateColumns(ctx, db, nil, "", &record)
//
// SQL: update t_user set `name`=:name,`email`=:email,`role_id`=:role_id,`createtime`=:createtime,`creater`=:creater,`modifytime`=:modifytime,`modifier`=:modifier,`version`=:version,`defunct`=:defunct,`deleted`=:deleted where `id`=?
//
// columns: []string{"id","version=version+1"} is also supported.
//
// !!!Don't Edit it!!!
func (model User) NamedUpdateColumns(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, clause string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update t_user set")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" `name`=:name,`email`=:email,`role_id`=:role_id,`createtime`=:createtime,`creater`=:creater,`modifytime`=:modifytime,`modifier`=:modifier,`version`=:version,`defunct`=:defunct,`deleted`=:deleted")
	} else {
		sqlBuilder.WriteString(" ")
		sqlBuilder.WriteString(sqlxmodel.FormatSetClause(columns[0]))
		for i := 1; i < len(columns); i++ {
			sqlBuilder.WriteString(",")
			sqlBuilder.WriteString(sqlxmodel.FormatSetClause(columns[i]))
		}
	}
	if clause == "" {
		sqlBuilder.WriteString(" where `id`=:id")
	} else {
		if sqlxmodel.IfClauseAppendWhere(clause) {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(clause)
	}
	if err := sqlxmodel.BeforeUpdate(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
}

// Insert insert a record
//
// Insert(ctx, db, &record)
//
// SQL: insert into t_user(`name`,`email`,`role_id`,`id`,`createtime`,`creater`,`modifytime`,`modifier`,`version`,`defunct`,`deleted`)values(:name,:email,:role_id,:id,:createtime,:creater,:modifytime,:modifier,:version,:defunct,:deleted)
//
// !!!Don't Edit it!!!
func (model User) Insert(ctx context.Context, db sqlxmodel.NamedExecContext, values interface{}) (sql.Result, error) {
	s := "insert into t_user(`name`,`email`,`role_id`,`id`,`createtime`,`creater`,`modifytime`,`modifier`,`version`,`defunct`,`deleted`)values(:name,:email,:role_id,:id,:createtime,:creater,:modifytime,:modifier,:version,:defunct,:deleted)"
	if err := sqlxmodel.BeforeInsert(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(s)
	}
	return db.NamedExecContext(ctx, s, values)
}

// SaveOnMysql insert a record
//
// SaveOnMysql(ctx, db, &record)
//
// SQL: insert into t_user(`name`,`email`,`role_id`,`id`,`createtime`,`creater`,`modifytime`,`modifier`,`version`,`defunct`,`deleted`)values(:name,:email,:role_id,:id,:createtime,:creater,:modifytime,:modifier,:version,:defunct,:deleted)
//
// !!!Don't Edit it!!!
func (model User) SaveOnMysql(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(256)
	sqlBuilder.WriteString("insert into t_user(`name`,`email`,`role_id`,`id`,`createtime`,`creater`,`modifytime`,`modifier`,`version`,`defunct`,`deleted`)values(:name,:email,:role_id,:id,:createtime,:creater,:modifytime,:modifier,:version,:defunct,:deleted) on duplicate key update")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" `name`=values(`name`),`email`=values(`email`),`role_id`=values(`role_id`),`createtime`=values(`createtime`),`creater`=values(`creater`),`modifytime`=values(`modifytime`),`modifier`=values(`modifier`),`version`=values(`version`),`defunct`=values(`defunct`),`deleted`=values(`deleted`)")
	} else {
		formatColumn := func(s string) string {
			return "`" + s + "`=values(`" + s + "`)"
		}
		sqlBuilder.WriteString(" ")
		sqlBuilder.WriteString(formatColumn(columns[0]))
		for i := 1; i < len(columns); i++ {
			sqlBuilder.WriteString(",")
			sqlBuilder.WriteString(formatColumn(columns[i]))
		}
	}
	if err := sqlxmodel.BeforeInsert(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
}

// DeleteByPrimaryKey delete one record by primary key
//
// DeleteByPrimaryKey(ctx, db, 100)
//
// SQL: delete from t_user where `id`=?
//
// !!!Don't Edit it!!!
func (model User) DeleteByPrimaryKey(ctx context.Context, db sqlxmodel.ExecContext, pk interface{}) (sql.Result, error) {
	s := "delete from t_user where `id`=?"
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(s)
	}
	return db.ExecContext(ctx, s, pk)
}

// Delete query records
//
// Delete(ctx, db, "where `id`=?", 100)
//
// SQL: delete from t_user where `id`=?
//
// !!!Don't Edit it!!!
func (model User) Delete(ctx context.Context, db sqlxmodel.ExecContext, clauseAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.WriteString("delete from t_user")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return nil, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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
// SQL: select count(1) as c from t_user A
//
// !!!Don't Edit it!!!
func (model User) Count(ctx context.Context, db sqlxmodel.QueryRowContext, clauseAndArgs ...interface{}) (int64, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("select count(1) as c from t_user A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return 0, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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

// Has has record
//
// Has(ctx, db, "id=1")
//
// SQL: select 1 from t_user A where id=1 limit 1
//
// !!!Don't Edit it!!!
func (model User) Has(ctx context.Context, db sqlxmodel.QueryRowContext, clauseAndArgs ...interface{}) (bool, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("select 1 from t_user A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return false, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
		}
	}
	sqlBuilder.WriteString(" limit 1")
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	row := db.QueryRowContext(ctx, sqlBuilder.String(), args...)
	var c int64
	err := row.Scan(&c)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return c == 1, nil
}

// RelatedWith
//
// RelatedWith(ctx, db, "Creater", 1)
//
// !!!Don't Edit it!!!
func (model *User) RelatedWith(ctx context.Context, db sqlxmodel.GetContext, field string, pk interface{}) error {
	return sqlxmodel.RelatedWith(ctx, db, model, field, pk)
}

// RelatedWithRef
//
// RelatedWithRef(ctx, db, "Creater", "CreaterID")
//
// !!!Don't Edit it!!!
func (model *User) RelatedWithRef(ctx context.Context, db sqlxmodel.GetContext, field string, ref ...string) error {
	return sqlxmodel.RelatedWithRef(ctx, db, model, field, ref...)
}

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
// SQL: select A.id,A.name from t_role where A.id=?
//
// !!!Don't Edit it!!!
func (model Role) QueryFirstByPrimaryKey(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, pk interface{}) error {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select A.id,A.name")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_role A where A.id=?")
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
// SQL: select A.id,A.name from t_role A where A.id=?
//
// !!!Don't Edit it!!!
func (model Role) QueryFirst(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, clauseAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select A.id,A.name")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_role A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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
// SQL: select A.id,A.name from t_role A where id>? order by A.id desc
//
// !!!Don't Edit it!!!
func (model Role) QueryList(ctx context.Context, db sqlxmodel.SelectContext, dest interface{}, selection string, clauseAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select A.id,A.name")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from t_role A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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
func (model Role) Update(ctx context.Context, db sqlxmodel.ExecContext, set string, clauseAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("update t_role A")
	if !sqlxmodel.HasPrefixToken(set, "set") {
		sqlBuilder.WriteString(" set ")
	} else {
		sqlBuilder.WriteString(" ")
	}
	sqlBuilder.WriteString(set)
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, strings.Count(set, "?"))
			sqlBuilder.WriteString(clause)
		} else {
			return nil, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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
func (model Role) NamedUpdate(ctx context.Context, db sqlxmodel.NamedExecContext, set string, clause string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update t_role")
	if set == "" {
		sqlBuilder.WriteString(" set `name`=:name")
	} else {
		if !sqlxmodel.HasPrefixToken(set, "set") {
			sqlBuilder.WriteString(" set ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(set)
	}
	if clause == "" {
		sqlBuilder.WriteString(" where `id`=:id")
	} else {
		if sqlxmodel.IfClauseAppendWhere(clause) {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(clause)
	}
	if err := sqlxmodel.BeforeUpdate(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
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
func (model Role) NamedUpdateColumns(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, clause string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update t_role set")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" `name`=:name")
	} else {
		sqlBuilder.WriteString(" ")
		sqlBuilder.WriteString(sqlxmodel.FormatSetClause(columns[0]))
		for i := 1; i < len(columns); i++ {
			sqlBuilder.WriteString(",")
			sqlBuilder.WriteString(sqlxmodel.FormatSetClause(columns[i]))
		}
	}
	if clause == "" {
		sqlBuilder.WriteString(" where `id`=:id")
	} else {
		if sqlxmodel.IfClauseAppendWhere(clause) {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(clause)
	}
	if err := sqlxmodel.BeforeUpdate(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
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
	if err := sqlxmodel.BeforeInsert(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(s)
	}
	return db.NamedExecContext(ctx, s, values)
}

// SaveOnMysql insert a record
//
// SaveOnMysql(ctx, db, &record)
//
// SQL: insert into t_role(`id`,`name`)values(:id,:name)
//
// !!!Don't Edit it!!!
func (model Role) SaveOnMysql(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(256)
	sqlBuilder.WriteString("insert into t_role(`id`,`name`)values(:id,:name) on duplicate key update")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" `name`=values(`name`)")
	} else {
		formatColumn := func(s string) string {
			return "`" + s + "`=values(`" + s + "`)"
		}
		sqlBuilder.WriteString(" ")
		sqlBuilder.WriteString(formatColumn(columns[0]))
		for i := 1; i < len(columns); i++ {
			sqlBuilder.WriteString(",")
			sqlBuilder.WriteString(formatColumn(columns[i]))
		}
	}
	if err := sqlxmodel.BeforeInsert(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
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
func (model Role) Delete(ctx context.Context, db sqlxmodel.ExecContext, clauseAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.WriteString("delete from t_role")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return nil, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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
// SQL: select count(1) as c from t_role A
//
// !!!Don't Edit it!!!
func (model Role) Count(ctx context.Context, db sqlxmodel.QueryRowContext, clauseAndArgs ...interface{}) (int64, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("select count(1) as c from t_role A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return 0, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
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

// Has has record
//
// Has(ctx, db, "id=1")
//
// SQL: select 1 from t_role A where id=1 limit 1
//
// !!!Don't Edit it!!!
func (model Role) Has(ctx context.Context, db sqlxmodel.QueryRowContext, clauseAndArgs ...interface{}) (bool, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("select 1 from t_role A")
	if len(clauseAndArgs) > 0 {
		args = clauseAndArgs[1:]
		if clause, ok := clauseAndArgs[0].(string); ok {
			if sqlxmodel.IfClauseAppendWhere(clause) {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			clause, args = sqlxmodel.WithIn(clause, args, 0)
			sqlBuilder.WriteString(clause)
		} else {
			return false, fmt.Errorf("expect string, but type %T", clauseAndArgs[0])
		}
	}
	sqlBuilder.WriteString(" limit 1")
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	row := db.QueryRowContext(ctx, sqlBuilder.String(), args...)
	var c int64
	err := row.Scan(&c)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return c == 1, nil
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

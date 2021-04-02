package sqlxmodel

import (
	"strings"
	"text/template"
)

func getTpl() string {
	return `// !!!Don't Edit it!!!
package {{ .PackageName }}

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

// {{ .Name | Title }}Model model of {{ .Name }}
var {{ .Name | Title }}Model = new({{ .Name }})

// QueryFirstByPrimaryKey query one record by primary key
//
// var records []*{{ .Name }}
//
// QueryFirstByPrimaryKey(ctx, db, &records, "", 100)
//
// SQL: select {{ JoinExpr .Fields "${.FormattedField}" }} from {{ .TableName }} where {{ FormattedField .PrimaryKey }}=? limit 1
func (model {{ .Name | Title }}) QueryFirstByPrimaryKey(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, pk interface{}) error {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select {{ JoinExpr .Fields "${.FormattedField}" }}")
	} else {
		if strings.Index(selection, "select") < 0 {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from {{ .TableName }} where {{ FormattedField .PrimaryKey }}=? limit 1")
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), pk)
}

// QueryFirst query one record
//
// var record {{ .Name }}
//
// QueryFirst(ctx, db, &record, "", "where {{ FormattedField .PrimaryKey }}=?", 100)
//
// SQL: select {{ JoinExpr .Fields "${.FormattedField}" }} from {{ .TableName }} where {{ FormattedField .PrimaryKey }}=? limit 1
func (model {{ .Name | Title }}) QueryFirst(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, whereAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select {{ JoinExpr .Fields "${.FormattedField}" }}")
	} else {
		if strings.Index(selection, "select") < 0 {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from {{ .TableName }}")
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if strings.Index(where, "where") < 0 {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn(where, args...)
			sqlBuilder.WriteString(where)
		} else {
			return fmt.Errorf("expect string, but type %T", whereAndArgs[0])
		}
	}
	sqlBuilder.WriteString(" limit 1")
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), args...)
}

// QueryList query all records
//
// var records []*{{ .Name }}
//
// QueryList(ctx, db, &records, "", "where {{ .PrimaryKey }}>? order by {{ .PrimaryKey }} desc", 100)
//
// SQL: select {{ JoinExpr .Fields "${.FormattedField}" }} from {{ .TableName }} where {{ .PrimaryKey }}>? order by {{ .PrimaryKey }} desc
func (model {{ .Name | Title }}) QueryList(ctx context.Context, db sqlxmodel.SelectContext, dest interface{}, selection string, whereAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select {{ JoinExpr .Fields "${.FormattedField}" }}")
	} else {
		if strings.Index(selection, "select") < 0 {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from {{ .TableName }}")
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if strings.Index(where, "where") < 0 {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn(where, args...)
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
// Update(ctx, db, "{{ JoinExpr .Fields "${.FormattedField}=?" .PrimaryKey }}", "where {{ .PrimaryKey }}=?", 100)
//
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=?" .PrimaryKey }} where {{ .PrimaryKey }}=?
func (model {{ .Name | Title }}) Update(ctx context.Context, db sqlxmodel.ExecContext, selection string, whereAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("update {{ .TableName }} set")
	sqlBuilder.WriteString(selection)
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if strings.Index(where, "where") < 0 {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn(where, args...)
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
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }} where {{ FormattedField .PrimaryKey }}=?
func (model {{ .Name | Title }}) NamedUpdate(ctx context.Context, db sqlxmodel.NamedExecContext, selection string, where string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update {{ .TableName }} set")
	if selection == "" {
		sqlBuilder.WriteString(" {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }}")
	} else {
		sqlBuilder.WriteString(selection)
	}
	if where == "" {
		sqlBuilder.WriteString(" where {{ FormattedField .PrimaryKey }}=:{{ .PrimaryKey }}")
	} else {
		if strings.Index(where, "where") < 0 {
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
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }} where {{ FormattedField .PrimaryKey }}=?
//
// columns: []string{"id","version=version+1"} is also supported.
func (model {{ .Name | Title }}) NamedUpdateColumns(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, where string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update {{ .TableName }} set")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }}")
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
				return "` + "`" + `" + s + "` + "`" + `=:" + s
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
		sqlBuilder.WriteString(" where {{ FormattedField .PrimaryKey }}=:{{ .PrimaryKey }}")
	} else {
		if strings.Index(where, "where") < 0 {
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
// SQL: insert into {{ .TableName }}({{ JoinExpr .Fields "${.FormattedField}" }})values({{ JoinExpr .Fields ":${.Field}" }})
func (model {{ .Name | Title }}) Insert(ctx context.Context, db sqlxmodel.NamedExecContext, values interface{}) (sql.Result, error) {
	s := "insert into {{ .TableName }}({{ JoinExpr .Fields "${.FormattedField}" }})values({{ JoinExpr .Fields ":${.Field}" }})"
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
// SQL: delete from {{ .TableName }} where {{ FormattedField .PrimaryKey }}=?
func (model {{ .Name | Title }}) DeleteByPrimaryKey(ctx context.Context, db sqlxmodel.ExecContext, pk interface{}) (sql.Result, error) {
	s := "delete from {{ .TableName }} where {{ FormattedField .PrimaryKey }}=?"
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(s)
	}
	return db.ExecContext(ctx, s, pk)
}

// Delete query records
//
// Delete(ctx, db, "where {{ FormattedField .PrimaryKey }}=?", 100)
//
// SQL: delete from {{ .TableName }} where {{ FormattedField .PrimaryKey }}=?
func (model {{ .Name | Title }}) Delete(ctx context.Context, db sqlxmodel.ExecContext, whereAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.WriteString("delete from {{ .TableName }}")
	if len(whereAndArgs) > 0 {
		args = whereAndArgs[1:]
		if where, ok := whereAndArgs[0].(string); ok {
			if strings.Index(where, "where") < 0 {
				sqlBuilder.WriteString(" where ")
			} else {
				sqlBuilder.WriteString(" ")
			}
			where, args = sqlxmodel.WithIn(where, args...)
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
`
}

func isEmpty(s string) bool {
	return len(s) <= 0
}

func join(fields []*ModelFieldInfo, sep string, fn func(e map[string]string) string, ignores ...string) string {
	var fs []*ModelFieldInfo
	for i := 0; i < len(fields); i++ {
		have := false
		for j := 0; j < len(ignores); j++ {
			if fields[i].FieldName == ignores[j] {
				have = true
				break
			}
		}
		if !have {
			fs = append(fs, fields[i])
		}
	}
	if len(fs) <= 0 {
		return ""
	}
	var s strings.Builder
	for i := 0; i < len(fs); i++ {
		if i > 0 {
			s.WriteString(sep)
		}
		e := map[string]string{
			"Field":          fs[i].FieldName,
			"NamedField":     namedField(fs[i].FieldName),
			"FormattedField": formattedField(fs[i].FieldName),
		}
		s.WriteString(fn(e))
	}
	return s.String()
}

func joinExpr(fields []*ModelFieldInfo, expr string, ignores ...string) string {
	tpl := template.New("")
	tpl.Delims("${", "}")
	template.Must(tpl.Parse(expr))
	return join(fields, ",", func(e map[string]string) string {
		var s strings.Builder
		tpl.Execute(&s, e)
		return s.String()
	}, ignores...)
}

func namedField(s string) string {
	return ":" + s
}

func formattedField(s string) string {
	return "`" + s + "`"
}

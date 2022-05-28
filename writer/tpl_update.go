package writer

var fnUpdate = `
// Update update a record
//
// Update(ctx, db, "{{ JoinExpr .Fields "${.FormattedField}=?" .PrimaryKey }}", "where {{ .PrimaryKey }}=?", 100)
//
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=?" .PrimaryKey }} where {{ .PrimaryKey }}=?
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) Update(ctx context.Context, db sqlxmodel.ExecContext, section string, whereAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("update {{ .TableName }}")
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
`

var fnNamedUpdate = `
// NamedUpdate update a record
//
// NamedUpdate(ctx, db, "", "", &record)
//
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }} where {{ FormattedField .PrimaryKey }}=?
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) NamedUpdate(ctx context.Context, db sqlxmodel.NamedExecContext, section string, where string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update {{ .TableName }}")
	if section == "" {
		sqlBuilder.WriteString(" set {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }}")
	} else {
		if !sqlxmodel.HasPrefixToken(section, "set") {
			sqlBuilder.WriteString(" set ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(section)
	}
	if where == "" {
		sqlBuilder.WriteString(" where {{ FormattedField .PrimaryKey }}=:{{ .PrimaryKey }}")
	} else {
		if !sqlxmodel.HasPrefixToken(where, "where") {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(where)
	}
	if err := sqlxmodel.BeforeUpdate(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
}
`

var fnNamedUpdateColumns = `
// NamedUpdateColumns update a record
//
// NamedUpdateColumns(ctx, db, nil, "", &record)
//
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }} where {{ FormattedField .PrimaryKey }}=?
//
// columns: []string{"id","version=version+1"} is also supported.
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) NamedUpdateColumns(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, where string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update {{ .TableName }} set")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }}")
	} else {
		formatColumn := func(s string) string {
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
		if !sqlxmodel.HasPrefixToken(where, "where") {
			sqlBuilder.WriteString(" where ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(where)
	}
	if err := sqlxmodel.BeforeUpdate(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String())
	}
	return db.NamedExecContext(ctx, sqlBuilder.String(), values)
}
`

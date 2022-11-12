package writer

var fnUpdate = `
// Update update a record
//
// Update(ctx, db, "{{ JoinExpr .Fields "${.FormattedField}=?" .PrimaryKey }}", "where {{ .PrimaryKey }}=?", 100)
//
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=?" .PrimaryKey }} where {{ .PrimaryKey }}=?
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) Update(ctx context.Context, db sqlxmodel.ExecContext, set string, clauseAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("update {{ .TableName }} A")
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
`

var fnNamedUpdate = `
// NamedUpdate update a record
//
// NamedUpdate(ctx, db, "", "", &record)
//
// SQL: update {{ .TableName }} set {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }} where {{ FormattedField .PrimaryKey }}=?
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) NamedUpdate(ctx context.Context, db sqlxmodel.NamedExecContext, set string, clause string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update {{ .TableName }}")
	if set == "" {
		sqlBuilder.WriteString(" set {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }}")
	} else {
		if !sqlxmodel.HasPrefixToken(set, "set") {
			sqlBuilder.WriteString(" set ")
		} else {
			sqlBuilder.WriteString(" ")
		}
		sqlBuilder.WriteString(set)
	}
	if clause == "" {
		sqlBuilder.WriteString(" where {{ FormattedField .PrimaryKey }}=:{{ .PrimaryKey }}")
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
func (model {{ .Name | Title }}) NamedUpdateColumns(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, clause string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	sqlBuilder.WriteString("update {{ .TableName }} set")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" {{ JoinExpr .Fields "${.FormattedField}=:${.Field}" .PrimaryKey }}")
	} else {
		sqlBuilder.WriteString(" ")
		sqlBuilder.WriteString(sqlxmodel.FormatSetClause(columns[0]))
		for i := 1; i < len(columns); i++ {
			sqlBuilder.WriteString(",")
			sqlBuilder.WriteString(sqlxmodel.FormatSetClause(columns[i]))
		}
	}
	if clause == "" {
		sqlBuilder.WriteString(" where {{ FormattedField .PrimaryKey }}=:{{ .PrimaryKey }}")
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
`

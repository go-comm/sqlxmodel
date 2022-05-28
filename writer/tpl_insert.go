package writer

var fnInsert = `
// Insert insert a record
//
// Insert(ctx, db, &record)
//
// SQL: insert into {{ .TableName }}({{ JoinExpr .Fields "${.FormattedField}" }})values({{ JoinExpr .Fields ":${.Field}" }})
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) Insert(ctx context.Context, db sqlxmodel.NamedExecContext, values interface{}) (sql.Result, error) {
	s := "insert into {{ .TableName }}({{ JoinExpr .Fields "${.FormattedField}" }})values({{ JoinExpr .Fields ":${.Field}" }})"
	if err := sqlxmodel.BeforeInsert(ctx, values); err != nil {
		return nil, err
	}
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(s)
	}
	return db.NamedExecContext(ctx, s, values)
}
`

var fnSaveOnMysql = `
// SaveOnMysql insert a record
//
// SaveOnMysql(ctx, db, &record)
//
// SQL: insert into {{ .TableName }}({{ JoinExpr .Fields "${.FormattedField}" }})values({{ JoinExpr .Fields ":${.Field}" }})
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) SaveOnMysql(ctx context.Context, db sqlxmodel.NamedExecContext, columns []string, values interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(256)
	sqlBuilder.WriteString("insert into {{ .TableName }}({{ JoinExpr .Fields "${.FormattedField}" }})values({{ JoinExpr .Fields ":${.Field}" }}) on duplicate key update")
	if len(columns) == 0 {
		sqlBuilder.WriteString(" {{ JoinExpr .Fields "` + "`" + `${.Field}` + "`" + `=values(` + "`" + `${.Field}` + "`" + `)" .PrimaryKey}}")
	} else {
		formatColumn := func(s string) string {
			return "` + "`" + `" + s + "` + "`" + `=values(` + "`" + `" + s + "` + "`" + `)"
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
`

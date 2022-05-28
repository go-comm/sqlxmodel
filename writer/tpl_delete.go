package writer

var fnDeleteByPrimaryKey = `
// DeleteByPrimaryKey delete one record by primary key
//
// DeleteByPrimaryKey(ctx, db, 100)
//
// SQL: delete from {{ .TableName }} where {{ FormattedField .PrimaryKey }}=?
//
// !!!Don't Edit it!!!
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
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) Delete(ctx context.Context, db sqlxmodel.ExecContext, whereAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.WriteString("delete from {{ .TableName }}")
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
`

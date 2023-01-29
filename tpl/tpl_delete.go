package tpl

var FnDeleteByPrimaryKey = `
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
		sqlxmodel.PrintSQL(s, pk)
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
func (model {{ .Name | Title }}) Delete(ctx context.Context, db sqlxmodel.ExecContext, clauseAndArgs ...interface{}) (sql.Result, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.WriteString("delete from {{ .TableName }}")
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
		sqlxmodel.PrintSQL(sqlBuilder.String(), args...)
	}
	return db.ExecContext(ctx, sqlBuilder.String(), args...)
}
`

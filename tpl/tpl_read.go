package tpl

var FnQueryFirstByPrimaryKey = `
// QueryFirstByPrimaryKey query one record by primary key
//
// var records []*{{ .Name }}
//
// QueryFirstByPrimaryKey(ctx, db, &records, "", 100)
//
// SQL: select {{ JoinExpr .Fields "A.${.Field}" }} from {{ .TableName }} where A.{{ .PrimaryKey }}=?
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) QueryFirstByPrimaryKey(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, pk interface{}) error {
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select {{ JoinExpr .Fields "A.${.Field}" }}")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from {{ .TableName }} A where A.{{ .PrimaryKey }}=?")
	if sqlxmodel.ShowSQL() {
		sqlxmodel.PrintSQL(sqlBuilder.String(), pk)
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), pk)
}
`

var FnQueryFirst = `
// QueryFirst query one record
//
// var record {{ .Name }}
//
// QueryFirst(ctx, db, &record, "", "where {{ FormattedField .PrimaryKey }}=?", 100)
//
// SQL: select {{ JoinExpr .Fields "A.${.Field}" }} from {{ .TableName }} A where A.{{ .PrimaryKey }}=?
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) QueryFirst(ctx context.Context, db sqlxmodel.GetContext, dest interface{}, selection string, clauseAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select {{ JoinExpr .Fields "A.${.Field}" }}")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from {{ .TableName }} A")
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
		sqlxmodel.PrintSQL(sqlBuilder.String(), args...)
	}
	return db.GetContext(ctx, dest, sqlBuilder.String(), args...)
}
`

var FnQueryList = `
// QueryList query all records
//
// var records []*{{ .Name }}
//
// QueryList(ctx, db, &records, "", "where {{ .PrimaryKey }}>? order by {{ .PrimaryKey }} desc", 100)
//
// SQL: select {{ JoinExpr .Fields "A.${.Field}" }} from {{ .TableName }} A where {{ .PrimaryKey }}>? order by A.{{ .PrimaryKey }} desc
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) QueryList(ctx context.Context, db sqlxmodel.SelectContext, dest interface{}, selection string, clauseAndArgs ...interface{}) error {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(128)
	if selection == "" {
		sqlBuilder.WriteString("select {{ JoinExpr .Fields "A.${.Field}" }}")
	} else {
		if !sqlxmodel.HasPrefixToken(selection, "select") {
			sqlBuilder.WriteString("select ")
		}
		sqlBuilder.WriteString(selection)
	}
	sqlBuilder.WriteString(" from {{ .TableName }} A")
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
		sqlxmodel.PrintSQL(sqlBuilder.String(), args...)
	}
	return db.SelectContext(ctx, dest, sqlBuilder.String(), args...)
}
`

var FnCount = `
// Count count
//
// Count(ctx, db, "")
//
// SQL: select count(1) as c from {{ .TableName }} A
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) Count(ctx context.Context, db sqlxmodel.QueryRowContext, clauseAndArgs ...interface{}) (int64, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("select count(1) as c from {{ .TableName }} A")
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
		sqlxmodel.PrintSQL(sqlBuilder.String(), args...)
	}
	row := db.QueryRowContext(ctx, sqlBuilder.String(), args...)
	var c int64
	err := row.Scan(&c)
	return c, err
}
`

var FnHas = `
// Has has record
//
// Has(ctx, db, "id=1")
//
// SQL: select 1 from {{ .TableName }} A where id=1 limit 1
//
// !!!Don't Edit it!!!
func (model {{ .Name | Title }}) Has(ctx context.Context, db sqlxmodel.QueryRowContext, clauseAndArgs ...interface{}) (bool, error) {
	var sqlBuilder strings.Builder
	var args []interface{}
	sqlBuilder.Grow(64)
	sqlBuilder.WriteString("select 1 from {{ .TableName }} A")
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
		sqlxmodel.PrintSQL(sqlBuilder.String(), args...)
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
`

var FnRelatedWith = `
// RelatedWith
//
// RelatedWith(ctx, db, "Creater", 1)
//
// !!!Don't Edit it!!!
func (model *{{ .Name | Title }}) RelatedWith(ctx context.Context, db sqlxmodel.GetContext, field string, pk interface{}) error {
	return sqlxmodel.RelatedWith(ctx, db, model, field, pk)
}
`

var FnRelatedWithRef = `
// RelatedWithRef
//
// RelatedWithRef(ctx, db, "Creater", "CreaterID")
//
// !!!Don't Edit it!!!
func (model *{{ .Name | Title }}) RelatedWithRef(ctx context.Context, db sqlxmodel.GetContext, field string, ref ...string) error {
	return sqlxmodel.RelatedWithRef(ctx, db, model, field, ref...)
}
`

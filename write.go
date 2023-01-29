package sqlxmodel

import (
	"io"
	"text/template"

	"github.com/go-comm/sqlxmodel/tpl"
)

var funcs = []struct {
	Syntax string
	lvl    int
}{
	{tpl.FnQueryFirstByPrimaryKey, LREAD},
	{tpl.FnQueryFirst, LREAD},
	{tpl.FnQueryList, LREAD},
	{tpl.FnUpdate, LUPDATE},
	{tpl.FnNamedUpdate, LUPDATE},
	{tpl.FnNamedUpdateColumns, LUPDATE},
	{tpl.FnInsert, LCREATE},
	{tpl.FnInsertIgnore, LCREATE},
	{tpl.FnSaveOnMysql, LCREATE},
	{tpl.FnDeleteByPrimaryKey, LDELETE},
	{tpl.FnCount, LREAD},
	{tpl.FnHas, LREAD},
	{tpl.FnRelatedWith, LREAD},
	{tpl.FnRelatedWithRef, LREAD},
}

func WriteHeader(t *template.Template, w io.Writer, data interface{}) error {
	if err := parseAndExecute(t, w, tpl.Header, data); err != nil {
		return err
	}
	return nil
}

func WriteBody(t *template.Template, w io.Writer, data interface{}, lvl int) error {
	if err := parseAndExecute(t, w, tpl.Model, data); err != nil {
		return err
	}
	for _, fn := range funcs {
		if fn.lvl&lvl == fn.lvl {
			if err := parseAndExecute(t, w, fn.Syntax, data); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseAndExecute(t *template.Template, w io.Writer, fmt string, data interface{}) error {
	var err error
	t, err = t.Parse(fmt)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

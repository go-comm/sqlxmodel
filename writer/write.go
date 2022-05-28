package writer

import (
	"io"
	"text/template"
)

var funcs = []string{
	fnQueryFirstByPrimaryKey,
	fnQueryFirst,
	fnQueryList,
	fnUpdate,
	fnNamedUpdate,
	fnNamedUpdateColumns,
	fnInsert,
	fnSaveOnMysql,
	fnDeleteByPrimaryKey,
	fnCount,
	fnHas,
	fnRelatedWith,
	fnRelatedWithRef,
}

func WriteHeader(t *template.Template, w io.Writer, data interface{}) error {
	if err := parseAndExecute(t, w, tplHeader, data); err != nil {
		return err
	}
	return nil
}

func WriteBody(t *template.Template, w io.Writer, data interface{}) error {
	if err := parseAndExecute(t, w, tplModel, data); err != nil {
		return err
	}
	for _, fn := range funcs {
		if err := parseAndExecute(t, w, fn, data); err != nil {
			return err
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

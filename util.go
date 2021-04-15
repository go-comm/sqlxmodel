package sqlxmodel

import (
	"reflect"
	"strings"
	"text/template"
	"unicode"
)

func WithIn(section string, where string, args ...interface{}) (string, []interface{}) {
	cnt := strings.Count(section, "?")
	if cnt <= 0 {
		cnt = 0
	}
	pIn := strings.Index(where, "in")
	if pIn < 0 {
		return where, args
	}
	isSpace := func(r byte) bool {
		switch r {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return true
		}
		return false
	}
	// find '?' after 'in'
	pQ := pIn + 2
	for ; pQ < len(where) && isSpace(where[pQ]); pQ++ {
	}
	if !(pQ < len(where) && where[pQ] == '?') {
		return where, args
	}
	c := strings.Count(where[:pIn], "?")
	c += cnt
	if c >= len(args) {
		return where, args
	}
	tv := reflect.TypeOf(args[c])
	if !(args[c] == nil || tv.Kind() == reflect.Slice || tv.Kind() == reflect.Array) {
		return where, args
	}
	rv := reflect.ValueOf(args[c])
	var s strings.Builder
	var nargs []interface{}
	s.WriteString(where[:pQ])
	nargs = append(nargs, args[:c]...)
	if args[c] == nil || rv.Len() <= 0 {
		s.WriteString("(NULL)")
	} else {
		s.WriteByte('(')
		for i := 0; i < rv.Len(); i++ {
			if i > 0 {
				s.WriteByte(',')
			}
			s.WriteByte('?')
			nargs = append(nargs, rv.Index(i).Interface())
		}
		s.WriteByte(')')
	}
	s.WriteString(where[pQ+1:])
	nargs = append(nargs, args[c+1:]...)
	return s.String(), nargs
}

func JoinSlice(x interface{}, args ...interface{}) []interface{} {
	s := make([]interface{}, 0, len(args)+1)
	s = append(s, x)
	s = append(s, args...)
	return s
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

func lowerTitle(s string) string {
	rs := []rune(s)
	if len(rs) <= 0 {
		return ""
	}
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}

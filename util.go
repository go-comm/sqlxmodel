package sqlxmodel

import (
	"reflect"
	"strings"
	"text/template"
	"unicode"
)

func isSpace(r byte) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

func HasPrefixToken(s string, token string) bool {
	i := 0
	for ; i < len(s) && isSpace(s[i]); i++ {
	}
	return strings.HasPrefix(s[i:], token)
}

func BeforeInsert(values interface{}) {
	switch vs := values.(type) {
	case interface{ BeforeInsert() }:
		vs.BeforeInsert()
	default:
		rv := reflect.Indirect(reflect.ValueOf(vs))
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			for i := rv.Len() - 1; i >= 0; i-- {
				e := rv.Index(i)
				if n, ok := e.Interface().(interface{ BeforeInsert() }); ok {
					n.BeforeInsert()
				}
			}
		}
	}
}

func BeforeUpdate(values interface{}) {
	switch vs := values.(type) {
	case interface{ BeforeUpdate() }:
		vs.BeforeUpdate()
	default:
		rv := reflect.Indirect(reflect.ValueOf(vs))
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			for i := rv.Len() - 1; i >= 0; i-- {
				e := rv.Index(i)
				if n, ok := e.Interface().(interface{ BeforeUpdate() }); ok {
					n.BeforeUpdate()
				}
			}
		}
	}
}

func WithIn(section string, where string, args ...interface{}) (string, []interface{}) {
	cnt := strings.Count(section, "?")
	if cnt < 0 {
		cnt = 0
	}
	if cnt >= len(args) {
		return where, args
	}
	var nargs []interface{}
	var s strings.Builder
	var off int = -1
	nargs = append(nargs, args[:cnt]...)
	args = args[cnt:]
	for {
		off = strings.IndexByte(where, '?')
		if off < 0 {
			s.WriteString(where)
			nargs = append(nargs, args...)
			break
		}
		if len(args) <= 0 {
			s.WriteString(where)
			break
		}
		rv := reflect.ValueOf(args[0])
		rt := rv.Type()
		if !(args[0] == nil || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array) {
			s.WriteString(where[:off+1])
			nargs = append(nargs, args[0])
			where = where[off+1:]
			args = args[1:]
			continue
		}
		if args[0] == nil || rv.Len() <= 0 {
			s.WriteString("(NULL)")
		} else {
			s.WriteString(where[:off])
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
		where = where[off+1:]
		args = args[1:]
	}
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

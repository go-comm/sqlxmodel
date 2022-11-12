package sqlxmodel

import (
	"context"
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

func HasAnyPrefixToken(s string, token ...string) bool {
	i := 0
	for ; i < len(s) && isSpace(s[i]); i++ {
	}
	switch len(token) {
	case 0:
		return false
	case 1:
		return strings.HasPrefix(s[i:], token[0])
	case 2:
		return strings.HasPrefix(s[i:], token[0]) || strings.HasPrefix(s[i:], token[1])
	case 3:
		return strings.HasPrefix(s[i:], token[0]) || strings.HasPrefix(s[i:], token[1]) || strings.HasPrefix(s[i:], token[2])
	}
	for _, tk := range token {
		if strings.HasPrefix(s[i:], tk) {
			return true
		}
	}
	return false
}

func IfClauseAppendWhere(s string) bool {
	i := 0
	for ; i < len(s) && isSpace(s[i]); i++ {
	}
	s = s[i:]
	if len(s) <= 0 {
		return false
	}
	if len(s) < 5 {
		return true
	}
	c := s[0]
	possible := c == 'w' || c == 'o' || c == 'j' || c == 'l' || c == 'r' || c == 'i' || c == 'h'
	if !possible {
		return true
	}
	i = 0
	for ; i < len(s) && !isSpace(s[i]); i++ {
		if i > 6 { // out of having
			return true
		}
	}
	token := s[:i]
	possible = token == "where" || token == "WHERE" ||
		token == "order" || token == "ORDER" ||
		token == "join" || token == "JOIN" ||
		token == "left" || token == "LEFT" ||
		token == "right" || token == "RIGHT" ||
		token == "inner" || token == "INNER" ||
		token == "having" || token == "HAVING"
	return !possible
}

func beforeInsert(ctx context.Context, values interface{}, deepth int) error {
	var err error
	switch vs := values.(type) {
	case interface{ BeforeInsert() }:
		vs.BeforeInsert()
	case interface{ BeforeInsert() error }:
		err = vs.BeforeInsert()
	case interface{ BeforeInsert(ctx context.Context) }:
		vs.BeforeInsert(ctx)
	case interface {
		BeforeInsert(ctx context.Context) error
	}:
		err = vs.BeforeInsert(ctx)
	default:
		if deepth > 1 {
			rv := reflect.Indirect(reflect.ValueOf(vs))
			if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
				for i := rv.Len() - 1; i >= 0; i-- {
					if err = beforeInsert(ctx, rv.Index(i).Interface(), deepth-1); err != nil {
						return err
					}
				}
			}
		}
	}
	return err
}

func BeforeInsert(ctx context.Context, values interface{}) error {
	return beforeInsert(ctx, values, 2)
}

func beforeUpdate(ctx context.Context, values interface{}, deepth int) error {
	var err error
	switch vs := values.(type) {
	case interface{ BeforeUpdate() }:
		vs.BeforeUpdate()
	case interface{ BeforeUpdate() error }:
		err = vs.BeforeUpdate()
	case interface{ BeforeUpdate(ctx context.Context) }:
		vs.BeforeUpdate(ctx)
	case interface {
		BeforeUpdate(ctx context.Context) error
	}:
		err = vs.BeforeUpdate(ctx)
	default:
		if deepth > 1 {
			rv := reflect.Indirect(reflect.ValueOf(vs))
			if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
				for i := rv.Len() - 1; i >= 0; i-- {
					if err = beforeUpdate(ctx, rv.Index(i).Interface(), deepth-1); err != nil {
						return err
					}
				}
			}
		}
	}
	return err
}

func BeforeUpdate(ctx context.Context, values interface{}) error {
	return beforeUpdate(ctx, values, 2)
}

func WithIn(clause string, args []interface{}, offset int) (string, []interface{}) {
	if len(args) <= 0 || len(args) <= offset {
		return clause, args
	}
	if offset < 0 {
		offset = 0
	}
	var nargs []interface{}
	var s strings.Builder
	var off int = -1
	nargs = append(nargs, args[:offset]...)
	args = args[offset:]
	for {
		off = strings.IndexByte(clause, '?')
		if off < 0 {
			s.WriteString(clause)
			nargs = append(nargs, args...)
			break
		}
		if len(args) <= 0 {
			s.WriteString(clause)
			break
		}
		rv := reflect.ValueOf(args[0])
		rt := rv.Type()
		if !(args[0] == nil || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array) {
			s.WriteString(clause[:off+1])
			nargs = append(nargs, args[0])
			clause = clause[off+1:]
			args = args[1:]
			continue
		}
		if args[0] == nil || rv.Len() <= 0 {
			s.WriteString("(NULL)")
		} else {
			s.WriteString(clause[:off])
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
		clause = clause[off+1:]
		args = args[1:]
	}
	return s.String(), nargs
}

func FormatSetClause(s string) string {
	var p int = -1
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			p = i
			break
		}
	}
	if p < 0 {
		return "`" + s + "`=:" + s
	}
	return s
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

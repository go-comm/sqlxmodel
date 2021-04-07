package sqlxmodel

import (
	"context"
	"database/sql"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"text/template"
)

type NamedExecContext interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type SelectContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type GetContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type ExecContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type QueryRowContext interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type QueryRowsContext interface {
	QueryRowsContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type ModelFieldInfo struct {
	FieldName       string
	StructFieldName string
}

type ModelInfo struct {
	Name                  string
	PackageName           string
	PrimaryKey            string
	PrimaryKeyStructField string
	TableName             string
	Fields                []*ModelFieldInfo
}

func NewSqlxModel(tagName string) *SqlxModel {
	return &SqlxModel{
		Mapper: NewReflectMapper(tagName),
		cached: make(map[string]*ModelInfo),
	}
}

type SqlxModel struct {
	mutex  sync.RWMutex
	cached map[string]*ModelInfo
	Mapper *ReflectMapper
}

func (m *SqlxModel) TryModel(e interface{}) *ModelInfo {
	v := reflect.Indirect(reflect.ValueOf(e))
	t := v.Type()
	mi := &ModelInfo{}
	mi.Name = t.Name()
	mi.PackageName = strings.Split(t.String(), ".")[0]
	// tablename
	if tableNamer, ok := e.(interface {
		TableName() string
	}); ok {
		mi.TableName = tableNamer.TableName()
	} else {
		mi.TableName = mi.Name
	}
	// primary key
	possiblePrimaryKey := ""
	if primaryKeyer, ok := e.(interface {
		PrimaryKey() string
	}); ok {
		possiblePrimaryKey = primaryKeyer.PrimaryKey()
	} else {
		log.Panicf("Func %v.PrimaryKey must be defined", mi.Name)
	}

	m.Mapper.TravelFieldsFunc(t, func(fi *FieldInfo) {
		tag := strings.TrimSpace(fi.Tag)
		if fi.Tag == "-" {
			return
		}
		mfi := &ModelFieldInfo{}
		if len(tag) <= 0 {
			tag = fi.Name
		}
		mfi.FieldName = tag
		mfi.StructFieldName = fi.Name
		mi.Fields = append(mi.Fields, mfi)
		if possiblePrimaryKey != "" && possiblePrimaryKey == fi.Tag {
			mi.PrimaryKey = fi.Tag
			mi.PrimaryKeyStructField = mfi.StructFieldName
		}
	})

	if len(mi.PrimaryKey) <= 0 {
		log.Panicf("%v.PrimaryKey can not be empty", mi.Name)
		return mi
	}

	m.mutex.Lock()
	m.cached[mi.Name] = mi
	m.mutex.Unlock()
	return mi
}

func (m *SqlxModel) WriteToFile(e interface{}, path string) error {
	mi := m.TryModel(e)
	tpl, err := template.New("").
		Funcs(template.FuncMap{
			"IsEmpty":        isEmpty,
			"Title":          strings.Title,
			"FormattedField": formattedField,
			"JoinExpr":       joinExpr,
		}).
		Parse(getTpl())
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return tpl.Execute(f, mi)
}

var (
	gShowSQL                           = true
	gSQLPrinter func(v ...interface{}) = log.Println
)

func SetShowSQL(t bool) {
	gShowSQL = t
}

func ShowSQL() bool {
	return gShowSQL
}

func PrintSQL(v ...interface{}) {
	if gShowSQL && gSQLPrinter != nil {
		gSQLPrinter(v...)
	}
}

func WithIn(where string, args ...interface{}) (string, []interface{}) {
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

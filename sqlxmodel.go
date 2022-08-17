package sqlxmodel

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"

	"github.com/go-comm/sqlxmodel/writer"
	"github.com/jmoiron/sqlx"
)

var errInvalidModel = errors.New("model invalid")

type Model interface {
	PrimaryKey() string

	TableName() string

	QueryFirstByPrimaryKey(ctx context.Context, db GetContext, dest interface{}, selection string, pk interface{}) error

	QueryFirst(ctx context.Context, db GetContext, dest interface{}, selection string, whereAndArgs ...interface{}) error

	QueryList(ctx context.Context, db SelectContext, dest interface{}, selection string, whereAndArgs ...interface{}) error

	Update(ctx context.Context, db ExecContext, section string, whereAndArgs ...interface{}) (sql.Result, error)

	NamedUpdate(ctx context.Context, db NamedExecContext, section string, where string, values interface{}) (sql.Result, error)

	NamedUpdateColumns(ctx context.Context, db NamedExecContext, columns []string, where string, values interface{}) (sql.Result, error)

	Insert(ctx context.Context, db NamedExecContext, values interface{}) (sql.Result, error)

	DeleteByPrimaryKey(ctx context.Context, db ExecContext, pk interface{}) (sql.Result, error)

	Delete(ctx context.Context, db ExecContext, whereAndArgs ...interface{}) (sql.Result, error)

	Count(ctx context.Context, db QueryRowContext, whereAndArgs ...interface{}) (int64, error)

	Has(ctx context.Context, db QueryRowContext, whereAndArgs ...interface{}) (bool, error)

	RelatedWith(ctx context.Context, db GetContext, field string, v interface{}) error
}

type QueryContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type QueryRowContext interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type ExecContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type NamedExecContext interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type GetContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type SelectContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type PrepareNamedContext interface {
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

type QueryxContext interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type QueryRowxContext interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type NamedQueryContext interface {
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

type DBContext interface {
	QueryContext
	ExecContext
	NamedExecContext
	GetContext
	SelectContext
	PrepareNamedContext
	QueryRowContext
	QueryxContext
	QueryRowxContext
	NamedQueryContext
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
		if len(tag) <= 0 || tag == "-" {
			return
		}
		mfi := &ModelFieldInfo{}
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

func (m *SqlxModel) WriteTo(w io.Writer, e0 interface{}, e1 ...interface{}) error {
	es := append([]interface{}(nil), e0)
	if len(e1) > 0 {
		es = append(es, e1...)
	}
	tpl := template.New("").
		Funcs(template.FuncMap{
			"IsEmpty":        isEmpty,
			"Title":          strings.Title,
			"LowerTitle":     lowerTitle,
			"FormattedField": formattedField,
			"JoinExpr":       joinExpr,
		})
	var err error
	firstWriteHeader := true
	for _, o := range es {
		mi := m.TryModel(o)
		if firstWriteHeader {
			err = writer.WriteHeader(tpl, w, mi)
			if err != nil {
				return err
			}
		}
		firstWriteHeader = false
		err = writer.WriteBody(tpl, w, mi)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *SqlxModel) WriteToFile(path string, e0 interface{}, e1 ...interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return m.WriteTo(f, e0, e1...)
}

var (
	gShowSQL       int32                  = 0
	gReflectMapper                        = NewReflectMapper("db")
	gSQLPrinter    func(v ...interface{}) = log.Println
)

func SetShowSQL(f bool) {
	if f {
		atomic.StoreInt32(&gShowSQL, 1)
	} else {
		atomic.StoreInt32(&gShowSQL, 0)
	}
}

func ShowSQL() bool {
	return atomic.LoadInt32(&gShowSQL) != 0
}

func PrintSQL(v ...interface{}) {
	if ShowSQL() && gSQLPrinter != nil {
		gSQLPrinter(v...)
	}
}

func Truncate(ctx context.Context, db ExecContext, model interface{}) error {
	if m, ok := model.(interface {
		TableName() string
	}); ok {
		_, err := db.ExecContext(ctx, "Truncate "+m.TableName())
		return err
	}
	return nil
}

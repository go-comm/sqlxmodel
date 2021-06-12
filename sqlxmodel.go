package sqlxmodel

import (
	"context"
	"database/sql"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
)

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

	RelatedWith(ctx context.Context, db GetContext, field string, v interface{}) error
}

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

func (m *SqlxModel) WriteToFile(e interface{}, path string) error {
	mi := m.TryModel(e)
	tpl, err := template.New("").
		Funcs(template.FuncMap{
			"IsEmpty":        isEmpty,
			"Title":          strings.Title,
			"LowerTitle":     lowerTitle,
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

func relatedWith(ctx context.Context, db GetContext, modelRefv reflect.Value, field string, pk interface{}) error {
	rv := reflect.Indirect(modelRefv)
	m := gReflectMapper.TryMap(rv.Type())
	fi, ok := m.Names[field]
	if !ok {
		return nil
	}
	newfv := reflect.New(Deref(fi.Type))
	ifv, ok := newfv.Interface().(interface {
		QueryFirstByPrimaryKey(ctx context.Context, db GetContext, dest interface{}, selection string, pk interface{}) error
	})
	if !ok {
		return nil
	}
	err := ifv.QueryFirstByPrimaryKey(ctx, db, ifv, "", pk)
	if err != nil {
		return err
	}
	fv := FieldByIndex(rv, fi.Index)
	if fv.Kind() == newfv.Kind() {
		fv.Set(newfv)
	} else {
		fv.Set(reflect.Indirect(newfv))
	}
	return nil
}

func RelatedWith(ctx context.Context, db GetContext, model interface{}, field string, pk interface{}) error {
	if pk == nil || reflect.ValueOf(pk).IsZero() {
		return nil
	}
	return relatedWith(ctx, db, reflect.ValueOf(model), field, pk)
}

func relatedWithRef(ctx context.Context, db GetContext, modelRefv reflect.Value, field []string, ref ...string) error {
	if len(ref) <= 0 || len(field) <= 0 {
		return nil
	}
	modelRefv = reflect.Indirect(modelRefv)
	rt := Deref(modelRefv.Type())
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < modelRefv.Len(); i++ {
			if err := relatedWithRef(ctx, db, modelRefv.Index(i), field, ref...); err != nil {
				return err
			}
		}
	case reflect.Struct:
		pk := modelRefv.FieldByName(ref[0])
		if pk.IsZero() {
			return nil
		}
		if err := relatedWith(ctx, db, modelRefv, field[0], pk.Interface()); err != nil {
			return err
		}
		if len(field) >= 1 && len(ref) >= 1 {
			return relatedWithRef(ctx, db, modelRefv.FieldByName(field[0]), field[1:], ref[1:]...)
		}
	default:
	}
	return nil
}

func RelatedWithRef(ctx context.Context, db GetContext, model interface{}, field string, ref ...string) error {
	if len(field) <= 0 || len(ref) <= 0 {
		return nil
	}
	return relatedWithRef(ctx, db, reflect.ValueOf(model), strings.Split(field, "."), ref...)
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

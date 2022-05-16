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

func (m *SqlxModel) WriteTo(e interface{}, w io.Writer) error {
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
	return tpl.Execute(w, mi)
}

func (m *SqlxModel) WriteToFile(e interface{}, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return m.WriteTo(e, f)
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
	store := getStore(ctx, fi.Type)

	if store != nil {
		if vv, ok := store[pk]; ok {
			if vv.NoRow {
				return sql.ErrNoRows
			}
			fv := FieldByIndex(rv, fi.Index)
			if fv.Kind() == vv.Value.Kind() {
				fv.Set(vv.Value)
			} else {
				fv.Set(reflect.Indirect(vv.Value))
			}
			return nil
		}
	}

	newfv := reflect.New(Deref(fi.Type))
	ifv, ok := newfv.Interface().(interface {
		QueryFirstByPrimaryKey(ctx context.Context, db GetContext, dest interface{}, selection string, pk interface{}) error
	})
	if !ok {
		return errInvalidModel
	}
	err := ifv.QueryFirstByPrimaryKey(ctx, db, ifv, "", pk)
	if err != nil {
		if err == sql.ErrNoRows {
			if store != nil {
				vv := &ctxEntryVal{}
				vv.NoRow = true
				store[pk] = vv
			}
			return err
		}
		return err
	}
	if store != nil {
		store[pk] = &ctxEntryVal{false, newfv}
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
	if model == nil {
		return errInvalidModel
	}
	if pk == nil || reflect.ValueOf(pk).IsZero() {
		return sql.ErrNoRows
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
				if err != sql.ErrNoRows {
					return err
				}
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
	if model == nil {
		return errInvalidModel
	}
	if len(field) <= 0 || len(ref) <= 0 {
		return nil
	}
	if !hasCtxEntry(ctx) {
		ctx = WithContext(ctx)
	}
	return relatedWithRef(ctx, db, reflect.ValueOf(model), strings.Split(field, "."), ref...)
}

func GetByPK(ctx context.Context, db GetContext, model interface{}, pk interface{}) error {
	if model == nil {
		return errInvalidModel
	}
	if pk == nil || reflect.ValueOf(pk).IsZero() {
		return sql.ErrNoRows
	}
	rv := reflect.Indirect(reflect.ValueOf(model))
	rt := rv.Type()

	store := getStore(ctx, rt)

	if store != nil {
		if vv, ok := store[pk]; ok {
			if vv.NoRow {
				return sql.ErrNoRows
			}
			if rv.Kind() == vv.Value.Kind() {
				rv.Set(vv.Value)
			} else {
				rv.Set(reflect.Indirect(vv.Value))
			}
			return nil
		}
	}

	m, ok := model.(interface {
		QueryFirstByPrimaryKey(ctx context.Context, db GetContext, dest interface{}, selection string, pk interface{}) error
	})
	if !ok {
		return errInvalidModel
	}
	err := m.QueryFirstByPrimaryKey(ctx, db, m, "", pk)
	if err != nil {
		if err == sql.ErrNoRows {
			if store != nil {
				vv := &ctxEntryVal{}
				vv.NoRow = true
				store[pk] = vv
			}
		}
		return err
	}
	if store != nil {
		store[pk] = &ctxEntryVal{false, rv}
	}
	return nil
}

type ctxKey struct{ Key int }

var ctxEntryKey = &ctxKey{1}

type ctxEntryVal struct {
	NoRow bool
	Value reflect.Value
}

type ctxEntry struct {
	Data map[reflect.Type]map[interface{}]*ctxEntryVal
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxEntryKey, &ctxEntry{})
}

func hasCtxEntry(ctx context.Context) bool {
	ice := ctx.Value(ctxEntryKey)
	if ice == nil {
		return false
	}
	_, ok := ice.(*ctxEntry)
	return ok
}

func getStore(ctx context.Context, tp reflect.Type) map[interface{}]*ctxEntryVal {
	ice := ctx.Value(ctxEntryKey)
	if ice == nil {
		return nil
	}
	ce, ok := ice.(*ctxEntry)
	if !ok {
		return nil
	}
	var b map[interface{}]*ctxEntryVal
	if ce.Data == nil {
		ce.Data = make(map[reflect.Type]map[interface{}]*ctxEntryVal)
		b = make(map[interface{}]*ctxEntryVal)
		ce.Data[tp] = b
	} else {
		b = ce.Data[tp]
		if b == nil {
			b = make(map[interface{}]*ctxEntryVal)
			ce.Data[tp] = b
		}
	}
	return b
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

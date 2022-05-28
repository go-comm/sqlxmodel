package sqlxmodel

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
)

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

package sqlxmodel

import (
	"reflect"
	"sync"
)

type FieldInfo struct {
	Index       int
	Name        string
	Tag         string
	StructField reflect.StructField
}

type StructMap struct {
	Index    []*FieldInfo
	TagNames map[string]*FieldInfo
}

func NewReflectMapper(tagName string) *ReflectMapper {
	return &ReflectMapper{
		tagName: tagName,
		cached:  make(map[reflect.Type]*StructMap),
	}
}

type ReflectMapper struct {
	mutex   sync.RWMutex
	tagName string
	cached  map[reflect.Type]*StructMap
}

func (mapper *ReflectMapper) TryMap(t reflect.Type) *StructMap {
	t = Deref(t)
	mapper.mutex.Lock()
	mapping, ok := mapper.cached[t]
	if !ok {
		mapping = mapper.getMapping(t)
		mapper.cached[t] = mapping
	}
	mapper.mutex.Unlock()
	return mapping
}

func (mapper *ReflectMapper) getMapping(t reflect.Type) *StructMap {
	var fieldinfos []*FieldInfo
	var queue []reflect.Type
	queue = append(queue, Deref(t))

	for len(queue) > 0 {
		t = queue[0]
		queue = queue[1:]

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Anonymous {
				queue = append(queue, Deref(f.Type))
				continue
			}
			tagVal := f.Tag.Get(mapper.tagName)
			fi := &FieldInfo{
				Index:       i,
				Name:        f.Name,
				Tag:         tagVal,
				StructField: f,
			}
			fieldinfos = append(fieldinfos, fi)
		}

	}

	m := &StructMap{Index: fieldinfos, TagNames: map[string]*FieldInfo{}}
	for i := 0; i < len(m.Index); i++ {
		f := m.Index[i]
		m.TagNames[f.Tag] = f
	}
	return m
}

func (mapper *ReflectMapper) TravelFieldFunc(t reflect.Type, fn func(fi *FieldInfo)) {
	m := mapper.TryMap(t)
	for _, v := range m.Index {
		fn(v)
	}
}

func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

package sqlxmodel

import (
	"reflect"
	"sync"
)

type FieldInfo struct {
	Index     []int
	Name      string
	Tag       string
	Anonymous bool
	Type      reflect.Type
	StructTag reflect.StructTag
}

type StructMap struct {
	Index []*FieldInfo
	Tags  map[string]*FieldInfo
	Names map[string]*FieldInfo
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
	mapper.mutex.RLock()
	mapping, ok := mapper.cached[t]
	mapper.mutex.RUnlock()
	if !ok {
		mapper.mutex.Lock()
		if !ok {
			mapping = mapper.getMapping(t)
			mapper.cached[t] = mapping
		}
		mapper.mutex.Unlock()
	}
	return mapping
}

func (mapper *ReflectMapper) getMapping(t reflect.Type) *StructMap {
	var fieldinfos []*FieldInfo
	var queue []*FieldInfo
	var head *FieldInfo
	queue = append(queue, &FieldInfo{Type: Deref(t)})

	for len(queue) > 0 {
		head = queue[0]
		queue = queue[1:]
		t = Deref(head.Type)
		if t.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			indexes := append([]int(nil), head.Index...)
			indexes = append(indexes, i)

			fi := &FieldInfo{
				Name:      f.Name,
				Index:     indexes,
				Type:      f.Type,
				StructTag: f.Tag,
				Anonymous: f.Anonymous,
			}

			if len(mapper.tagName) > 0 {
				fi.Tag = f.Tag.Get(mapper.tagName)
			}

			if f.Anonymous {
				queue = append(queue, fi)
				continue
			}
			fieldinfos = append(fieldinfos, fi)
		}
	}

	m := &StructMap{Index: fieldinfos, Names: map[string]*FieldInfo{}, Tags: map[string]*FieldInfo{}}
	for i := 0; i < len(m.Index); i++ {
		f := m.Index[i]
		m.Names[f.Name] = f
		if len(f.Tag) > 0 {
			m.Tags[f.Tag] = f
		}
	}
	return m
}

func (mapper *ReflectMapper) TravelFieldsFunc(t reflect.Type, fn func(*FieldInfo)) {
	m := mapper.TryMap(t)
	for _, v := range m.Index {
		fn(v)
	}
}

func (mapper *ReflectMapper) TravelFieldsByTagsFunc(t reflect.Type, fn func(*FieldInfo)) {
	m := mapper.TryMap(t)
	for _, v := range m.Tags {
		fn(v)
	}
}

func (mapper *ReflectMapper) FieldByName(v reflect.Value, name string) (reflect.Value, bool) {
	m := mapper.TryMap(v.Type())
	fi, ok := m.Names[name]
	if !ok {
		return reflect.Value{}, false
	}
	return FieldByIndex(v, fi.Index), true
}

func (mapper *ReflectMapper) FieldByTag(v reflect.Value, tag string) (reflect.Value, bool) {
	m := mapper.TryMap(v.Type())
	fi, ok := m.Tags[tag]
	if !ok {
		return reflect.Value{}, false
	}
	return FieldByIndex(v, fi.Index), true
}

func (mapper *ReflectMapper) TraversalsByNamesFunc(t reflect.Type, names []string, fn func(*FieldInfo)) {
	m := mapper.TryMap(t)
	for _, name := range names {
		fi, ok := m.Names[name]
		if ok {
			fn(fi)
		}
	}
}

func (mapper *ReflectMapper) TraversalsByNames(t reflect.Type, names []string) (ls []*FieldInfo) {
	mapper.TraversalsByNamesFunc(t, names, func(fi *FieldInfo) {
		ls = append(ls, fi)
	})
	return
}

func (mapper *ReflectMapper) TraversalsByName(t reflect.Type, name string) (*FieldInfo, bool) {
	m := mapper.TryMap(t)
	fi, ok := m.Names[name]
	return fi, ok
}

func (mapper *ReflectMapper) TraversalsByTagsFunc(t reflect.Type, tags []string, fn func(*FieldInfo)) {
	m := mapper.TryMap(t)
	for _, name := range tags {
		fi, ok := m.Tags[name]
		if ok {
			fn(fi)
		}
	}
}

func (mapper *ReflectMapper) TraversalsByTags(t reflect.Type, tags []string) (ls []*FieldInfo) {
	mapper.TraversalsByNamesFunc(t, tags, func(fi *FieldInfo) {
		ls = append(ls, fi)
	})
	return
}

func (mapper *ReflectMapper) TraversalsByTag(t reflect.Type, tag string) (*FieldInfo, bool) {
	m := mapper.TryMap(t)
	fi, ok := m.Tags[tag]
	return fi, ok
}

func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func FieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		v = reflect.Indirect(v).Field(i)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			alloc := reflect.New(Deref(v.Type()))
			v.Set(alloc)
		}
		if v.Kind() == reflect.Map && v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
	}
	return v
}

func FieldByIndexReadOnly(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		v = reflect.Indirect(v).Field(i)
	}
	return v
}

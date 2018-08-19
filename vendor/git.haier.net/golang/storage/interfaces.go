package storage

import (
	"context"
	"fmt"
	"reflect"
)

type Interface interface {
	Get(ctx context.Context, key string, out interface{}) error
	Create(ctx context.Context, key string, obj interface{}, ttl uint64) error
	BulkCreate(ctx context.Context, key string, c chan ChannelObj, ttl uint64) error
	Delete(ctx context.Context, key string, out interface{}) error
	DeleteByQuery(ctx context.Context, key string, keyword interface{}) (deleted, conflict int64, err error)
	List(ctx context.Context, key string, sp *SelectionPredicate, obj interface{}) ([]interface{}, error)
	Update(ctx context.Context, key string, resourceVersion int64, obj interface{}, ttl uint64) error
	Upsert(ctx context.Context, key string, resourceVersion int64, update_obj, insert_obj interface{}, ttl uint64) error
}

type ChannelObj struct {
	Data interface{}
	Id   string
}

func SetResourceVersion(out interface{}, v int64) {
	switch reflect.TypeOf(out).Kind() {
	case reflect.Ptr:
		v_typ := reflect.TypeOf(out).Elem()
		if v_typ.Kind() == reflect.Map {
			reflect.ValueOf(out).Elem().SetMapIndex(reflect.ValueOf("_version"), reflect.ValueOf(v))
		} else if v_typ.Kind() == reflect.Struct {
			_, ok := reflect.TypeOf(out).Elem().FieldByName("ResourceVersion")
			if ok {
				version_v := reflect.ValueOf(out).Elem().FieldByName("ResourceVersion")
				if version_v.Kind() == reflect.Int64 || version_v.Kind() == reflect.Int {
					version_v.SetInt(v)
				} else if version_v.Kind() == reflect.String {
					version_v.SetString(fmt.Sprintf("%d", v))
				}
			} else {
				//fmt.Printf("not found ResourceVersion\n")
			}
		} else {
			//fmt.Printf("unknown kind: %s\n", v_typ.Kind())
		}
	case reflect.Map:
		reflect.ValueOf(out).SetMapIndex(reflect.ValueOf("_version"), reflect.ValueOf(v))
	default:

	}

}

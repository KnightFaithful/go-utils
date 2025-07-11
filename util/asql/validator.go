package asql

import (
	"example.com/m/util/utilerror"
	"reflect"
)

func checkItemListType(itemList interface{}) *utilerror.UtilError {
	t := reflect.TypeOf(itemList)
	if t.Kind() != reflect.Slice {
		return utilerror.NewError("param is not sliceType pointer")
	}

	if t.Elem().Kind() != reflect.Ptr {
		return utilerror.NewError("param is not pointer")
	}

	if t.Elem().Elem().Kind() != reflect.Struct {
		return utilerror.NewError("param is not slice[*struct] pointer")
	}

	return nil
}

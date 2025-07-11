package helper

import (
	"example.com/m/util/fileutil"
	"example.com/m/util/stringutil"
	"example.com/m/util/utilerror"
	"reflect"
)

const PageMaxLimit = int64(10 * 1024 * 1024)

func CalculateRepositoryEveryPageSize(list interface{}) (int64, *utilerror.UtilError) {
	ref := reflect.ValueOf(list)
	kind := ref.Type().Kind()
	if kind == reflect.Ptr {
		ref = ref.Elem()
	}
	kind = ref.Type().Kind()
	if kind != reflect.Slice {
		return 0, utilerror.NewError("PageList kind not reflect.Slice")
	}
	len := int64(ref.Len())
	if len == 0 {
		return 0, nil
	}
	size := fileutil.GetStringMemorySize(stringutil.Object2String(ref.Index(0).Interface()))
	size = size * len
	pageSize := PageMaxLimit*len/size + 1
	return pageSize, nil
}

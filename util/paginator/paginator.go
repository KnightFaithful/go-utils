package paginator

import (
	"example.com/m/util/expression"
	"example.com/m/util/utilerror"
	"gorm.io/gorm"
	"reflect"
)

type PageIn struct {
	Pageno     int64  // 页码
	Count      int64  // 数量
	OrderBy    string // 为空字符串 代表不需要排序
	IsGetTotal bool   // 为False代表不需要获取总数
}

var LimitOnePage = &PageIn{Pageno: 1, Count: 1}

func Paginator(qs *gorm.DB, pageIn *PageIn, out interface{}) (int64, *utilerror.UtilError) {
	/*
		如果pageIn 为nil，则不需要分页，当pageIn 不为nil时，pageno和count都大于0时才会用到offset和limit
	*/
	if pageIn == nil {
		err := qs.Find(out).Error
		if err != nil {
			return 0, utilerror.NewError(err.Error())
		}
		return 0, nil
	}
	total := int64(0)
	if pageIn.IsGetTotal {
		err := qs.Count(&total).Error
		if err != nil {
			return 0, utilerror.NewError(err.Error())
		}
	}
	if pageIn.OrderBy != "" {
		qs = qs.Order(pageIn.OrderBy)
	}
	if pageIn.Pageno > 0 && pageIn.Count > 0 {
		qs = qs.Offset(int((pageIn.Pageno - 1) * pageIn.Count)).Limit(int(pageIn.Count))
	}
	err := qs.Find(out).Error
	if err != nil {
		return 0, utilerror.NewError(err.Error())
	}
	return total, nil
}

// 内存分页
func PageList(l interface{}, page, count int64) {
	if l == nil || page <= 0 {
		return
	}
	kind := reflect.ValueOf(l).Type().Kind()
	if kind != reflect.Ptr {
		panic("PageList kind not reflect.Ptr")
	}
	kind = reflect.ValueOf(l).Elem().Type().Kind()
	if kind != reflect.Slice {
		panic("PageList kind not reflect.Slice")
	}

	total := reflect.ValueOf(l).Elem().Len()

	if total == 0 {
		return
	}

	begin := int((page - 1) * count)
	end := int(page * count)

	if total > end {
		l2 := reflect.ValueOf(l).Elem().Slice(begin, end)
		reflect.ValueOf(l).Elem().Set(l2)
	} else if total > begin {
		l2 := reflect.ValueOf(l).Elem().Slice(begin, total)
		reflect.ValueOf(l).Elem().Set(l2)
	} else {
		l2 := reflect.ValueOf(l).Elem().Slice(0, 0)
		reflect.ValueOf(l).Elem().Set(l2)
	}
}

func PageListAndReturn(l interface{}, page, count int64) interface{} {
	if l == nil || page <= 0 {
		return nil
	}
	ref := reflect.ValueOf(l)
	kind := ref.Type().Kind()
	if kind == reflect.Ptr {
		ref = ref.Elem()
	}
	kind = ref.Type().Kind()
	if kind != reflect.Slice {
		panic("PageList kind not reflect.Slice")
	}

	total := ref.Len()

	if total == 0 {
		return nil
	}

	begin := int((page - 1) * count)
	end := int(page * count)
	if total > end {
		slice := reflect.ValueOf(l).Slice(begin, end)
		// 将切片转换为与原始数组类型相同的接口类型
		result := reflect.MakeSlice(reflect.TypeOf(l), slice.Len(), slice.Len())
		reflect.Copy(result, slice)
		return result.Interface()
	} else if total > begin {
		slice := reflect.ValueOf(l).Slice(begin, total)
		// 将切片转换为与原始数组类型相同的接口类型
		result := reflect.MakeSlice(reflect.TypeOf(l), slice.Len(), slice.Len())
		reflect.Copy(result, slice)
		return result.Interface()
	} else {
		slice := reflect.ValueOf(l).Slice(0, 0)
		// 将切片转换为与原始数组类型相同的接口类型
		result := reflect.MakeSlice(reflect.TypeOf(l), slice.Len(), slice.Len())
		reflect.Copy(result, slice)
		return result.Interface()
	}
}

// 面向基础类型分页
func PagingBaseType(l interface{}, page, count int64) interface{} {

	switch t := l.(type) {
	case []int64:
		totalCount := int64(len(l.([]int64)))
		return t[(page-1)*count : expression.IfInt64(page*count < totalCount, page*count, totalCount)]
	case []string:
		totalCount := int64(len(l.([]string)))
		return t[(page-1)*count : expression.IfInt64(page*count < totalCount, page*count, totalCount)]
	case []float64:
		totalCount := int64(len(l.([]float64)))
		return t[(page-1)*count : expression.IfInt64(page*count < totalCount, page*count, totalCount)]
	default:
		return pageSlice(l, page, count)
	}
}

// support []Struct{}, []int64{}...
func pageSlice(l interface{}, page int64, count int64) interface{} {
	kind := reflect.ValueOf(l).Type().Kind()
	if kind != reflect.Slice {
		panic("PageList kind not reflect.Slice")
	}
	total := reflect.ValueOf(l).Len()
	begin := int((page - 1) * count)
	end := int(page * count)
	if total > end {
		return reflect.ValueOf(l).Slice(begin, end).Interface()
	} else if total > begin {
		return reflect.ValueOf(l).Slice(begin, total).Interface()
	} else {
		return reflect.ValueOf(l).Slice(0, 0).Interface()
	}
}

const (
	PageNoDefault    int64 = 1
	PageCountDefault int64 = 20
)

func GetValidPageInfo(pageNo, count int64) (int64, int64) {
	if pageNo == 0 {
		pageNo = PageNoDefault
	}
	if count == 0 {
		count = PageCountDefault
	}
	return pageNo, count
}

func SplitInt64OutputOffset(total int64, pageSize int64) []int64 {
	if total <= pageSize {
		return []int64{0}
	}
	var res []int64
	for i := int64(0); i <= total; i += pageSize {
		res = append(res, i)
	}
	return res
}

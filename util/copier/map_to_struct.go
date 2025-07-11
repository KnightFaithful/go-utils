package copier

import (
	"example.com/m/util/utilerror"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func SetStructFieldByMap(ptr interface{}, fields map[string]interface{}) (interface{}, *utilerror.UtilError) {

	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return nil, utilerror.NewError("SetStructFieldByMap attr ptr is not point")
	}
	v := reflect.ValueOf(ptr).Elem()

	for i := 0; i < v.NumField(); i++ {

		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("json")

		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		name = strings.Split(name, ",")[0]

		if value, ok := fields[name]; ok {
			valueKind := reflect.ValueOf(value).Kind()
			if reflect.ValueOf(value).Kind() == v.FieldByName(fieldInfo.Name).Kind() {
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			} else if v.FieldByName(fieldInfo.Name).Kind() == reflect.Int64 && valueKind == reflect.String {
				v2, _ := strconv.Atoi(value.(string))
				v3 := int64(v2)
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(v3))
			} else if v.FieldByName(fieldInfo.Name).Kind() == reflect.Slice && valueKind == reflect.String {
				v2, _ := reflect.TypeOf(ptr).Elem().FieldByName(fieldInfo.Name)
				valueList := strings.Split(value.(string), ",")
				if v2.Type.Elem().Kind() == reflect.String {
					var value2List []string
					value2List = append(value2List, valueList...)
					v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value2List))
				} else if v2.Type.Elem().Kind() == reflect.Int64 {
					var value2List []int64
					for _, x := range valueList {
						x2, err := strconv.Atoi(x)
						if err != nil {
							return nil, utilerror.NewError(err.Error())
						}
						value2List = append(value2List, int64(x2))
					}
					v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value2List))
				}
			} else if v.FieldByName(fieldInfo.Name).Kind() == reflect.Ptr && valueKind == reflect.String {

				// TODO 还需要判断指针对应的类型是string还是int
				fmt.Println(v.FieldByName(fieldInfo.Name).Kind())

				field := v.FieldByName(fieldInfo.Name).Type()

				fmt.Println(field)
				v2 := value.(string)
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(&v2))
			}
		}
	}
	return v, nil
}

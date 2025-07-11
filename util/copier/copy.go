package copier

import (
	"encoding/json"
	"example.com/m/util/utilerror"
	"reflect"
	"strings"
)

/*
*
1. dest 必须为指针引用，否则 Unmarshal 会失败
2. 错误由 json 包 Marshal 和 Unmarshal 返回即可
*/
func Copy(source interface{}, dest interface{}) *utilerror.UtilError {
	err := copy(source, dest)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

func copy(source interface{}, dest interface{}) error {

	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(sourceBytes, dest)
	if err != nil {
		return err
	}
	return nil
}

func NewCopyWithIgnore(source interface{}, dest interface{}, ignoreField []string) *utilerror.UtilError {
	err := CopyWithIgnore(source, dest, ignoreField)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

// Deprecated
func CopyWithIgnore(source interface{}, dest interface{}, ignoreField []string) error {
	var srcMap map[string]interface{}
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(sourceBytes, &srcMap)
	if err != nil {
		return err
	}

	for _, field := range ignoreField {
		delete(srcMap, field)
	}

	err = Copy(srcMap, dest)
	if err != nil {
		return err
	}
	return nil
}

// source不变，修改dest
// source中的属性优先，覆盖dest中的属性
func MergeWithIgnore(source interface{}, dest interface{}, ignoreField []string) error {
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	var srcMap map[string]interface{}
	err = json.Unmarshal(sourceBytes, &srcMap)
	if err != nil {
		return err
	}

	for _, field := range ignoreField {
		delete(srcMap, field)
	}

	destBytes, err := json.Marshal(dest)
	if err != nil {
		return err
	}
	var destMap map[string]interface{}
	err = json.Unmarshal(destBytes, &destMap)
	if err != nil {
		return err
	}

	for key, field := range destMap {
		if _, ok := srcMap[key]; !ok {
			srcMap[key] = field
		}
	}

	err = Copy(srcMap, dest)
	if err != nil {
		return err
	}
	return nil
}

func Merge(source interface{}, dest interface{}) error {
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	var srcMap map[string]interface{}
	err = json.Unmarshal(sourceBytes, &srcMap)
	if err != nil {
		return err
	}

	destBytes, err := json.Marshal(dest)
	if err != nil {
		return err
	}
	var destMap map[string]interface{}
	err = json.Unmarshal(destBytes, &destMap)
	if err != nil {
		return err
	}

	for key, field := range destMap {
		if _, ok := srcMap[key]; !ok {
			srcMap[key] = field
		}
	}

	err = Copy(srcMap, dest)
	if err != nil {
		return err
	}
	return nil
}

// 根据json标签，生成map
// json序列化的方式
func CopyToJsonMap(source interface{}) (map[string]interface{}, error) {
	var srcMap map[string]interface{}
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(sourceBytes, &srcMap)
	if err != nil {
		return nil, err
	}
	return srcMap, nil
}

// 根据gorm标签，生成map
// 使用struct反射的方式
func CopyToDbMap(source interface{}) (map[string]interface{}, error) {
	srcMap := map[string]interface{}{}
	if source == nil {
		return nil, utilerror.NewError("source can not be nil")
	}
	fromType := reflect.TypeOf(source)
	if fromType.Kind() == reflect.Ptr {
		fromType = fromType.Elem()
	}
	fromValue := reflect.ValueOf(source)
	if fromValue.Kind() == reflect.Ptr {
		fromValue = fromValue.Elem()
	}

	for i := 0; i < fromType.NumField(); i++ {
		// 判断是否为可导出字段
		if !fromValue.Field(i).CanInterface() {
			continue
		}
		tag := fromType.Field(i).Tag.Get("gorm")
		tag = strings.Split(strings.Split(strings.Split(tag, ",")[0], ":")[1], ";")[0]
		fromValue := fromValue.Field(i).Interface()
		srcMap[tag] = fromValue
	}
	return srcMap, nil
}

func JsonCopy(source interface{}, dest interface{}) *utilerror.UtilError {
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	err = json.Unmarshal(sourceBytes, dest)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

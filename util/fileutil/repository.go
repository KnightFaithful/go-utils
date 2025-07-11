package fileutil

import (
	"encoding/json"
	"example.com/m/util/expression"
	"example.com/m/util/printhelper"
	"example.com/m/util/utilerror"

	"os"
	"reflect"
	"strings"
)

//func Save(target interface{}, path, fileName string) *utilerror.UtilError {
//	if !strings.Contains(fileName, ".") {
//		return utilerror.NewError("file name not correct")
//	}
//	if target == nil {
//		return utilerror.NewError("List cannot be null")
//	}
//	kind := reflect.ValueOf(target).Type().Kind()
//	if kind != reflect.Ptr {
//		return utilerror.NewError("PageList kind not reflect.Ptr")
//	}
//	kind = reflect.ValueOf(target).Elem().Type().Kind()
//	if kind != reflect.Slice {
//		return utilerror.NewError("PageList kind not reflect.Slice")
//	}
//	total := reflect.ValueOf(target).Elem().Len()
//	if total == 0 {
//		return nil
//	}
//	err := fileutil.WriteString(stringutil.Object2String(target), path+fileName)
//	if err != nil {
//		return err.Mark()
//	}
//
//	return nil
//}
//
//func read(target interface{}, path, fileName string) *utilerror.UtilError {
//	if target == nil {
//		return utilerror.NewError("List cannot be null")
//	}
//	kind := reflect.ValueOf(target).Type().Kind()
//	if kind != reflect.Ptr {
//		return utilerror.NewError("PageList kind not reflect.Ptr")
//	}
//	exists := fileutil.FileExists(path + fileName)
//	if !exists {
//		return utilerror.NewError("file not exist")
//	}
//	res, err := fileutil.ReadFile(path + fileName)
//	if err != nil {
//		return err.Mark()
//	}
//	err2 := json.Unmarshal(res, target)
//	if err2 != nil {
//		return utilerror.NewError(err2.Error())
//	}
//	return nil
//}

func SaveJson(target interface{}, path, fileName string) *utilerror.UtilError {
	if !strings.Contains(fileName, ".") {
		return utilerror.NewError("file name not correct")
	}
	if target == nil {
		return utilerror.NewError("List cannot be null")
	}
	if err := CreateFolderIfNotExist(path); err != nil {
		return err.Mark()
	}
	ref := reflect.ValueOf(target)
	kind := ref.Type().Kind()
	if kind == reflect.Ptr {
		ref = ref.Elem()
	}
	kind = ref.Type().Kind()
	if kind != reflect.Slice {
		return utilerror.NewError("PageList kind not reflect.Slice")
	}
	total := ref.Len()
	if total == 0 {
		return nil
	}
	outputFile, err := os.Create(path + fileName)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	defer outputFile.Close()

	// 创建JSON编码器
	encoder := json.NewEncoder(outputFile)

	// 将数据结构序列化为JSON并写入文件
	err = encoder.Encode(target)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

func ReadFolderAllJsonReturnMap(path string, do func(filePath string) (interface{}, *utilerror.UtilError)) ([]interface{}, *utilerror.UtilError) {
	files, err := ReadFolder(expression.IfString(path == "", ".", path))
	if err != nil {
		return nil, err.Mark()
	}
	var res []interface{}
	for _, file := range files {
		printhelper.Println("reading " + file)
		temp, err := do(file)
		if err != nil {
			return nil, err.Mark()
		}
		res = append(res, temp)
	}
	return res, nil
}

func ReadJson(target interface{}, filePathName string) *utilerror.UtilError {
	if !strings.HasSuffix(filePathName, ".json") {
		return nil
	}
	if target == nil {
		return utilerror.NewError("List cannot be null")
	}
	kind := reflect.ValueOf(target).Type().Kind()
	if kind != reflect.Ptr {
		return utilerror.NewError("PageList kind not reflect.Ptr")
	}
	exists := FileExists(filePathName)
	if !exists {
		return utilerror.NewError("file not exist")
	}
	file, err := os.Open(filePathName)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	defer file.Close()

	// 创建JSON解码器
	decoder := json.NewDecoder(file)

	// 解码JSON数据到slice中
	err = decoder.Decode(target)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}

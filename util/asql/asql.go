package asql

import (
	"bytes"
	"example.com/m/util/collection"
	"example.com/m/util/iters"
	"example.com/m/util/utilerror"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GenerateBatchUpdateSQL(tableName string, itemList interface{}) ([]string, *utilerror.UtilError) {

	errCheck := checkItemListType(itemList)
	if errCheck != nil {
		return nil, errCheck.Mark()
	}

	fieldValue := reflect.ValueOf(itemList)
	fieldType := reflect.TypeOf(itemList).Elem().Elem()
	sliceLength := fieldValue.Len()
	fieldNum := fieldType.NumField()

	// 检验结构体标签是否为空和重复
	verifyTagDuplicate := make(map[string]string)
	count := 0
	for i := 0; i < fieldNum; i++ {
		fieldTag := fieldType.Field(i).Tag.Get("gorm")

		fieldName := getFieldName(fieldTag)
		if len(strings.TrimSpace(fieldName)) == 0 {
			return nil, utilerror.NewError("the structure attribute should have tag")
		}

		if strings.HasPrefix(fieldName, "id;") {
			count++
		}

		_, ok := verifyTagDuplicate[fieldName]
		if !ok {
			verifyTagDuplicate[fieldName] = fieldName
		} else {
			return nil, utilerror.NewError("the structure attribute %v tag is not allow duplication", fieldName)
		}
	}

	if count != 1 {
		return nil, utilerror.NewError("the structure attribute should have a primary key")
	}

	IDSet := collection.NewStringSet()
	var IDList []string
	updateMap := make(map[string][]*string)
	for i := 0; i < sliceLength; i++ {
		// 得到某一个具体的结构体的
		structValue := fieldValue.Index(i).Elem()
		for j := 0; j < fieldNum; j++ {
			elem := structValue.Field(j)
			gormTag := fieldType.Field(j).Tag.Get("gorm")
			fieldTag := getFieldName(gormTag)

			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				updateMap[fieldTag] = append(updateMap[fieldTag], nil) // 如果为nil的指针，则填入nil，保持每个field中数组数量的一致性
			} else {
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				var temp string
				switch elem.Kind() {
				case reflect.Int64:
					temp = strconv.FormatInt(elem.Int(), 10)
				case reflect.String:
					if strings.Contains(elem.String(), "'") {
						temp = fmt.Sprintf("'%v'", strings.ReplaceAll(elem.String(), "'", "\\'"))
					} else {
						temp = fmt.Sprintf("'%v'", elem.String())
					}
				case reflect.Float64:
					temp = strconv.FormatFloat(elem.Float(), 'f', -1, 64)
				case reflect.Bool:
					temp = strconv.FormatBool(elem.Bool())
				default:
					return nil, utilerror.NewError("type conversion error, param is %v", fieldType.Field(j).Tag.Get("json"))
				}

				if strings.HasPrefix(fieldTag, "id;") {
					id, err := strconv.ParseInt(temp, 10, 64)
					if err != nil {
						return nil, utilerror.NewError(err.Error())
					}
					// id 的合法性校验
					if id < 1 {
						return nil, utilerror.NewError("this structure should have a primary key and gt 0")
					}
					if IDSet.Contains(temp) {
						return nil, utilerror.NewError("this structure data id can not repeat: %v", temp)
					}
					IDSet.Add(temp)
					IDList = append(IDList, temp)
					continue
				}
				updateMap[fieldTag] = append(updateMap[fieldTag], &temp)
			}

		}
	}
	// 过滤掉 updateMap 中都是 nil 的字段，不用更新
	for fieldTag, valList := range updateMap {
		isFieldAllNil := true
		for _, val := range valList {
			if val != nil {
				isFieldAllNil = false
				break
			}
		}
		if isFieldAllNil { // 如果全为nil的，废弃掉
			delete(updateMap, fieldTag)
		}
	}

	var newIDList []string
	iters.From(IDList).Select(func(i interface{}) interface{} {
		return i.(string)
	}).Distinct().ToSlice(&newIDList)

	if len(IDList) != len(newIDList) {
		var repeatedIDList []string
		iters.From(IDList).Except(iters.From(newIDList)).ToSlice(&repeatedIDList)
		return nil, utilerror.NewError("this structure data id %v can not repeat", strings.Join(repeatedIDList, ","))
	}

	length := len(IDList)
	size := batchCreateMaxNum
	SQLQuantity := getSQLQuantity(length, size)
	var SQLArray []string
	k := 0

	updateFieldCount := len(updateMap)
	for i := 0; i < SQLQuantity; i++ {
		count := 0

		var record bytes.Buffer
		record.WriteString("UPDATE " + tableName + " SET ")

		for fieldName, fieldValueList := range updateMap {
			record.WriteString(fieldName)
			record.WriteString(" = CASE " + "id")

			for j := k; j < len(IDList) && j < len(fieldValueList) && j < size+k; j++ {
				if fieldValueList[j] == nil { // 如果要更新的值为nil，说明不用更新，设置为它自己（这里不能conine，避出现）
					record.WriteString(" WHEN " + IDList[j] + " THEN " + fieldName)
				} else {
					record.WriteString(" WHEN " + IDList[j] + " THEN " + *fieldValueList[j])
				}
			}
			record.WriteString(" ELSE " + fieldName) // 如果没变更则设置成db原来字段的值
			count++
			if count != updateFieldCount {
				record.WriteString(" END, ")
			}
		}
		record.WriteString(" END WHERE ")
		record.WriteString("id" + " IN (")
		min := size + k
		if len(IDList) < min {
			min = len(IDList)
		}
		record.WriteString(strings.Join(IDList[k:min], ","))
		record.WriteString(");")

		k += size
		SQLArray = append(SQLArray, record.String())
	}

	return SQLArray, nil
}

func GenerateInsertSQL(tableName string, itemList interface{}) (string, *utilerror.UtilError) {

	errCheck := checkItemListType(itemList)
	if errCheck != nil {
		return "", errCheck.Mark()
	}

	fieldValue := reflect.ValueOf(itemList)
	fieldType := reflect.TypeOf(itemList).Elem().Elem()
	sliceLength := fieldValue.Len()
	fieldNum := fieldType.NumField()

	var record bytes.Buffer
	record.WriteString("INSERT INTO" + " " + tableName + " " + "(")
	var tempRecord bytes.Buffer
	// 防止gorm中没有ID标签
	filterIDIndex := -1
	verifyTagDuplicate := make(map[string]string)
	for i := 0; i < fieldNum; i++ {
		gormTag := fieldType.Field(i).Tag.Get("gorm")

		fieldName := getFieldName(gormTag)
		if len(strings.TrimSpace(fieldName)) == 0 {
			return "", utilerror.NewError("the structure attribute should have tag")
		}

		_, ok := verifyTagDuplicate[fieldName]
		if !ok {
			verifyTagDuplicate[fieldName] = fieldName
		} else {
			return "", utilerror.NewError("the structure attribute %v tag is not allow duplication", fieldName)
		}

		if strings.HasPrefix(fieldName, "id;") {
			filterIDIndex = i
			continue
		}

		tempRecord.WriteString(fieldName)
		tempRecord.WriteString(",")
	}

	s := strings.TrimRight(tempRecord.String(), ",")
	record.WriteString(s)
	record.WriteString(")" + " " + "VALUES" + " ")
	for i := 0; i < sliceLength; i++ {
		record.WriteString("(")
		var tempSQLRecord bytes.Buffer
		// 得到某一个具体的结构体的
		structValue := fieldValue.Index(i).Elem()

		for j := 0; j < fieldNum; j++ {
			if filterIDIndex != -1 && j == filterIDIndex {
				continue
			}

			elem := structValue.Field(j)

			var temp string
			switch elem.Kind() {
			case reflect.Int64:
				temp = strconv.FormatInt(elem.Int(), 10)
			case reflect.String:
				if strings.Contains(elem.String(), "'") {
					temp = fmt.Sprintf("'%v'", strings.ReplaceAll(elem.String(), "'", "\\'"))
				} else {
					temp = fmt.Sprintf("'%v'", elem.String())
				}
			case reflect.Float64:
				temp = strconv.FormatFloat(elem.Float(), 'f', -1, 64)
			case reflect.Int:
				temp = strconv.Itoa(elem.Interface().(int))
			case reflect.Int8:
				temp = strconv.FormatInt(int64(elem.Interface().(int8)), 10)
			case reflect.Int16:
				temp = strconv.FormatInt(int64(elem.Interface().(int16)), 10)
			case reflect.Int32:
				temp = strconv.FormatInt(int64(elem.Interface().(int32)), 10)
			case reflect.Uint:
				temp = strconv.Itoa(int(elem.Interface().(uint)))
			case reflect.Uint8:
				temp = strconv.FormatUint(elem.Uint(), 10)
			case reflect.Uint16:
				temp = strconv.FormatUint(uint64(elem.Interface().(uint16)), 10)
			case reflect.Uint32:
				temp = strconv.FormatUint(elem.Uint(), 10)
			case reflect.Uint64:
				temp = strconv.FormatUint(elem.Uint(), 10)
			case reflect.Float32:
				temp = fmt.Sprint(elem.Interface())
			case reflect.Bool:
				temp = strconv.FormatBool(elem.Bool())
			default:
				return "", utilerror.NewError("type conversion error, param is %v", fieldType.Field(j).Tag.Get("json"))
			}
			tempSQLRecord.WriteString(temp)
			tempSQLRecord.WriteString(",")
		}
		str := strings.TrimRight(tempSQLRecord.String(), ",")
		record.WriteString(str)
		record.WriteString(")")
		if i != sliceLength-1 {
			record.WriteString(",")
		}
	}
	record.WriteString(";")
	return record.String(), nil
}

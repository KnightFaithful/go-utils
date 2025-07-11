package excel

import (
	"bytes"
	"example.com/m/util/collection"
	"example.com/m/util/convert"
	"example.com/m/util/utilerror"
	"github.com/xuri/excelize/v2"

	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ExcelSheetTab struct {
	SheetName     string      `gorm:"column:sheet_name" json:"sheet_name"`
	Data          interface{} `gorm:"column:data" json:"data"`
	ExcludeTitles []string    `gorm:"-" json:"exclude_titles"`
}

type GenerateExcelOpt struct {
	IsEncrypt     *bool
	MergeCellList []string
}

const CommonEnumYES = "YES"
const EncryptChar = "*"

func GenerateExcelBytes(fileName string, exportDataList []*ExcelSheetTab) ([]byte, string, *utilerror.UtilError) {
	fileBytes, filePath, gErr := GenerateStreamExcelWithMultipleSheet(exportDataList, fileName, &GenerateExcelOpt{
		IsEncrypt: convert.Bool(false),
	})
	if gErr != nil {
		return nil, "", gErr.Mark()
	}
	return fileBytes, filePath, nil
}

func GenerateStreamExcelWithMultipleSheet(sheetList []*ExcelSheetTab, excelName string, opt *GenerateExcelOpt) ([]byte, string, *utilerror.UtilError) {
	excel, err := generateStreamExcelSheetList(sheetList, opt)
	if err != nil {
		return nil, "", err.Mark()
	}
	fileBytes, filePath, err := generateFileBytesAndPath(excel, excelName)
	if err != nil {
		return nil, "", err.Mark()
	}

	return fileBytes, filePath, nil
}

func generateStreamExcelSheetList(sheetList []*ExcelSheetTab, opt *GenerateExcelOpt) (*excelize.File, *utilerror.UtilError) {
	// 生成Excel的文件
	excel := excelize.NewFile()
	for i := 0; i < len(sheetList); i++ {
		var err *utilerror.UtilError
		excel, err = generateStreamExcelSheet(i, sheetList[i], excel, opt)
		if err != nil {
			return nil, err.Mark()
		}
	}
	return excel, nil
}

func generateFileBytesAndPath(excel *excelize.File, excelName string) ([]byte, string, *utilerror.UtilError) {

	if !strings.HasSuffix(excelName, ".xlsx") {
		return nil, "", utilerror.NewError("The Excel name should end in .xlsx")
	}

	//store excel 存储Excel
	exportDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, "", utilerror.NewError(err.Error())
	}

	// 保存目录
	filePath := exportDir + excelName
	err = excel.SaveAs(filePath)
	if err != nil {
		return nil, "", utilerror.NewError(err.Error())
	}

	fileBytes, errRead := ioutil.ReadFile(filePath)
	if errRead != nil {
		return nil, "", utilerror.NewError("download excel fail")
	}
	return fileBytes, filePath, nil
}

func generateStreamExcelSheet(i int, sheet *ExcelSheetTab, excel *excelize.File, opt *GenerateExcelOpt) (*excelize.File, *utilerror.UtilError) {
	itemList := sheet.Data
	sheetName := sheet.SheetName
	if i == 0 {
		if sheetName != "Sheet1" {
			excel.SetSheetName("Sheet1", sheetName)
		}
	} else {
		excel.NewSheet(sheetName)
	}
	errCheck := CheckItemListType(itemList)
	if errCheck != nil {
		return nil, errCheck.Mark()
	}

	fieldValue := reflect.ValueOf(itemList)
	fieldType := reflect.TypeOf(itemList).Elem().Elem()
	sliceLength := fieldValue.Len()
	fieldNum := fieldType.NumField()

	excludeTitles := collection.NewStringSet(sheet.ExcludeTitles...)
	// 获取流式写入器
	streamWriter, err := excel.NewStreamWriter(sheetName)
	if err != nil {
		return nil, utilerror.NewError(err.Error())
	}

	titleList := make([]interface{}, 0)

	// 设置头部
	for i := 0; i < fieldNum; i++ {
		name := fieldType.Field(i).Tag.Get("title")
		if name == "" {
			continue
		}
		if excludeTitles.Contains(name) {
			continue
		}
		titleList = append(titleList, excelize.Cell{Value: name})
	}
	if err := streamWriter.SetRow("A1", titleList); err != nil {
		return nil, utilerror.NewError(err.Error())
	}
	//默认时间格式
	defaultStyle, _ := excel.NewStyle(&excelize.Style{NumFmt: 22, Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"}})

	//set rows
	for i := 0; i < sliceLength; i++ {
		// 得到某一个具体的结构体的
		structValue := fieldValue.Index(i).Elem()

		rowList := make([]interface{}, 0)

		//跳过的字段
		unUsedColumns := 0
		for j := 0; j < fieldNum; j++ {
			name := fieldType.Field(j).Tag.Get("title")
			if name == "" {
				unUsedColumns++
				continue
			}
			if excludeTitles.Contains(name) {
				unUsedColumns++
				continue
			}
			structValue := structValue.Field(j)

			// 直接写入对于列表类型中间分隔符为空格
			//rowList = append(rowList, structValue)
			structType := fieldType.Field(j)
			var cellStyle *int
			if structValue.Kind() == reflect.Struct {
				switch structValue.Interface().(type) {
				case time.Time:
					timeFormat := fieldType.Field(j).Tag.Get("format")
					if timeFormat != "" { //自定义时间格式
						if _, ok := timeStyleDict[timeFormat]; ok {
							cellStyle = convert.Int(timeStyleDict[timeFormat])
							break
						}
					}
					cellStyle = convert.Int(defaultStyle)
				default:
					break
				}
			}
			fieldIsEncrypt := fieldType.Field(j).Tag.Get("encrypt")
			isEncrypt := getIsEncrypt(fieldIsEncrypt, opt)
			elem, err := getCell(structValue, structType, cellStyle, isEncrypt)
			if err != nil {
				return nil, err.Mark()
			}
			rowList = append(rowList, elem)
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := streamWriter.SetRow(cell, rowList); err != nil {
			return nil, utilerror.NewError(err.Error())
		}
	}
	if fErr := streamWriter.Flush(); fErr != nil {
		return nil, utilerror.NewError(fErr.Error())
	}

	return excel, nil
}

func CheckItemListType(itemList interface{}) *utilerror.UtilError {
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

var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF",
	"AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ"}

// 自定义时间格式，存在于 styles.go 文件
var timeStyleDict = map[string]int{
	"mm-dd-yy":      14,
	"d-mmm-yy":      15,
	"d-mmm":         16,
	"mmm-yy":        17,
	"h:mm am/pm":    18,
	"h:mm:ss am/pm": 19,
	"h:mm":          20,
	"h:mm:ss":       21,
}

func getIsEncrypt(fieldIsEncrypt string, opt *GenerateExcelOpt) bool {
	return opt != nil && convert.BoolValue(opt.IsEncrypt) && fieldIsEncrypt == CommonEnumYES
}

func getCell(cellValue reflect.Value, cellType reflect.StructField, cellStyle *int, isEncrypt bool) (*excelize.Cell, *utilerror.UtilError) {
	result := &excelize.Cell{}
	var val interface{}
	var err *utilerror.UtilError
	if isEncrypt {
		val = EncryptChar
	} else {
		val, err = getCellValue(cellValue, cellType)
		if err != nil {
			return nil, err.Mark()
		}
	}
	result.Value = val
	if cellStyle != nil {
		result.StyleID = *cellStyle
	}
	return result, nil
}

func getCellValue(cellValue reflect.Value, cellType reflect.StructField) (interface{}, *utilerror.UtilError) {
	var result interface{}

	switch cellValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		result = cellValue.Int()
	case reflect.Int64:
		if cellValue.Int() > 19700101000 { //比较长的整数转字符串不然会变成科学计数法
			result = strconv.FormatInt(cellValue.Int(), 10)
		} else {
			result = cellValue.Int()
		}
	case reflect.Float64, reflect.Float32:
		result = cellValue.Float()
	case reflect.String:
		result = cellValue.String()
	case reflect.Slice:
		var buf bytes.Buffer
		elemKind := cellType.Type.Elem().Kind()
		switch elemKind {
		case reflect.Int64:
			for k := 0; k < cellValue.Len(); k++ {
				value := strconv.FormatInt(cellValue.Index(k).Int(), 10)
				buf.WriteString(value)
				if k != cellValue.Len()-1 {
					buf.WriteString(",")
				}
			}
		case reflect.Float64:
			for k := 0; k < cellValue.Len(); k++ {
				value := strconv.FormatFloat(cellValue.Index(k).Float(), 'f', -1, 64)
				buf.WriteString(value)
				if k != cellValue.Len()-1 {
					buf.WriteString(",")
				}
			}
		case reflect.String:
			for k := 0; k < cellValue.Len(); k++ {
				buf.WriteString(cellValue.Index(k).String())
				if k != cellValue.Len()-1 {
					buf.WriteString(",")
				}
			}
		default:
			return "", utilerror.NewError("type does not meet the requirements")
		}
		if buf.String() != "" {
			result = buf.String()
		}
	case reflect.Struct:
		switch cellValue.Interface().(type) {
		case time.Time:
			//这里应该是 time.Time,后续需要处理excel数据设置为时间格式 其他结构体需要自行处理或者在这补充
			orignalTime := cellValue.Interface().(time.Time)
			if orignalTime.Year() == 1970 || orignalTime.Year() == 1 { //1970的时间不处理
				timStr := "-"
				return timStr, nil
			}
			//excel处理只能传utc时间，所以将时间转为utc时间，时分秒没变，相当于如果是东八区时间，调小八小时后转为utc时间
			result = time.Date(orignalTime.Year(), orignalTime.Month(), orignalTime.Day(), orignalTime.Hour(), orignalTime.Minute(), orignalTime.Second(), orignalTime.Nanosecond(), time.FixedZone("CST", 0)).UTC()
		default:
			return "", utilerror.NewError("type does not meet the requirements")
		}

	}
	return result, nil
}

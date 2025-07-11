package expression

import (
	"example.com/m/util/convert"
	"example.com/m/util/utilerror"
	"fmt"
	"github.com/shopspring/decimal"

	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	MaxInt64 = 1<<63 - 1
	MinInt64 = -1 << 63
	MAXInt   = int(^uint(0) >> 1) //     最大值，根据二进制补码，第一位为0，其余为1
	MINInt   = ^MAXInt            // 最小值，第一位为1，其余为0，最大值取反即可
)

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfInt64(condition bool, trueVal, falseVal int64) int64 {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfString(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func Max(nums ...int64) int64 {
	var maxNum int64 = MinInt64
	for _, num := range nums {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}

func MaxInt(nums ...int) int {
	var maxNum int = MINInt
	for _, num := range nums {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}

func Min(nums ...int64) int64 {
	var minNum int64 = MaxInt64
	for _, num := range nums {
		if num < minNum {
			minNum = num
		}
	}
	return minNum
}

func MinInt(nums ...int) int {
	var minNum = MAXInt
	for _, num := range nums {
		if num < minNum {
			minNum = num
		}
	}
	return minNum
}

func Default(src, def interface{}) interface{} {
	// 确保src 和def是相同的指针类型
	if src == nil || reflect.ValueOf(src).IsNil() {
		return def
	}
	return src
}

func Abs(num int64) int64 {
	if num < 0 {
		return -num
	}
	return num
}

func IfNilUseDefaultInt64(pointer *int64, defaultVal int64) int64 {
	if pointer == nil {
		return defaultVal
	}
	return *pointer
}

// 判断nums是否每个元素都在assertFunc中返回true
func AllMatch(assertFunc func(int64) bool, nums ...int64) bool {
	for _, n := range nums {
		if !assertFunc(n) {
			return false
		}
	}
	return true
}

func IsTheTimePeriodCrossed(periodStart1, periodEnd1, periodStart2, periodEnd2 int64) bool {
	res := periodStart1 <= periodStart2 && periodStart2 < periodEnd1 ||
		periodStart1 < periodEnd2 && periodEnd2 <= periodEnd1 ||
		periodStart2 <= periodStart1 && periodStart1 < periodEnd2 ||
		periodStart2 < periodEnd1 && periodEnd1 <= periodEnd2
	//fmt.Printf("IsTheTimePeriodCrossed %v %v %v %v = %v\n", periodStart1, periodEnd1, periodStart2, periodEnd2, res)
	return res
}

func MinFloat64(a, b string) (string, *utilerror.UtilError) {
	floatA, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return "", utilerror.NewError(err.Error())
	}
	floatB, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return "", utilerror.NewError(err.Error())
	}
	if floatA < floatB {
		return a, nil
	}
	return b, nil
}

func Sum(list []int64) int64 {
	res := int64(0)
	for _, a := range list {
		res += a
	}
	return res
}

func Median(list ...int64) float64 {
	if len(list) == 0 {
		return 0
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	if len(list)%2 == 0 {
		return float64(list[len(list)/2-1]+list[len(list)/2]) / 2
	}
	return float64(list[len(list)/2])
}

func Keep1DecimalPlaces(val float64) string {
	str := fmt.Sprintf("%.1f", val)
	//res, _ := strconv.ParseFloat(str, 64)
	return str
}

func Keep1DecimalPlacesString(val string) string {
	_, err := convert.StringToFloat(val)
	if err != nil {
		return val
	}
	parts := strings.Split(val, ".")
	if len(parts) < 2 {
		return val + ".0"
	}
	decimalPart := parts[1]
	if len(decimalPart) > 1 {
		decimalPart = decimalPart[:1]
	}
	return parts[0] + "." + decimalPart
}

func ReverseMatrixRowsAndCols(matrix [][]int64) [][]int64 {
	if len(matrix) == 0 {
		return matrix
	}
	res := make([][]int64, len(matrix[0]))
	for i := range res {
		res[i] = make([]int64, len(matrix))
	}
	for i := range matrix {
		for j := range matrix[0] {
			res[j][i] = matrix[i][j]
		}
	}
	return res
}

func SumDecimal(list ...decimal.Decimal) decimal.Decimal {
	res := decimal.NewFromInt(0)
	for _, a := range list {
		res = res.Add(a)
	}
	return res
}

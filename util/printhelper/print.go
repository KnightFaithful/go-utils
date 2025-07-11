package printhelper

import (
	"example.com/m/util/stringutil"
	"fmt"
	"github.com/shopspring/decimal"

	"gorm.io/gorm/logger"
	"runtime"
	"strings"
)

func PrintIntList(arr []int) {
	// 计算数组长度和最大值
	length := len(arr)
	maxValue := 0
	for _, value := range arr {
		if value > maxValue {
			maxValue = value
		}
	}

	// 计算每个元素的最大宽度
	valueWidth := len(fmt.Sprint(maxValue * 10))

	printTop := func() {
		fmt.Println()
		for i := 0; i < length; i++ {
			index := fmt.Sprintf("%v", strings.Repeat("-", valueWidth+1))
			fmt.Print(index)
		}
		fmt.Println()
	}

	printTop()

	// 输出索引
	for i := 0; i < length; i++ {
		index := fmt.Sprintf("|%*d", valueWidth, i)
		fmt.Print(index)
	}

	printTop()

	// 输出值
	for _, value := range arr {
		valueStr := fmt.Sprintf("|%*d", valueWidth, value)
		fmt.Print(valueStr)
	}

	printTop()

}

func Printf(format string, a ...interface{}) {
	fmt.Println()
	fmt.Println(logger.Blue, fmt.Sprintf(format, a...), logger.Reset)
}

func Println(a ...interface{}) {
	fmt.Println()

	pc, _, line, _ := runtime.Caller(1)
	// 使用反射获取 PC 对应的函数
	fn := runtime.FuncForPC(pc)

	// 使用反射获取函数名
	methodName := fn.Name()

	a = append([]interface{}{methodName, fmt.Sprintf("%v行: ", line)}, a...)
	fmt.Println(logger.Green, fmt.Sprintln(a...), logger.Reset)
}

func PrintErrorF(format string, a ...interface{}) {
	fmt.Println()
	fmt.Println(logger.Red, fmt.Sprintf(format, a...), logger.Reset)
}

func PrintlnDecimal(list ...decimal.Decimal) {
	fmt.Println()
	var strList []string
	for _, item := range list {
		strList = append(strList, item.String())
	}
	str := stringutil.JoinIgnoreEmpty(strList, ",")
	fmt.Println(logger.YellowBold, fmt.Sprintf(str), logger.Reset)
}

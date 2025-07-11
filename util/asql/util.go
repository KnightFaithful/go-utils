package asql

import (
	"math"
	"strings"
)

func getFieldName(fieldTag string) string {
	fieldTagArr := strings.Split(fieldTag, ":")
	if len(fieldTagArr) == 0 {
		return ""
	}

	fieldName := fieldTagArr[len(fieldTagArr)-1]

	return fieldName
}

func getSQLQuantity(length, size int) int {
	SQLQuantity := int(math.Ceil(float64(length) / float64(size)))
	return SQLQuantity
}

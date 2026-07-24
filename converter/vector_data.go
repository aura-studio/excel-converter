package converter

import "strings"

func isVectorEgUseExcel(excelName string) bool {
	excelName = strings.TrimSuffix(excelName, ".xlsx")
	return strings.HasSuffix(excelName, "VectorEgUse")
}

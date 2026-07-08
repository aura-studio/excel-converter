package converter

import "strings"

func skipVectorDataNode(node Node) bool {
	if node == nil {
		return false
	}
	return skipVectorDataExcelPath(node.ExcelPathName())
}

func skipVectorDataExcelPath(excelPath string) bool {
	excelPath = strings.ReplaceAll(excelPath, "\\", "/")
	parts := strings.Split(excelPath, "/")
	if len(parts) == 0 {
		return false
	}
	excelName := parts[len(parts)-1]
	return strings.Contains(excelName, "VectorEg")
}

func skipVectorDataLinkPath(linkPath LinkPath) bool {
	return skipVectorDataExcelPath(linkPath.ExcelName)
}

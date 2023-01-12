package converter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExcelBase struct {
	fieldType FieldType
	relPath   string
	file      *excelize.File
	sheetMap  map[SheetType]map[string]Sheet
	nodes     []Node
}

func NewExcelBase(path Path, relPath string, fieldType FieldType) *ExcelBase {
	return &ExcelBase{
		fieldType: fieldType,
		relPath:   relPath,
		sheetMap:  map[SheetType]map[string]Sheet{},
	}
}

func (e *ExcelBase) String() string {
	return e.relPath
}

func (e *ExcelBase) Read() {

}

func (e *ExcelBase) Preprocess() {
	for _, sheet := range e.sheetMap[SheetTypeRegular] {
		sheet.FormatHeader(e.fieldType)
	}
	for _, sheet := range e.sheetMap[SheetTypeInferior] {
		sheet.FormatHeader(e.fieldType)
	}
	for _, sheet := range e.sheetMap[SheetTypeSettings] {
		sheet.FormatHeader(e.fieldType)
	}
}

func (e *ExcelBase) SheetMap() map[SheetType]map[string]Sheet {
	return e.sheetMap
}

func (e *ExcelBase) ReadFile() map[string][][]string {
	var data = make(map[string][][]string)
	var err error
	e.file, err = excelize.OpenFile(filepath.Join(path.ImportAbsPath(), e.relPath))
	if err != nil {
		Exit("[%v] Read %v error, %v", e, e, err)
	}
	sheetMap := e.file.GetSheetMap()
	for _, sheetName := range sheetMap {
		rows, err := e.file.GetRows(sheetName)
		if err != nil {
			Exit("[%v] Read %v error, %v", e, sheetName, err)
		}
		data[sheetName] = rows
	}
	return data
}

func (e *ExcelBase) PackageName() string {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	if len(strs) < 3 {
		Exit("[%v] Error rel path", e)
	}
	return strs[0]
}

func (e *ExcelBase) DomainName() string {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	length := len(strs)
	switch {
	case length < 3:
		Exit("[%v] Error rel path", e)
	case length == 3: // excel当作Domain
		return e.FixedName()
	case length > 3:
		return format.ToUpper(strs[2])
	}
	return FlagDefault
}

func (e *ExcelBase) IndirectName() string {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	length := len(strs)
	switch {
	case length < 3:
		Exit("[%v] Error rel path", e)
	case length == 3: // excel当作Domain
		return FlagDefault
	case length > 3:
		return e.FixedName()
	}
	return FlagDefault
}

func (e *ExcelBase) Name() string {
	return filepath.Base(e.relPath)
}

func (e *ExcelBase) FixedName() string {
	return format.ToUpper(e.Name())
}

func (e *ExcelBase) Type() string {
	Exit("[Main] Invalid call ExcelBase.Type")
	return ""
}

func (e *ExcelBase) GetSheet(sheetName string) Sheet {
	for _, sheetMap := range e.sheetMap {
		if sheet, ok := sheetMap[sheetName]; ok {
			return sheet
		}
	}
	Exit("[%v] Sheet %s not found", e, sheetName)
	return nil
}

func (e *ExcelBase) GetHeaderSize(sheetName string) int {
	sheet := e.GetSheet(sheetName)
	return sheet.HeaderSize()
}

func (e *ExcelBase) GetHeaderField(sheetName string, key interface{}) HeaderField {
	sheet := e.GetSheet(sheetName)
	return sheet.GetHeaderField(key)
}

func (e *ExcelBase) Build() {

}

func (e *ExcelBase) Nodes() []Node {
	return e.nodes
}

func (e *ExcelBase) ForServer() bool {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	for _, str := range strs {
		if format.SpecificClient(str) {
			return false
		}
	}
	return true
}

func (e *ExcelBase) ForClient() bool {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	for _, str := range strs {
		if format.SpecificServer(str) {
			return false
		}
	}
	return true
}

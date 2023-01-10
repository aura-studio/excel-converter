package converter

import (
	"path/filepath"
	"sort"
	"strings"
)

type ExcelRegular struct {
	*ExcelBase
}

func NewExcelRegular(path Path, relPath string, fieldType FieldType) *ExcelRegular {
	return &ExcelRegular{
		ExcelBase: NewExcelBase(path, relPath, fieldType),
	}
}

func (e *ExcelRegular) SheetType(sheetName string) SheetType {
	switch {
	case strings.Contains(sheetName, FlagComment):
		return SheetTypeComment
	case strings.HasPrefix(sheetName, FlagUnderScore):
		return SheetTypeInferior
	default:
		return SheetTypeRegular
	}
}

func (e *ExcelRegular) Read() {
	data := e.ReadFile()
	for sheetName, rows := range data {
		sheetType := e.SheetType(sheetName)
		if e.sheetMap[sheetType] == nil {
			e.sheetMap[sheetType] = make(map[string]Sheet, 0)
		}
		sheet := sheetCreators[sheetType](e, sheetName, rows)
		sheet.Read()
		e.sheetMap[sheetType][sheetName] = sheet
	}
}

func (e *ExcelRegular) Build() {
	regularSheets := e.sheetMap[SheetTypeRegular]
	for _, sheet := range regularSheets {
		structure := e.parseStructrue(sheet)
		sheet.ParseContent(structure)
		node := NewNode(NewNodeContext(), sheet, structure)
		e.nodes = append(e.nodes, node)
		for _, key := range e.parseKeys(sheet) {
			ctx := NewNodeContext()
			ctx.key = key
			node := NewNode(ctx, sheet, StructureTypeMap)
			e.nodes = append(e.nodes, node)
		}
	}
	sort.Slice(e.nodes, func(i, j int) bool {
		return e.nodes[i].RootName() < e.nodes[j].RootName()
	})
}

func (*ExcelRegular) Type() ExcelType {
	return ExcelTypeRegular
}

func (e *ExcelRegular) parseStructrue(sheet Sheet) StructureType {
	ext := filepath.Ext(sheet.Name())
	if ext == "" {
		return StructureTypeStructs
	}
	switch structure := StructureType(ext[1:]); structure {
	case StructureTypeStructs, StructureTypeRows, StructureTypeCols,
		StructureTypeMap, StructureTypeRowMap, StructureTypeColMap:
		return structure
	default:
		Exit("[%v] Unsupported structure %s", e, structure)
	}
	return ""
}

func (e *ExcelRegular) parseKeys(sheet Sheet) []string {
	var keys []string
	for index := 0; index < sheet.HeaderSize(); index++ {
		field := NewField(sheet, sheet.GetHeaderField(index))
		if field.Key {
			keys = append(keys, field.Name)
		}
	}
	return keys
}

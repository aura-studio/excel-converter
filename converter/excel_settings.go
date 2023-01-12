package converter

type ExcelSettings struct {
	*ExcelBase
}

func NewExcelSettings(path Path, relPath string, fieldType FieldType) *ExcelSettings {
	return &ExcelSettings{
		ExcelBase: NewExcelBase(path, relPath, fieldType),
	}
}

func (e *ExcelSettings) Read() {
	data := e.ReadFile()
	for sheetName, rows := range data {
		sheetType := SheetTypeSettings
		if e.sheetMap[sheetType] == nil {
			e.sheetMap[sheetType] = make(map[string]Sheet, 0)
		}
		sheet := sheetCreators[sheetType](e, sheetName, rows)
		sheet.Read()
		e.sheetMap[sheetType][sheetName] = sheet
	}
}

func (e *ExcelSettings) Preprocess() {
	e.ExcelBase.Preprocess()
	settingsSheets := e.sheetMap[SheetTypeSettings]
	for _, sheet := range settingsSheets {
		structure := e.parseStructure(sheet)
		sheet.ParseContent(structure)
	}
}

func (e *ExcelSettings) IndirectName() string {
	return FlagSettings
}

func (e *ExcelSettings) PackageName() string {
	return FlagDefault
}

func (e *ExcelSettings) DomainName() string {
	return FlagDefault
}

func (*ExcelSettings) Type() ExcelType {
	return ExcelTypeSettings
}

func (e *ExcelSettings) parseStructure(sheet Sheet) StructureType {
	switch sheet.Name() {
	case FlagVarian:
		return StructureTypeStructs
	case FlagLink:
		return StructureTypeRows
	case FlagCategory:
		return StructureTypeCols
	default:
		Exit("[%v] Unsupported structure for settings name %s", e, sheet.Name())
	}
	return ""
}

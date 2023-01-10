package converter

type ExcelType string

const (
	ExcelTypeComment  ExcelType = "comment"
	ExcelTypeSettings ExcelType = "settings"
	ExcelTypeTemplate ExcelType = "template"
	ExcelTypeRegular  ExcelType = "regular"
)

type Excel interface {
	String() string
	Read()
	Preprocess()
	SheetMap() map[SheetType]map[string]Sheet
	Name() string
	FixedName() string
	PackageName() string
	DomainName() string
	IndirectName() string
	Type() ExcelType
	GetSheet(string) Sheet
	GetHeaderSize(string) int
	GetHeaderField(string, interface{}) HeaderField
	Build()
	Nodes() []Node
	ForServer() bool
	ForClient() bool
}

var excelCreators map[ExcelType]func(Path, string, FieldType) Excel

func init() {
	excelCreators = map[ExcelType]func(Path, string, FieldType) Excel{
		ExcelTypeComment: func(path Path, relPath string, fieldType FieldType) Excel {
			return NewExcelComment(path, relPath, fieldType)
		},
		ExcelTypeSettings: func(path Path, relPath string, fieldType FieldType) Excel {
			return NewExcelSettings(path, relPath, fieldType)
		},
		ExcelTypeTemplate: func(path Path, relPath string, fieldType FieldType) Excel {
			return NewExcelTemplate(path, relPath, fieldType)
		},
		ExcelTypeRegular: func(path Path, relPath string, fieldType FieldType) Excel {
			return NewExcelRegular(path, relPath, fieldType)
		},
	}
}

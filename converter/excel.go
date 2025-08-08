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
	Category() string
	Type() ExcelType
	GetSheet(string) Sheet
	GetHeaderSize(string) int
	GetHeaderField(string, any) HeaderField
	Build()
	Nodes() []Node
	ForServer() bool
	ForClient() bool
}

var excelCreators map[ExcelType]func(Path, string) Excel

func init() {
	excelCreators = map[ExcelType]func(Path, string) Excel{
		ExcelTypeComment: func(path Path, relPath string) Excel {
			return NewExcelComment(path, relPath)
		},
		ExcelTypeSettings: func(path Path, relPath string) Excel {
			return NewExcelSettings(path, relPath)
		},
		ExcelTypeTemplate: func(path Path, relPath string) Excel {
			return NewExcelTemplate(path, relPath)
		},
		ExcelTypeRegular: func(path Path, relPath string) Excel {
			return NewExcelRegular(path, relPath)
		},
	}
}

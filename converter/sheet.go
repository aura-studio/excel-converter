package converter

type SheetType string

const (
	SheetTypeComment  SheetType = "comment"
	SheetTypeInferior SheetType = "inferior"
	SheetTypeRegular  SheetType = "regular"
	SheetTypeSettings SheetType = "settings"
)

type Sheet interface {
	String() string
	Excel() Excel
	Read()
	Type() SheetType
	Name() string
	FixedName() string
	IndirectName() string
	HeaderSize() int
	VerticleSize() int
	GetHeaderField(interface{}) HeaderField
	GetHorizon(int) []string
	GetVerticle(int) []string
	GetCell(HeaderField, int) string
	GetIndex(interface{}) int
	ParseContent(StructureType)
	FormatHeader(FieldType)
	FormatContent()
	ForServer() bool
	ForClient() bool
}

var sheetCreators map[SheetType]func(excel Excel, name string, rows [][]string) Sheet

func init() {
	sheetCreators = map[SheetType]func(excel Excel, name string, rows [][]string) Sheet{
		SheetTypeComment: func(excel Excel, name string, data [][]string) Sheet {
			return NewSheetComment(excel, name, data)
		},
		SheetTypeInferior: func(excel Excel, name string, data [][]string) Sheet {
			return NewSheetInferior(excel, name, data)
		},
		SheetTypeRegular: func(excel Excel, name string, data [][]string) Sheet {
			return NewSheetRegular(excel, name, data)
		},
		SheetTypeSettings: func(excel Excel, name string, data [][]string) Sheet {
			return NewSheetSettings(excel, name, data)
		},
	}
}

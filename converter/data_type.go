package converter

type DataExportType int

var dataExportType DataExportType

const (
	DataExportTypeLiteral DataExportType = iota
	DataExportTypeJSON
)

func SetDataExportType(t DataExportType) {
	dataExportType = t
}

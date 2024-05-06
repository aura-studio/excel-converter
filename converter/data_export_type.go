package converter

type DataType int

var dataType DataType

const (
	DataTypeLiteral DataType = iota
	DataTypeJSON
)

func SetDataType(t DataType) {
	dataType = t
}

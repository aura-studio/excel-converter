package converter

type SourceType string

const (
	SourceTypeNil            SourceType = "nil"
	SourceTypeTable          SourceType = "table"            // 二维表单
	SourceTypeRecordStruct   SourceType = "record_struct"    // 一维记录
	SourceTypeRecordSlice    SourceType = "record_slice"     // 一维记录
	SourceTypeRecordSliceMap SourceType = "record_slice_map" // 一维记录
	SourceTypeTextArray      SourceType = "text_array"       // 单字符串
	SourceTypeTextTable      SourceType = "text_table"       // 单字符串
	SourceTypeTextDict       SourceType = "text_dict"        // 单字符串
	SourceTypeTextList       SourceType = "text_list"        // 单单字符串
	SourceTypeTextElement    SourceType = "text_element"     // 单字符串
	SourceTypeTextNil        SourceType = "text_nil"         // 空引用
)

type Source interface {
	Type() SourceType
	Sources() []Source
	Content() string
	AppendSource(Source)
	Sheet() Sheet
	Excel() Excel
	Empty() bool
}

type SourceContext struct {
	key string
}

func NewSourceContext() *SourceContext {
	return &SourceContext{}
}

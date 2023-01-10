package converter

type SectionType string

const (
	SectionTypeComment SectionType = "comment"
	SectionTypeHeader  SectionType = "header"
	SectionTypeRows    SectionType = "rows"
	SectionTypeColumns SectionType = "columns"
	SectionTypeBlock   SectionType = "block"
	SectionTypeContent SectionType = "content"
)

type Section interface {
	Type() SectionType
	Append(int, []string)
	Size() int
	GetHorizon(int, []int) []string
	GetVertical(int, []int) []string
	Format(int)
	GetIndex(interface{}) int
}

var sectionCreators map[SectionType]func(Sheet) Section

func init() {
	sectionCreators = map[SectionType]func(Sheet) Section{
		SectionTypeComment: func(sheet Sheet) Section {
			return NewSectionComment(sheet)
		},
		SectionTypeHeader: func(sheet Sheet) Section {
			return NewSectionHeader(sheet)
		},
		SectionTypeRows: func(sheet Sheet) Section {
			return NewSectionRows(sheet)
		},
		SectionTypeColumns: func(sheet Sheet) Section {
			return NewSectionColumns(sheet)
		},
		SectionTypeBlock: func(sheet Sheet) Section {
			return NewSectionBlock(sheet)
		},
		SectionTypeContent: func(sheet Sheet) Section {
			return NewSectionContent(sheet)
		},
	}
}

func NewSection(sheet Sheet, sectionType SectionType) Section {
	sectionCreator, ok := sectionCreators[sectionType]
	if !ok {
		Exit("[%v] Unsupported section type %v", sheet, sectionType)
	}
	return sectionCreator(sheet)
}

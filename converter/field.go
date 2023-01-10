package converter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

var (
	arrayRegExp    = regexp.MustCompile(`^array:(.+)$`)
	tableRegExp    = regexp.MustCompile(`^table:\[([^\]]*)\](.*)$`)
	keyRegExp      = regexp.MustCompile(`^(.+):key$`)
	inferiorRegExp = regexp.MustCompile(`^([^:].*):(_.*)$`)
	listRegExp     = regexp.MustCompile(`^list\[(.+):(\d+)\]$`)
	dictRegExp     = regexp.MustCompile(`^dict\[(.+)]$`)
)

type Field struct {
	HeaderField
	Sheet     Sheet
	Name      string
	Structure StructureType
	Args      []string
	Size      int
	Key       bool
}

func NewField(sheet Sheet, v interface{}) *Field {
	f := &Field{}
	f.Sheet = sheet
	switch v := v.(type) {
	case StructureType:
		f.Structure = v
	case HeaderField:
		f.HeaderField = v
		f.Name = v.Name
	case string:
		f.HeaderField = HeaderField{
			Index: -1,
			Name:  "",
			Type:  v,
		}
	}
	f.parseType()
	f.parseSize()
	return f
}

func (f *Field) String() string {
	if f.HeaderField.Type == "" {
		return string(f.Structure)
	}
	return fmt.Sprintf("%v(%v)", f.Structure, f.Name)
}

func (f *Field) parseType() {
	if f.HeaderField.Type == "" {
		return
	}
	if structure, ok := f.parseKey(); ok {
		if structure != StructureTypeInt && structure != StructureTypeString && structure != StructureTypeBool && structure != StructureTypeFloat {
			Exit("[%v] Invalid key type %s", f.Sheet, f.HeaderField.Type)
		}
		f.Structure = structure
		f.Key = true
		return
	}
	if structure, subType, ok := f.parseArray(); ok {
		f.Structure = structure
		f.Args = append(f.Args, subType)
		return
	}
	if structure, keyType, valType, ok := f.parseTable(); ok {
		f.Structure = structure
		f.Args = append(f.Args, keyType, valType)
		return
	}
	if structure, itemType, itemCount, ok := f.parseList(); ok {
		f.Structure = structure
		f.Args = append(f.Args, itemType, itemCount)
		return
	}
	if structure, args, ok := f.parseDict(); ok {
		f.Structure = structure
		f.Args = append(f.Args, args...)
		return
	}
	if structure, InferiorSheetName, ok := f.parseInferior(); ok {
		if structure != StructureTypeStruct && structure != StructureTypeCol && structure != StructureTypeRow && structure != StructureTypeRowMapVal && structure != StructureTypeColMapVal && structure != StructureTypeMapVal &&
			structure != StructureTypeStructs && structure != StructureTypeCols && structure != StructureTypeRows && structure != StructureTypeRowMap && structure != StructureTypeColMap && structure != StructureTypeMap {
			Exit("[%v] Invalid inferior type %s", f.Sheet, f.HeaderField.Type)
		}
		f.Structure = structure
		f.Args = append(f.Args, InferiorSheetName)
		return
	}
	f.Structure = StructureType(format.ToLowerRaw(f.HeaderField.Type))
}

func (f *Field) parseKey() (StructureType, bool) {
	if submatches := keyRegExp.FindAllStringSubmatch(f.HeaderField.Type, -1); len(submatches) == 1 {
		return StructureType(format.ToLowerRaw(submatches[0][1])), true
	}
	return "", false
}

func (f *Field) parseArray() (StructureType, string, bool) {
	if submatches := arrayRegExp.FindAllStringSubmatch(f.HeaderField.Type, -1); len(submatches) == 1 {
		return StructureTypeArray, submatches[0][1], true
	}
	return "", "", false
}

func (f *Field) parseTable() (StructureType, string, string, bool) {
	if submatches := tableRegExp.FindAllStringSubmatch(f.HeaderField.Type, -1); len(submatches) == 1 {
		return StructureTypeTable, submatches[0][1], submatches[0][2], true
	}
	return "", "", "", false
}

func (f *Field) parseInferior() (StructureType, string, bool) {
	if submatches := inferiorRegExp.FindAllStringSubmatch(f.HeaderField.Type, -1); len(submatches) == 1 {
		return StructureType(format.ToLowerRaw(submatches[0][1])), submatches[0][2], true
	}
	return "", "", false
}

func (f *Field) parseList() (StructureType, string, string, bool) {
	if submatches := listRegExp.FindAllStringSubmatch(f.HeaderField.Type, -1); len(submatches) == 1 {
		return StructureTypeList, submatches[0][1], submatches[0][2], true
	}
	return "", "", "", false
}

func (f *Field) parseDict() (StructureType, []string, bool) {
	var splitDict = func(s string) []string {
		depth := 1
		commaIndex := -1
		strs := []string{}
		for index := 0; index < len(s); index++ {
			if s[index] == '[' {
				depth += 1
			} else if s[index] == ']' {
				depth -= 1
			}
			if depth > 1 {
				continue
			}
			if s[index] == ',' {
				strs = append(strs, strings.Trim(s[commaIndex+1:index], " "))
				commaIndex = index
			}
		}
		strs = append(strs, strings.Trim(s[commaIndex+1:], " "))
		return strs
	}
	if submatches := dictRegExp.FindAllStringSubmatch(f.HeaderField.Type, -1); len(submatches) == 1 {
		return StructureTypeDict, splitDict(submatches[0][1]), true
	}

	return "", nil, false
}

func (f *Field) parseSize() {
	f.Size = f.recursiveSize()
}

func (f *Field) recursiveSize() int {
	switch f.Structure {
	case StructureTypeList:
		var size int
		if f.Index < 0 {
			size = 0
		} else {
			size = 1
		}
		arg := f.Args[0]
		field := NewField(f.Sheet, arg)
		return size + field.recursiveSize()*cast.ToInt(f.Args[1])
	case StructureTypeDict:
		var size int
		if f.Index < 0 {
			size = 0
		} else {
			size = 1
		}
		for _, arg := range f.Args {
			if sz, err := cast.ToIntE(arg); err == nil {
				size += sz
			} else {
				field := NewField(f.Sheet, arg)
				size += field.recursiveSize()
			}
		}
		return size
	default:
		return 1
	}
}

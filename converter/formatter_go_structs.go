package converter

import "fmt"

type FormatterGoStructs struct {
	*FormatterBase
	used       bool
	identifier *Identifier
}

func NewFormatterGoStructs(identifier *Identifier) *FormatterGoStructs {
	f := &FormatterGoStructs{
		FormatterBase: NewFormatterBase(),
		identifier:    identifier,
	}
	f.WriteString(`// Package structs <important: auto generate by excel-to-go converter, do not modify>
package structs

import (
	"fmt"
	"math/big"
	"time"
)
`)
	f.WriteString(fmt.Sprintf(`
var timeLocation = time.FixedZone("SYS", %s)
`, FlagTimeZone))

	f.WriteString(`
func NewTime(year, month, day, hour, minute, second int) *time.Time {
	tm := time.Date(year, time.Month(month), day, hour, minute, second, 0, timeLocation)
	return &tm
}

func NewBigInt(s string) *big.Int {
	i := new(big.Int)
	_, ok := i.SetString(s, 10)
	if !ok {
		panic(fmt.Errorf("big int:%s error in excel", s))
	}
	return i
}

func NewBigFloat(s string) *big.Float {
	f := new(big.Float)
	_, ok := f.SetString(s)
	if !ok {
		panic(fmt.Errorf("big float:%s error in excel", s))
	}
	return f
}

func NewBigRat(s string) *big.Rat {
	r := new(big.Rat)
	_, ok := r.SetString(s)
	if !ok {
		panic(fmt.Errorf("big rat:%s error in excel", s))
	}
	return r
}

type (
`)
	return f
}

func (f *FormatterGoStructs) Close() string {
	if !f.used {
		return ""
	}
	f.WriteString(")")
	return f.String()
}

func (f *FormatterGoStructs) FormatStruct() {
	f.used = true
	for _, node := range f.identifier.OriginNodes {
		f.WriteString("\t// ")
		f.WriteString(f.identifier.NodeStructMap[node.ID()])
		f.WriteString(" comment\n")
		f.WriteString("\t")
		f.WriteString(f.identifier.NodeStructMap[node.ID()])
		f.WriteString(" struct {\n")
		for _, node := range node.Nodes() {
			f.WriteString("\t\t")
			f.WriteString(node.FieldName())
			f.WriteString(" ")
			f.WriteString(f.identifier.NodeTypeMap[node.ID()])
			f.WriteString("\n")
		}
		f.WriteString("\t}\n")
	}
}

func (f *FormatterGoStructs) FormatStructEqual() {
	f.used = true
	for _, structNames := range f.identifier.StructEquals {
		f.WriteString("\t// ")
		f.WriteString(structNames[0])
		f.WriteString(" comment\n")
		f.WriteString("\t")
		f.WriteString(structNames[0])
		f.WriteString(" = ")
		f.WriteString(structNames[1])
		f.WriteString("\n")
		f.WriteString("\n")
	}
}

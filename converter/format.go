package converter

import (
	"bytes"
	"fmt"
	"math/big"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
)

const (
	dateTimeFormat = "2006-1-2 15:04:05"
	dateFormat     = "2006-1-2"
)

var (
	templateRegExp   = regexp.MustCompile(`^\d+_(.*)$`)
	contentKeyRegExp = regexp.MustCompile(`^\[(.*)\]$`)
)

type Format struct{}

var format Format

func (f Format) ToUpperRaw(s string) string {
	return strings.ToUpper(s)
}

func (f Format) ToLowerRaw(s string) string {
	return strings.ToLower(s)
}

func (f Format) ToUpper(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSuffix(s, filepath.Ext(s))
	s = strings.TrimPrefix(s, FlagClient)
	s = strings.TrimPrefix(s, FlagServer)
	if submatches := templateRegExp.FindAllStringSubmatch(s, -1); len(submatches) > 0 {
		s = submatches[0][1]
	}

	var buf bytes.Buffer
	buf.WriteString(strings.ToUpper(s[:1]))
	isAfterUnderLine := false
	for i := 1; i < len(s); i++ {
		if s[i] == '_' {
			isAfterUnderLine = true
		} else {
			if isAfterUnderLine {
				buf.WriteString(strings.ToUpper(s[i : i+1]))
			} else {
				buf.WriteByte(s[i])
			}
			isAfterUnderLine = false
		}
	}

	s = buf.String()

	if strings.HasSuffix(s, "Id") {
		return fmt.Sprintf("%sID", strings.TrimSuffix(s, "Id"))
	} else if strings.HasSuffix(s, FlagDefault) {
		return strings.TrimSuffix(s, FlagDefault)
	}

	return s
}

func (f Format) ToLower(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSuffix(s, filepath.Ext(s))
	if strings.HasPrefix(s, FlagClient) {
		return strings.TrimPrefix(s, FlagClient)
	} else if strings.HasPrefix(s, FlagServer) {
		return strings.TrimPrefix(s, FlagServer)
	}

	if submatches := templateRegExp.FindAllStringSubmatch(s, -1); len(submatches) > 0 {
		s = submatches[0][1]
	}

	allUpper := true
	for _, b := range s {
		if b < 65 || b > 90 {
			allUpper = false
			break
		}
	}
	runes := make([]rune, 0)
	if allUpper {
		for _, b := range s {
			runes = append(runes, b+32)
		}
	} else {
		for ix, b := range s {
			if b >= 65 && b <= 90 {
				if ix != 0 {
					runes = append(runes, 95)
				}
				runes = append(runes, b+32)
			} else {
				runes = append(runes, b)
			}
		}
	}

	s = string(runes)

	if strings.HasSuffix(s, "_i_d") {
		return fmt.Sprintf("%s_id", strings.TrimSuffix(s, "_i_d"))
	} else if strings.HasSuffix(s, "_default") {
		return strings.TrimSuffix(s, "_default")
	}

	return strings.ReplaceAll(s, "__", "_")
}

func (f Format) Indent(count int) string {
	var buf bytes.Buffer
	for ix := 0; ix < count; ix++ {
		buf.WriteString("\t")
	}
	return buf.String()
}

func (f Format) ToExcelCol(index int) string {
	rns := make([]rune, 0)
	for index != 0 {
		tmp := index % 26
		rns = append(rns, rune(tmp+64))
		index /= 26
	}
	// reverse
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func (f Format) ToExcelIndex(col string) int {
	col = strings.ToUpper(col)
	index := 0
	for ix := 0; ix < len(col); ix++ {
		index = index*26 + int(rune(col[ix])-64)
	}
	return index
}

func (f Format) ParseRange(content string) []interface{} {
	if content == "" {
		return nil
	}
	contents := strings.Split(content, FlagComma)
	var keys []interface{}
	for _, content := range contents {
		content = strings.Trim(content, " ")
		if content == "" {
			continue
		}
		if strings.Contains(content, FlagDash) {
			contents = strings.Split(content, FlagDash)
			for index := 0; index < len(contents); index++ {
				contents[index] = strings.Trim(contents[index], " ")
			}
			var keyFrom int
			if n, err := cast.ToIntE(contents[0]); err == nil {
				keyFrom = n
			} else {
				keyFrom = f.ToExcelIndex(contents[0])
			}

			var keyTo int
			if n, err := cast.ToIntE(contents[1]); err == nil {
				keyTo = n
			} else {
				keyTo = f.ToExcelIndex(contents[1])
			}

			for key := keyFrom; key <= keyTo; key++ {
				keys = append(keys, key)
			}
		} else {
			keys = append(keys, content)
		}
	}
	return keys
}

func (f Format) SpecificServer(s string) bool {
	return strings.HasPrefix(s, FlagServer)
}

func (f Format) SpecificClient(s string) bool {
	return strings.HasPrefix(s, FlagClient)
}

func (f Format) ToGoPackageCase(s string) string {
	return strings.ReplaceAll(format.ToLowerRaw(s), FlagUnderScore, "")
}

func (f Format) ToLuaPackageCase(s string) string {
	return format.ToLowerRaw(s)
}

func (f Format) ParseTime(s string) (time.Time, error) {
	var format string
	if len(s) > 10 {
		format = dateTimeFormat
	} else {
		format = dateFormat
	}
	return time.Parse(format, s)
}

func (f Format) BigIntToLua(s string) (string, bool) {
	i := new(big.Int)
	_, ok := i.SetString(s, 10)
	if !ok {
		return "", false
	}
	n, _ := new(big.Float).SetInt(i).Float64()
	return strconv.FormatFloat(n, 'E', -1, 64), true
}

func (f Format) BigFloatToLua(s string) (string, bool) {
	i := new(big.Float)
	_, ok := i.SetString(s)
	if !ok {
		return "", false
	}
	n, _ := i.Float64()
	return strconv.FormatFloat(n, 'E', -1, 64), true
}

func (f Format) ParseContentKey(s string) (string, bool) {
	if submatches := contentKeyRegExp.FindAllStringSubmatch(s, -1); len(submatches) == 1 {
		return submatches[0][1], true
	}
	return "", false
}

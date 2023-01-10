package converter

import (
	"strings"

	"github.com/spf13/cast"
)

type SectionBase struct {
	sheet    Sheet
	data     [][]string
	indexMap map[interface{}]int // key: int or string, val: index
}

func NewSectionBase(sheet Sheet) *SectionBase {
	return &SectionBase{
		sheet:    sheet,
		indexMap: make(map[interface{}]int),
	}
}

func (s *SectionBase) GetIndex(key interface{}) int {
	var indexKey = key
	if n, err := cast.ToIntE(key); err == nil {
		indexKey = n
	}
	index, ok := s.indexMap[indexKey]
	if !ok {
		Exit("[%v] Invalid index key %v", s.sheet, key)
	}
	return index
}

func (s *SectionBase) Append(index int, data []string) {
	s.data = append(s.data, data)
	keys := []interface{}{
		index + 1,
		cast.ToString(index + 1),
		format.ToExcelCol(index + 1),
	}
	val := len(s.data) - 1
	for _, key := range keys {
		s.indexMap[key] = val
	}
}

func (s *SectionBase) Size() int {
	return len(s.data)
}

func (s *SectionBase) GetHorizon(index int, fieldIndexes []int) []string {
	if index >= len(s.data) {
		Exit("[%v] Row index %d is out of length %d", s.sheet, index, len(s.data))
	}
	return s.data[index]
}

func (s *SectionBase) GetVertical(index int, fieldIndexes []int) []string {
	index = fieldIndexes[index]
	var vertical = make([]string, 0, len(s.data))
	for _, row := range s.data {
		if index >= len(row) {
			vertical = append(vertical, "")
		} else {
			vertical = append(vertical, row[index])
		}
	}
	return vertical
}

func (s *SectionBase) Format(int) {

}

// 删除每个格子的左右空格
func (s *SectionBase) TrimCells() {
	for i := 0; i < len(s.data); i++ {
		for j := 0; j < len(s.data[i]); j++ {
			s.data[i][j] = strings.Trim(s.data[i][j], " ")  // 空格
			s.data[i][j] = strings.Trim(s.data[i][j], "	")  // Tab
			s.data[i][j] = strings.Trim(s.data[i][j], "\n") // 换行
			s.data[i][j] = strings.Trim(s.data[i][j], "\r") // 换行
		}
	}
}

// 删除每行末尾的空格
func (s *SectionBase) TrimColumns() {
	for i := 0; i < len(s.data); i++ {
		for len(s.data[i]) > 0 {
			if s.data[i][len(s.data[i])-1] != "" {
				break
			}
			s.data[i] = s.data[i][:len(s.data[i])-1]
		}
	}
}

// 删除整行为空的行
func (s *SectionBase) CutRows() {
	for i := 0; i < len(s.data); i++ {
		row := s.data[i]
		isEmpty := true
		for _, cell := range row {
			if cell != "" {
				isEmpty = false
				break
			}
		}

		if isEmpty {
			s.data = append(s.data[:i], s.data[i+1:]...)
			i--
		}
	}
}

// 按照最长的行填补列
func (s *SectionBase) FillColumns(headerSize int) {
	var maxSize int
	if headerSize > 0 {
		maxSize = headerSize
	} else {
		for i := 0; i < len(s.data); i++ {
			if maxSize < len(s.data[i]) {
				maxSize = len(s.data[i])
			}
		}
	}
	for i := 0; i < len(s.data); i++ {
		for j := len(s.data[i]); j < maxSize; j++ {
			s.data[i] = append(s.data[i], "")
		}
	}
}

// 行列互转
func (s *SectionBase) Rotate() {
	var maxSize int
	for i := 0; i < len(s.data); i++ {
		if maxSize < len(s.data[i]) {
			maxSize = len(s.data[i])
		}
	}
	var data = make([][]string, maxSize)
	for i := 0; i < len(data); i++ {
		data[i] = make([]string, len(s.data))
	}
	for i := 0; i < len(s.data); i++ {
		for j := 0; j < len(s.data[i]); j++ {
			data[j][i] = s.data[i][j]
		}
	}
	s.data = data

	lineMap := make(map[interface{}]int)
	for _, index := range s.indexMap {
		lineMap[index+1] = index
		lineMap[format.ToExcelCol(index+1)] = index
	}
	s.indexMap = lineMap
}

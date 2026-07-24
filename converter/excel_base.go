package converter

import (
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExcelBase struct {
	relPath  string
	file     *excelize.File
	sheetMap map[SheetType]map[string]Sheet
	nodes    []Node
}

func NewExcelBase(path Path, relPath string) *ExcelBase {
	return &ExcelBase{
		relPath:  relPath,
		sheetMap: map[SheetType]map[string]Sheet{},
	}
}

func (e *ExcelBase) String() string {
	return e.relPath
}

func (e *ExcelBase) Read() {

}

func (e *ExcelBase) Preprocess() {
	for _, sheet := range e.sheetMap[SheetTypeRegular] {
		sheet.FormatHeader()
	}
	for _, sheet := range e.sheetMap[SheetTypeInferior] {
		sheet.FormatHeader()
	}
	for _, sheet := range e.sheetMap[SheetTypeSettings] {
		sheet.FormatHeader()
	}
}

func (e *ExcelBase) SheetMap() map[SheetType]map[string]Sheet {
	return e.sheetMap
}

func (e *ExcelBase) ReadFile() map[string][][]string {
	return e.ReadFileRows(0)
}

func (e *ExcelBase) ReadFileRows(maxRows int) map[string][][]string {
	var data = make(map[string][][]string)
	var err error
	e.file, err = excelize.OpenFile(filepath.Join(path.ImportAbsPath(), e.relPath))
	if err != nil {
		Exit("[%v] Read %v error, %v", e, e, err)
	}
	defer func() {
		// All generated sheets and nodes retain only the extracted row data.
		// Drop the parsed workbook promptly so a full-repository conversion
		// doesn't retain thousands of excelize workbooks until process exit.
		e.file = nil
	}()
	if err := e.formatAllCellsAsText(); err != nil {
		Exit("[%v] Format cells as text error, %v", e, err)
	}
	sheetMap := e.file.GetSheetMap()
	for _, sheetName := range sheetMap {
		rows, err := e.readSheetRows(sheetName, maxRows)
		if err != nil {
			Exit("[%v] Read %v error, %v", e, sheetName, err)
		}
		data[sheetName] = rows
	}
	return data
}

func (e *ExcelBase) readSheetRows(sheetName string, maxRows int) ([][]string, error) {
	iter, err := e.file.Rows(sheetName)
	if err != nil {
		return nil, err
	}
	capacity := maxRows
	if capacity <= 0 {
		capacity = 64
	}
	rows := make([][]string, 0, capacity)
	for iter.Next() {
		row, err := iter.Columns()
		if err != nil {
			return rows, err
		}
		rows = append(rows, row)
		if maxRows > 0 && len(rows) >= maxRows {
			break
		}
	}
	if err := iter.Error(); err != nil {
		return rows, err
	}
	return rows, nil
}

// formatAllCellsAsText performs the programmatic equivalent of selecting all
// populated cells in Excel and converting them to Text before import. Numeric
// values are committed as text at Excel's 15-significant-digit precision so
// binary floating-point tails are not emitted into generated source files.
func (e *ExcelBase) formatAllCellsAsText() error {
	textStyle, err := e.file.NewStyle(&excelize.Style{NumFmt: 49})
	if err != nil {
		return err
	}
	for _, sheetName := range e.file.GetSheetMap() {
		// GetCellValue loads the worksheet into excelize.File.Sheet so its
		// original cell type and stored value are available below.
		if _, err := e.file.GetCellValue(sheetName, "A1"); err != nil {
			return err
		}
	}
	for _, sheet := range e.file.Sheet {
		if sheet == nil {
			continue
		}
		for rowIndex := range sheet.SheetData.Row {
			row := &sheet.SheetData.Row[rowIndex]
			for cellIndex := range row.C {
				cell := &row.C[cellIndex]
				cell.S = textStyle
				if cell.V != "" && (cell.T == "" || cell.T == "n") {
					cell.V = normalizeExcelNumber(cell.V)
					cell.T = "str"
				}
			}
		}
	}
	return nil
}

func normalizeExcelNumber(value string) string {
	if !strings.ContainsAny(value, ".eE") {
		return value
	}
	number, _, err := big.ParseFloat(value, 10, 256, big.ToNearestEven)
	if err != nil {
		return value
	}
	return number.Text('g', 15)
}

func (e *ExcelBase) PackageName() string {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	if len(strs) < 3 {
		Exit("[%v] Error rel path", e)
	}
	return strs[0]
}

func (e *ExcelBase) DomainName() string {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	length := len(strs)
	switch {
	case length < 3:
		Exit("[%v] Error rel path", e)
	case length == 3: // excel当作Domain
		return e.FixedName()
	case length > 3:
		return format.ToUpper(strs[2])
	}
	return FlagDefault
}

func (e *ExcelBase) IndirectName() string {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	length := len(strs)
	switch {
	case length < 3:
		Exit("[%v] Error rel path", e)
	case length == 3: // excel当作Domain
		return FlagDefault
	case length > 3:
		return e.FixedName()
	}
	return FlagDefault
}

func (e *ExcelBase) Category() string {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	length := len(strs)
	switch {
	case length < 3:
		Exit("[%v] Error rel path", e)
	}
	return format.ToUpper(strs[0])
}

func (e *ExcelBase) Name() string {
	return filepath.Base(e.relPath)
}

func (e *ExcelBase) FixedName() string {
	return format.ToUpper(e.Name())
}

func (e *ExcelBase) Type() string {
	Exit("[Main] Invalid call ExcelBase.Type")
	return ""
}

func (e *ExcelBase) GetSheet(sheetName string) Sheet {
	for _, sheetMap := range e.sheetMap {
		if sheet, ok := sheetMap[sheetName]; ok {
			return sheet
		}
	}
	Exit("[%v] Sheet %s not found", e, sheetName)
	return nil
}

func (e *ExcelBase) GetHeaderSize(sheetName string) int {
	sheet := e.GetSheet(sheetName)
	return sheet.HeaderSize()
}

func (e *ExcelBase) GetHeaderField(sheetName string, key any) HeaderField {
	sheet := e.GetSheet(sheetName)
	return sheet.GetHeaderField(key)
}

func (e *ExcelBase) Build() {

}

func (e *ExcelBase) Nodes() []Node {
	return e.nodes
}

func (e *ExcelBase) ForServer() bool {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	for _, str := range strs {
		if format.SpecificClient(str) {
			return false
		}
	}
	return true
}

func (e *ExcelBase) ForClient() bool {
	strs := strings.Split(e.relPath, string(os.PathSeparator))
	for _, str := range strs {
		if format.SpecificServer(str) {
			return false
		}
	}
	return true
}

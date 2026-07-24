package converter

import (
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func TestFormatAllCellsAsTextBeforeImport(t *testing.T) {
	file := excelize.NewFile()
	file.NewSheet("Second")
	if err := file.SetCellValue("Sheet1", "A1", 2.4000000000000004); err != nil {
		t.Fatal(err)
	}
	if err := file.SetCellValue("Sheet1", "B1", "001234567890123456"); err != nil {
		t.Fatal(err)
	}
	if err := file.SetCellValue("Sheet1", "C3", 240.00000000000003); err != nil {
		t.Fatal(err)
	}
	if err := file.SetCellValue("Second", "B2", 7.000000000000001); err != nil {
		t.Fatal(err)
	}
	numericStyle, err := file.NewStyle(&excelize.Style{NumFmt: 2})
	if err != nil {
		t.Fatal(err)
	}
	if err := file.SetCellStyle("Sheet1", "A1", "A1", numericStyle); err != nil {
		t.Fatal(err)
	}

	excel := NewExcelBase(nil, "test.xlsx")
	excel.file = file
	if err := excel.formatAllCellsAsText(); err != nil {
		t.Fatal(err)
	}

	textStyle, err := file.GetCellStyle("Sheet1", "A1")
	if err != nil {
		t.Fatal(err)
	}
	for _, cell := range []struct {
		sheet string
		axis  string
	}{
		{sheet: "Sheet1", axis: "A1"},
		{sheet: "Sheet1", axis: "B1"},
		{sheet: "Sheet1", axis: "C3"},
		{sheet: "Second", axis: "B2"},
	} {
		style, err := file.GetCellStyle(cell.sheet, cell.axis)
		if err != nil {
			t.Fatal(err)
		}
		if style != textStyle {
			t.Fatalf("cell %s!%s style = %d, want text style %d", cell.sheet, cell.axis, style, textStyle)
		}
	}

	rows, err := excel.readSheetRows("Sheet1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if got := rows[0][0]; got != "2.4" {
		t.Fatalf("imported decimal = %q, want normalized text 2.4", got)
	}
	if got := rows[0][1]; got != "001234567890123456" {
		t.Fatalf("original text = %q, want leading-zero text unchanged", got)
	}
	if got := rows[2][2]; got != "240" {
		t.Fatalf("imported decimal = %q, want normalized text 240", got)
	}
}

func TestFormatAllCellsAsTextRunsBeforeLimitedImport(t *testing.T) {
	file := excelize.NewFile()
	for row := 1; row <= 3; row++ {
		axis, err := excelize.CoordinatesToCellName(1, row)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.SetCellValue("Sheet1", axis, row); err != nil {
			t.Fatal(err)
		}
	}

	excel := NewExcelBase(nil, "test.xlsx")
	excel.file = file
	if err := excel.formatAllCellsAsText(); err != nil {
		t.Fatal(err)
	}
	rows, err := excel.readSheetRows("Sheet1", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d imported rows, want 2", len(rows))
	}
	firstStyle, err := file.GetCellStyle("Sheet1", "A1")
	if err != nil {
		t.Fatal(err)
	}
	thirdStyle, err := file.GetCellStyle("Sheet1", "A3")
	if err != nil {
		t.Fatal(err)
	}
	if thirdStyle != firstStyle {
		t.Fatalf("row outside import limit style = %d, want text style %d", thirdStyle, firstStyle)
	}
}

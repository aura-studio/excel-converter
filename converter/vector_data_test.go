package converter

import "testing"

func TestVectorExcelReadLimits(t *testing.T) {
	tests := []struct {
		name        string
		excel       string
		vectorEgUse bool
		maxRows     int
	}{
		{
			name:  "VectorEg exports all data",
			excel: "S_VectorEg.xlsx",
		},
		{
			name:  "VectorEgPP exports all data",
			excel: "S_VectorEgPP.xlsx",
		},
		{
			name:        "VectorEgUse exports schema only",
			excel:       "FortuneTigerVectorEgUse.xlsx",
			vectorEgUse: true,
			maxRows:     schemaOnlyExcelMaxRows,
		},
		{
			name:  "unrelated EgUse exports all data",
			excel: "FeatureEgUse.xlsx",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isVectorEgUseExcel(test.excel); got != test.vectorEgUse {
				t.Fatalf("isVectorEgUseExcel(%q) = %v, want %v", test.excel, got, test.vectorEgUse)
			}
			if got := excelReadMaxRows(test.excel); got != test.maxRows {
				t.Fatalf("excelReadMaxRows(%q) = %d, want %d", test.excel, got, test.maxRows)
			}
		})
	}
}

func TestCommentedVectorEgIsClassifiedAsCommentExcel(t *testing.T) {
	converter := NewConverter()
	for _, excelName := range []string{"#S_VectorEg.xlsx", "#S_VectorEgPP.xlsx"} {
		if got := converter.ExcelType(excelName); got != ExcelTypeComment {
			t.Fatalf("ExcelType(%q) = %q, want %q", excelName, got, ExcelTypeComment)
		}
	}
}

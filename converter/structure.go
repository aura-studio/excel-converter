package converter

type StructureType string

const (
	StructureTypeEmpty StructureType = ""
	// 基础类型
	StructureTypeString   StructureType = "string"
	StructureTypeInt      StructureType = "int"
	StructureTypeFloat    StructureType = "float"
	StructureTypeBool     StructureType = "bool"
	StructureTypeBigInt   StructureType = "bigint"
	StructureTypeBigRat   StructureType = "bigrat"
	StructureTypeBigFloat StructureType = "bigfloat"
	StructureTypeTime     StructureType = "time"
	// 显式结构 int 指代基础类型，struct 指代所有类型
	StructureTypeArray StructureType = "array" // internal: golang []int
	StructureTypeTable StructureType = "table" // internal: golang map[index]int
	StructureTypeList  StructureType = "list"  // internal: golang []struct
	StructureTypeDict  StructureType = "dict"  // internal: golang struct
	// 隐式结构 int 指代基础类型，struct 指代所有类型
	StructureTypeStruct    StructureType = "struct"    // external：golang struct
	StructureTypeRow       StructureType = "row"       // external：golang []int
	StructureTypeCol       StructureType = "col"       // external：golang []int
	StructureTypeMapVal    StructureType = "mapval"    // external：golang map[index] = struct
	StructureTypeRowMapVal StructureType = "rowmapval" // external：golang map[index] = []int
	StructureTypeColMapVal StructureType = "colmapval" // external：golang map[index] = []int
	// Sheet类型 & 隐式结构 int 指代基础类型，struct 指代所有类型
	StructureTypeStructs StructureType = "structs" // external：golang []struct
	StructureTypeRows    StructureType = "rows"    // external：golang [][]int
	StructureTypeCols    StructureType = "cols"    // external：golang [][]int
	StructureTypeMap     StructureType = "map"     // external：golang map[index]struct
	StructureTypeRowMap  StructureType = "rowmap"  // external：golang map[index][]int
	StructureTypeColMap  StructureType = "colmap"  // external：golang map[index][]int
)

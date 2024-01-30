package converter

import (
	"fmt"
	"sort"
	"strings"
)

const (
	CollectionHeaderSrcExcel = "src_path"
	CollectionHeaderDstExcel = "dst_path"
	CollectionHeaderSrcSheet = "src_sheet"
	CollectionHeaderDstSheet = "dst_sheet"
)

type StoragePath struct {
	PackageName string
	ExcelName   string
	SheetName   string
}

func (s *StoragePath) String() string {
	return fmt.Sprintf("%s:%s:%s", s.PackageName, s.ExcelName, s.SheetName)
}

type StorageVar struct {
	PackageName string
	VarName     string
}

type LinkPath struct {
	ExcelName string
	SheetName string
}

func (l *LinkPath) String() string {
	return fmt.Sprintf("%s:%s", l.ExcelName, l.SheetName)
}

type Link struct {
	SrcLinkPath LinkPath
	DstLinkPath LinkPath
}

type Storage struct {
	StoragePath StoragePath
	StorageVar  StorageVar
}

type Collection struct {
	storageMap map[StoragePath]StorageVar
	linkMap    map[LinkPath]LinkPath
	categories []string
}

func NewCollection() *Collection {
	return &Collection{
		storageMap: make(map[StoragePath]StorageVar),
		linkMap:    make(map[LinkPath]LinkPath),
		categories: []string{FlagBase},
	}
}

func (l *Collection) PackageNames() []string {
	var packageNames = make([]string, 0, len(l.storageMap))
	var packageMap = make(map[string]bool)
	for storagePath := range l.storageMap {
		packageMap[storagePath.PackageName] = true
	}
	for packageName := range packageMap {
		packageNames = append(packageNames, packageName)
	}
	sort.Strings(packageNames)
	return packageNames
}

func (l *Collection) Storages() []*Storage {
	var storages = make([]*Storage, 0, len(l.storageMap))
	for storagePath, storageVar := range l.storageMap {
		storages = append(storages, &Storage{
			StoragePath: storagePath,
			StorageVar:  storageVar,
		})
	}
	sort.Slice(storages, func(i, j int) bool {
		return storages[i].StoragePath.String() < storages[j].StoragePath.String()
	})
	return storages
}

func (l *Collection) Categories() []string {
	return l.categories
}

func (l *Collection) Links() []*Link {
	var links = make([]*Link, 0, len(l.linkMap))
	for dstLink, srcLink := range l.linkMap {
		links = append(links, &Link{
			SrcLinkPath: srcLink,
			DstLinkPath: dstLink,
		})
	}
	sort.Slice(links, func(i, j int) bool {
		return links[i].DstLinkPath.String() < links[j].DstLinkPath.String()
	})
	return links
}

func (l *Collection) ReadNode(node Node) {
	packageName := node.Excel().PackageName()
	varName := node.RootName()
	excelName := node.ExcelPathName()
	sheetName := node.SheetPathName()
	storagePath := StoragePath{packageName, excelName, sheetName}
	storageVar := StorageVar{packageName, varName}
	l.storageMap[storagePath] = storageVar
}

func (l *Collection) ReadLink(sheet Sheet) {
	switch sheet.Name() {
	case FlagLink:
		for index := 0; index < sheet.VerticleSize(); index++ {
			srcExcelName := sheet.GetCell(sheet.GetHeaderField(format.ToUpper(CollectionHeaderSrcExcel)), index)
			dstExcelName := sheet.GetCell(sheet.GetHeaderField(format.ToUpper(CollectionHeaderDstExcel)), index)
			srcSheetName := sheet.GetCell(sheet.GetHeaderField(format.ToUpper(CollectionHeaderSrcSheet)), index)
			dstSheetName := sheet.GetCell(sheet.GetHeaderField(format.ToUpper(CollectionHeaderDstSheet)), index)

			srcLinkPath := LinkPath{srcExcelName, srcSheetName}
			dstLinkPath := LinkPath{dstExcelName, dstSheetName}

			l.linkMap[dstLinkPath] = srcLinkPath
		}
	case FlagVarian:
		for horizonIndex := 1; horizonIndex < sheet.HeaderSize(); horizonIndex++ {
			for verticleIndex := 1; verticleIndex < sheet.VerticleSize(); verticleIndex++ {
				srcExcelName := sheet.GetHeaderField(horizonIndex).Name
				dstExcelName := sheet.GetHeaderField(horizonIndex).Name
				srcSheetName := sheet.GetCell(sheet.GetHeaderField(horizonIndex), verticleIndex)
				dstSheetName := fmt.Sprintf("%s/%s", sheet.GetCell(sheet.GetHeaderField(horizonIndex), 0), sheet.GetCell(sheet.GetHeaderField(0), verticleIndex))

				srcLinkPath := LinkPath{srcExcelName, srcSheetName}
				dstLinkPath := LinkPath{dstExcelName, dstSheetName}

				l.linkMap[dstLinkPath] = srcLinkPath
			}
		}
	case FlagCategory:
		var keys []string
		var index = 0
		for index < sheet.HeaderSize() {
			layerKeys := sheet.GetHorizon(index)
			if len(layerKeys) == 0 {
				break
			}

			if index == 0 {
				keys = layerKeys
			} else {
				var newKeys []string
				for _, key := range keys {
					for _, layerKey := range layerKeys {
						newKeys = append(newKeys, fmt.Sprintf("%s_%s", key, format.ToUpper(layerKey)))
					}
				}
				keys = newKeys
			}

			for _, key := range keys {
				l.categories = append(l.categories, key)
				Debug("[%s] category found, %v", key, strings.Split(key, "_"))
			}

			index++
		}

		sort.Strings(l.categories)
	}
}

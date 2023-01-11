package converter

import (
	"os"
	"path/filepath"
)

type Path interface {
	Path() string
	Abs(relPath string) string
	Rel(path string) string
	ImportAbsPath() string
	ImportRelPath() string
	ExportAbsPath() string
	ExportRelPath() string
}

type RootPath struct {
	root          string
	relImportPath string
	relExportPath string
}

var path = NewRootPath()

func NewRootPath() *RootPath {
	return &RootPath{}
}

func (p *RootPath) Init(relImportPath, relExportPath string) {
	p.root = p.findRoot()
	p.relImportPath = relImportPath
	p.relExportPath = relExportPath
}

func (p *RootPath) Dirname() string {
	return filepath.Base(p.root)
}

func (p *RootPath) Path() string {
	return p.root
}

func (p *RootPath) ImportAbsPath() string {
	return p.Abs(p.relImportPath)
}

func (p *RootPath) ImportRelPath() string {
	return p.relImportPath
}

func (p *RootPath) ExportAbsPath() string {
	return p.Abs(p.relExportPath)
}

func (p *RootPath) ExportRelPath() string {
	return p.relExportPath
}

func (p *RootPath) findRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		Exit("[Main] find root error, %v", err)
	}

	return cwd
}

func (p *RootPath) Abs(relPath string) string {
	return filepath.Join(p.root, relPath)
}

func (p *RootPath) Rel(path string) string {
	relPath, err := filepath.Rel(p.root, path)
	if err != nil {
		Exit("[Main] Get rel file path error %v", path)
	}
	return relPath
}

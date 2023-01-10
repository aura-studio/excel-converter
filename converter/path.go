package converter

import (
	"os"
	"path/filepath"
	"rocket-nano/internal/util/system"
)

const (
	RelPathExcel    = "excel"
	RelPathExternal = "external"
	RelPathGo       = "../internal/exported/config"
	RelPathLua      = "Client_Config/Config"
)

type Path interface {
	Abs(relPath string) string
	Rel(absPath string) string
}

type PathExternal struct {
	root string
}

func NewPathExternal() *PathExternal {
	p := &PathExternal{}
	p.root = p.findRoot()
	return p
}

func (p *PathExternal) findRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		Exit("[Main] find root error, %v", err)
	}

	possibleRoots := []string{
		cwd,
		filepath.Join(p.findProjectDirectory(), RelPathExternal),
	}

	for _, root := range possibleRoots {
		if system.Exists(filepath.Join(root, RelPathExcel)) {
			return root
		}
	}
	Exit("[Main] Root path not found, possible roots: %v", possibleRoots)
	return ""
}

func (*PathExternal) findProjectDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		Exit("[Main] get cwd error: %v", err)
	}
	projectDirectory, err := system.MatchParentDir(workingDirectory, "rocket-nano")
	if err != nil {
		return workingDirectory
	}
	return projectDirectory
}

func (p *PathExternal) Abs(relPath string) string {
	return filepath.Join(p.root, relPath)
}

func (p *PathExternal) Rel(path string) string {
	relPath, err := filepath.Rel(p.root, path)
	if err != nil {
		Exit("[Main] Get rel file path error %v", path)
	}
	return relPath
}

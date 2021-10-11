package nodejs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type NodeResolver struct {
	Extensions []string
	Modules    []string
	MainFields []string
}

func (r *NodeResolver) SetDefaults() {
	if r.MainFields == nil {
		r.MainFields = []string{
			"browser",
			"jsnext:main",
			"module",
			"main",
		}
	}
	if r.Modules == nil {
		r.Modules = []string{
			"node_modules",
		}
	}

	if r.Extensions == nil {
		r.Extensions = []string{".tsx", ".ts", ".mjs", ".js"}
	}
}

func (r *NodeResolver) Resolve(path string, at string) (string, bool) {
	if path[0] == '/' {
		full := path
		if fi, exists := pathExists(full); exists {
			if fi.IsDir() {
				return r.resolveInDir(full)
			}
			return full, true
		}
		return "", false
	}

	if strings.HasPrefix(path, "./") {
		full := filepath.Join(filepath.Dir(at), path)
		if fi, exists := pathExists(full); exists {
			if fi.IsDir() {
				return r.resolveInDir(full)
			}
			return full, true
		}
		return "", false
	}

	for i := range r.Modules {
		m := r.Modules[i]

		if m[0] == '/' {
			if found, exists := r.Resolve(filepath.Join(m, path), at); exists {
				return found, exists
			}
			continue
		}

		pkgRoot := filepath.Dir(at)

		for pkgRoot != "/" {
			packageJSONFile := filepath.Join(pkgRoot, "package.json")
			if _, packageJSONExists := pathExists(packageJSONFile); packageJSONExists {
				p := filepath.Join(pkgRoot, m, path)
				if found, exists := r.Resolve(p, at); exists {
					return found, exists
				}
			}
			pkgRoot = filepath.Join(pkgRoot, "../")
		}
	}

	return "", false
}

func (r *NodeResolver) resolveInDir(dir string) (string, bool) {
	packageJSONFile := filepath.Join(dir, "package.json")
	if _, packageJSONExists := pathExists(packageJSONFile); packageJSONExists {
		data, _ := os.ReadFile(packageJSONFile)
		pkgJSON := make(map[string]interface{})
		if err := json.Unmarshal(data, &pkgJSON); err != nil {
			panic(err)
		}
		for _, mainField := range r.MainFields {
			if mf, ok := pkgJSON[mainField]; ok {
				if mainFilePath, ok := mf.(string); ok {
					mainFile := filepath.Join(dir, mainFilePath)
					if _, mainFileExists := pathExists(mainFile); mainFileExists {
						return mainFile, true
					}
				}
			}
		}
	} else {
		for i := range r.Extensions {
			indexFile := filepath.Join(dir, "index"+r.Extensions[i])
			if _, indexFileExists := pathExists(indexFile); indexFileExists {
				return indexFile, true
			}
		}
	}
	return "", false
}

func pathExists(path string) (os.FileInfo, bool) {
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false
		}
		panic(err)
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		r, err := os.Readlink(path)
		if err != nil {
			return nil, false
		}
		return pathExists(filepath.Join(filepath.Dir(path), r))
	}
	return fi, true
}

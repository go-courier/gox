package core

import (
	"container/list"
	"crypto/sha1"
	"fmt"
	"path/filepath"
)

type SourceFiles struct {
	list *list.List
	m    map[string]*SourceFile
}

func (sfs *SourceFiles) PutSourceFile(sf *SourceFile) {
	if sfs.list == nil {
		sfs.list = list.New()
	}
	if sfs.m == nil {
		sfs.m = map[string]*SourceFile{}
	}

	if _, ok := sfs.m[sf.filename]; !ok {
		sfs.m[sf.filename] = sf
		sfs.list.PushBack(sf)
	}
}

type ManifestConvert interface {
	OutputRequestURL(f string) string
}

func (sfs *SourceFiles) ListManifests(entrypoint string, c ManifestConvert) []Manifest {
	manifests := make([]Manifest, 0, sfs.list.Len())

	for e := sfs.list.Front(); e != nil; e = e.Next() {
		sf := e.Value.(*SourceFile)

		url := c.OutputRequestURL(sf.OutputFilename)

		if entrypoint == sf.Filename() {
			url = filepath.Base(entrypoint)
		}

		manifests = append(manifests, Manifest{
			URL:        url,
			SourceFile: sf,
		})
	}

	return manifests
}

func NewSourceFile(filename string, raw []byte, async bool) *SourceFile {
	return &SourceFile{
		filename: filename,
		raw:      raw,
		Async:    async,
	}
}

type SourceFile struct {
	filename string
	raw      []byte

	Async          bool
	OutputFilename string
	Transformed    Output
}

func (sf *SourceFile) Filename() string {
	if sf.Transformed != nil {
		return sf.Transformed.Filename()
	}
	return sf.filename
}

func (sf *SourceFile) Code() []byte {
	return sf.raw
}

func (sf *SourceFile) OutputCode() []byte {
	if sf.Transformed != nil {
		return sf.Transformed.Code()
	}
	return sf.raw
}

func hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func fileRename(filename string, rename func(base string) string) string {
	ext := filepath.Ext(filename)
	base := filepath.Base(filename)
	return filepath.Join(filepath.Dir(filename), rename(base[0:len(base)-len(ext)])+ext)
}

func OutputFilename(o Output) string {
	return fileRename(o.Filename(), func(base string) string {
		return base + "-" + hash(o.Code())
	})
}

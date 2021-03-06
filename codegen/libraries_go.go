package codegen

import (
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/raml"
	log "github.com/Sirupsen/logrus"
)

// library defines an RAML library
// it is implemented as package in Go
type goLibrary struct {
	*raml.Library
	PackageName string
	baseDir     string // root directory
	dir         string // library directory
}

// create new library instance
func newGoLibrary(lib *raml.Library, baseDir string) *goLibrary {
	l := goLibrary{
		Library: lib,
		baseDir: baseDir,
	}

	// package name  : base of filename without the extension
	l.PackageName = strings.TrimSuffix(filepath.Base(l.Filename), filepath.Ext(l.Filename))
	l.PackageName = normalizePkgName(l.PackageName)

	// package directory : filename without the extension
	relDir := libRelDir(l.Filename)
	l.dir = normalizePkgName(filepath.Join(l.baseDir, relDir))

	return &l
}

// generate code of all libraries
func generateLibraries(libraries map[string]*raml.Library, baseDir string) error {
	for _, ramlLib := range libraries {
		l := newGoLibrary(ramlLib, baseDir)
		if err := l.generate(); err != nil {
			return err
		}
	}
	return nil
}

// generate code of this library
func (l *goLibrary) generate() error {
	if err := checkCreateDir(l.dir); err != nil {
		return err
	}

	// generate all Type structs
	if err := generateStructs(l.Types, l.dir, l.PackageName, langGo); err != nil {
		return err
	}

	// security schemes
	if err := generateSecurity(l.SecuritySchemes, l.dir, l.PackageName, langGo); err != nil {
		return err
	}

	// included libraries
	for _, ramlLib := range l.Libraries {
		childLib := newGoLibrary(ramlLib, l.baseDir)
		if err := childLib.generate(); err != nil {
			return err
		}
	}
	return nil
}

// get relative lib directory from library filename
// this relative directory will be used as:
// - lib package name
// - lib files directory
func libRelDir(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// get library import path from a type
func libImportPath(rootImportPath, typ string) string {
	// library use '.', return nothing if it is not a library
	if strings.Index(typ, ".") < 0 {
		return ""
	}

	// library name in the current document
	libName := strings.Split(typ, ".")[0]

	if libName == "goraml" { // special package name, reserved for goraml
		return filepath.Join(rootImportPath, "goraml")
	}

	// raml file of this lib
	libRAMLFile := globAPIDef.FindLibFile(denormalizePkgName(libName))

	if libRAMLFile == "" {
		log.Fatalf("can't find library : %v", libName)
	}

	// relative lib package
	libPkg := libRelDir(libRAMLFile)

	return filepath.Join(rootImportPath, normalizePkgName(libPkg))
}

// normalize package name because not all characters can be used as package name
func normalizePkgName(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

// inverse of normalizePkgName
func denormalizePkgName(name string) string {
	return strings.Replace(name, "_", "-", -1)
}

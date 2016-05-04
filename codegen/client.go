package codegen

import (
	"strings"

	"github.com/Jumpscale/go-raml/raml"
)

const (
	clientResourceTemplate       = "./templates/client_resource.tmpl"
	clientHelperResourceTemplate = "./templates/client_helper_resource.tmpl"
)

// API client definition
type clientDef struct {
	Name    string
	BaseURI string
	Methods []methodInterface
}

// create client definition from RAML API definition
func newClientDef(apiDef *raml.APIDefinition) clientDef {
	cd := clientDef{
		Name:    normalizeURI(apiDef.Title),
		BaseURI: apiDef.BaseURI,
	}
	if strings.Index(cd.BaseURI, "{version}") > 0 {
		cd.BaseURI = strings.Replace(cd.BaseURI, "{version}", apiDef.Version, -1)
	}
	return cd
}

// GenerateClient generates client library
func GenerateClient(apiDef *raml.APIDefinition, dir, lang string) error {
	//check create dir
	if err := checkCreateDir(dir); err != nil {
		return err
	}

	// creates base client struct
	cd := newClientDef(apiDef)

	for _, v := range apiDef.Resources {
		rd := newResourceDef(apiDef, normalizeURITitle(apiDef.Title), "main")
		rd.generateMethods(&v, lang)
		cd.Methods = append(cd.Methods, rd.Methods...)
	}

	switch lang {
	case langGo:
		gc := goClient{clientDef: cd}
		return gc.generate(apiDef, dir)
	case langPython:
		pc := pythonClient{clientDef: cd}
		return pc.generate(dir)
	}
	return errInvalidLang
}

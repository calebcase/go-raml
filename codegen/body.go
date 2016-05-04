package codegen

import (
	"github.com/Jumpscale/go-raml/raml"
)

const (
	reqBodySuffix  = "ReqBody"
	respBodySuffix = "RespBody"
)

// generate all body struct from an RAML definition
func generateBodyStructs(apiDef *raml.APIDefinition, dir, packageName, lang string) error {
	// generate
	for _, v := range apiDef.Resources {
		if err := generateStructsFromResourceBody("", dir, packageName, lang, &v); err != nil {
			return err
		}
	}

	return nil
}

// generate all structs from resource's method's request & response body
func generateStructsFromResourceBody(resourcePath, dir, packageName, lang string, r *raml.Resource) error {
	if r == nil {
		return nil
	}

	structName := normalizeURITitle(resourcePath + r.URI)

	// build
	var methods = []struct {
		Name   string
		Method *raml.Method
	}{
		{Name: "Get", Method: r.Get},
		{"Post", r.Post},
		{"Head", r.Head},
		{"Put", r.Put},
		{"Delete", r.Delete},
		{"Patch", r.Patch},
	}
	for _, v := range methods {
		if err := buildBodyFromMethod(structName, v.Name, dir, packageName, lang, v.Method); err != nil {
			return err
		}
	}

	// build request/response body of child resources
	for _, v := range r.Nested {
		if err := generateStructsFromResourceBody(resourcePath+r.URI, dir, packageName, lang, v); err != nil {
			return err
		}
	}

	return nil
}

// build request and reponse body of a method.
// in python case, we only need to build it for request body because we only need it for validator
func buildBodyFromMethod(structName, methodName, dir, packageName, lang string, method *raml.Method) error {
	if method == nil {
		return nil
	}

	//generate struct for request body
	switch lang {
	case langGo:
		if err := generateStructFromBody(structName+methodName, dir, packageName, &method.Bodies, true); err != nil {
			return err
		}
	case langPython:
		if !hasJSONBody(&method.Bodies) {
			return nil
		}
		pc := newPythonClass(structName+methodName+reqBodySuffix, "", method.Bodies.ApplicationJSON.Properties)
		return pc.generate(dir)
	}

	//generate struct for response body
	for _, val := range method.Responses {
		if err := generateStructFromBody(structName+methodName, dir, packageName, &val.Bodies, false); err != nil {
			return err
		}

	}

	return nil
}

// check if this raml.Bodies has JSON body that need to be generated
func hasJSONBody(body *raml.Bodies) bool {
	return body.ApplicationJSON != nil && len(body.ApplicationJSON.Properties) > 0
}

// generate a struct from an RAML request/response body
func generateStructFromBody(structNamePrefix, dir, packageName string, body *raml.Bodies, isGenerateRequest bool) error {
	if !hasJSONBody(body) {
		return nil
	}

	// construct struct from body
	structDef := newStructDefFromBody(body, structNamePrefix, packageName, isGenerateRequest)

	// generate
	return structDef.generate(dir)
}

package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/raml"
)

const (
	structTemplateLocation         = "./templates/struct.tmpl"
	inputValidatorTemplateLocation = "./templates/struct_input_validator.tmpl"
	inputValidatorFileResult       = "struct_input_validator.go"
	maxStringLen                   = 2147483647
)

// FieldDef defines a field of a struct
type fieldDef struct {
	Name          string // field name
	Type          string // field type
	IsComposition bool   // composition type
	IsOmitted     bool   // omitted empty
	UniqueItems   bool

	Validators string
}

func (fd *fieldDef) buildValidators(p raml.Property) {
	validators := ""
	// string
	if p.MinLength != nil {
		validators += fmt.Sprintf(",min=%v", *p.MinLength)
	}
	if p.MaxLength != nil {
		validators += fmt.Sprintf(",max=%v", *p.MaxLength)
	}
	if p.Pattern != nil {
		validators += fmt.Sprintf(",regexp=%v", *p.Pattern)
	}

	// Number
	if p.Minimum != nil {
		validators += fmt.Sprintf(",min=%v", *p.Minimum)
	}

	if p.Maximum != nil {
		validators += fmt.Sprintf(",max=%v", *p.Maximum)
	}

	if p.MultipleOf != nil {
		validators += fmt.Sprintf(",multipleOf=%v", *p.MultipleOf)
	}

	//if p.Format != nil {
	//}

	// Array & Map
	if p.MinItems != nil {
		validators += fmt.Sprintf(",min=%v", *p.MinItems)
	}
	if p.MaxItems != nil {
		validators += fmt.Sprintf(",max=%v", *p.MaxItems)
	}
	if p.UniqueItems {
		fd.UniqueItems = true
	}

	// Required
	if !fd.IsOmitted {
		validators += ",nonzero"
	}

	if validators != "" {
		fd.Validators = validators[1:]
	}
}

// StructDef defines a struct
type structDef struct {
	T           raml.Type           // raml.Type of this struct
	Name        string              // struct's name
	Description []string            // structs description
	PackageName string              // package name
	Fields      map[string]fieldDef // all struct's fields
	OneLineDef  string              // not empty if this struct can be defined in one line

	Validators []string
}

// true if this struct need to import 'fmt' package
func (sd structDef) NeedFmt() bool {
	// array type min items and max items
	if sd.T.MinItems > 0 || sd.T.MaxItems > 0 {
		return true
	}
	for _, f := range sd.Fields {
		if f.UniqueItems {
			return true
		}
	}
	return false
}

// true if this struct is not an alias of `interface{}`
func (sd structDef) NotBareInterface() bool {
	return !strings.HasSuffix(sd.OneLineDef, " interface{}")
}

// create new struct def
func newStructDef(name, packageName, description string, properties map[string]interface{}) structDef {
	// generate struct's fields from type properties
	fields := make(map[string]fieldDef)
	for k, v := range properties {
		prop := raml.ToProperty(k, v)
		fd := fieldDef{
			Name:      strings.Title(prop.Name),
			Type:      convertToGoType(prop.Type),
			IsOmitted: !prop.Required,
		}

		fd.buildValidators(prop)
		fields[prop.Name] = fd
	}
	return structDef{
		Name:        name,
		PackageName: packageName,
		Fields:      fields,
		Description: commentBuilder(description),
	}
}

// create struct definition from RAML Type node
func newStructDefFromType(t raml.Type, sName, packageName, lang string) structDef {
	sd := newStructDef(sName, packageName, t.Description, t.Properties)
	sd.T = t

	// handle advanced type on raml1.0
	sd.handleAdvancedType()

	return sd
}

// create struct definition from RAML Body node
func newStructDefFromBody(body *raml.Bodies, structNamePrefix, packageName string, isGenerateRequest bool) structDef {
	// set struct name based on request or response
	structName := structNamePrefix + respBodySuffix
	if isGenerateRequest {
		structName = structNamePrefix + reqBodySuffix
	}

	return newStructDef(structName, packageName, "", body.ApplicationJson.Properties)
}

// generate Go struct
func (sd structDef) generate(dir string) error {
	fileName := filepath.Join(dir, sd.Name+".go")
	if err := generateFile(sd, structTemplateLocation, "struct_template", fileName, false); err != nil {
		return err
	}
	return runGoFmt(fileName)
}

// generate all structs from an RAML api definition
func generateStructs(apiDefinition *raml.APIDefinition, dir, packageName, lang string) error {
	for name, t := range apiDefinition.Types {
		sd := newStructDefFromType(t, name, packageName, lang)
		if err := sd.generate(dir); err != nil {
			return err
		}
	}
	return nil
}

// handle advance type type into structField
// example:
//   Mammal:
//     type: Animal
//     properties:
//       name:
//         type: string
// the additional fieldDef would be Animal composition
func (sd *structDef) handleAdvancedType() {
	if sd.T.Type == nil {
		sd.T.Type = "object"
	}

	strType := interfaceToString(sd.T.Type)

	switch {
	case len(strings.Split(strType, ",")) > 1: //multiple inheritance
		sd.addMultipleInheritance(strType)
	case sd.T.IsUnion():
		sd.buildUnion()
	case sd.T.IsArray(): // arary type
		sd.buildArray()
	case strings.ToLower(strType) == "object": // plain type
		return
	case sd.T.IsEnum(): // enum
		sd.buildEnum()
	case len(sd.T.Properties) == 0: // specialization
		sd.buildSpecialization()
	default: // single inheritance
		sd.addSingleInheritance(strType)
	}
}

// add single inheritance
// inheritance is implemented as composition
// spec : http://docs.raml.org/specs/1.0/#raml-10-spec-inheritance-and-specialization
func (sd *structDef) addSingleInheritance(strType string) {
	fd := fieldDef{
		Name:          strType,
		IsComposition: true,
	}
	sd.Fields[strType] = fd

}

// construct multiple inheritance to Go type
// example :
// Anggora:
//	 type: [ Animal , Cat ]
//	 properties:
//		color:
//			type: string
// The additional fielddef would be a composition of Animal & Cat
// http://docs.raml.org/specs/1.0/#raml-10-spec-multiple-inheritance
func (sd *structDef) addMultipleInheritance(strType string) {
	for _, s := range strings.Split(strType, ",") {
		fieldType := strings.TrimSpace(s)
		fd := fieldDef{
			Name:          fieldType,
			IsComposition: true,
		}

		sd.Fields[fd.Name] = fd
	}
}

// buildEnum based on http://docs.raml.org/specs/1.0/#raml-10-spec-enums
// example result  `type TypeName []data_type`
func (sd *structDef) buildEnum() {
	if _, ok := sd.T.Type.(string); !ok {
		return
	}

	sd.buildOneLine(convertToGoType(sd.T.Type.(string)))
}

// build array type
// spec http://docs.raml.org/specs/1.0/#raml-10-spec-array-types
// example result  `type TypeName []something`
func (sd *structDef) buildArray() {
	sd.buildOneLine(convertToGoType(sd.T.Type.(string)))
}

// build union type
// union type is implemented as `interface{}`
// example result `type sometype interface{}`
func (sd *structDef) buildUnion() {
	sd.buildOneLine(convertUnion(sd.T.Type.(string)))
}

func (sd *structDef) buildSpecialization() {
	sd.buildOneLine(convertToGoType(sd.T.Type.(string)))
}

func (sd *structDef) buildOneLine(tipe string) {
	sd.OneLineDef = "type " + sd.Name + " " + tipe
}

// generate input validator helper file
func generateInputValidator(packageName, dir string) error {
	var ctx = struct {
		PackageName string
	}{
		PackageName: packageName,
	}
	fileName := filepath.Join(dir, inputValidatorFileResult)
	if err := generateFile(ctx, inputValidatorTemplateLocation, "struct_input_validator_template", fileName, true); err != nil {
		return err
	}
	return runGoFmt(fileName)
}

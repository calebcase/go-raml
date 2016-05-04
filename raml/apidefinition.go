package raml

// APIDefinition describes the basic information of an API, such as its
// title and base URI, and describes how to define common schema references.
type APIDefinition struct {
	RAMLVersion string `yaml:"-"`

	// A short, plain-text label for the API.
	Title string `yaml:"title" validate:"nonzero"`

	// The version of the API, for example "v1"
	Version string `yaml:"version"`

	// A URI that serves as the base for URIs of all resources.
	// Often used as the base of the URL of each resource containing the location of the API.
	// Can be a template URI.
	// The OPTIONAL baseUri property specifies a URI as an identifier for the API as a whole,
	// and MAY be used the specify the URL at which the API is served (its service endpoint),
	// and which forms the base of the URLs of each of its resources.
	// The baseUri property's value is a string that MUST conform to the URI specification RFC2396 or a Template URI.
	BaseURI string `yaml:"baseUri"`

	// Named parameters used in the baseUri (template).
	BaseURIParameters map[string]NamedParameter `yaml:"baseUriParameters"`

	// The protocols supported by the API.
	// The OPTIONAL protocols property specifies the protocols that an API supports.
	// If the protocols property is not explicitly specified, one or more protocols
	// included in the baseUri property is used;
	// if the protocols property is explicitly specified,
	// the property specification overrides any protocol included in the baseUri property.
	// The protocols property MUST be a non-empty array of strings, of values HTTP and/or HTTPS, and is case-insensitive.
	Protocols []string `yaml:"protocols"`

	// The default media types to use for request and response bodies (payloads),
	// for example "application/json".
	// Specifying the OPTIONAL mediaType property sets the default for return by API
	// requests having a body and for the expected responses. You do not need to specify the media type within every body definition.
	// The value of the mediaType property MUST be a sequence of
	// media type strings or a single media type string.
	// The media type applies to requests having a body,
	// the expected responses, and examples using the same sequence of media type strings.
	// Each value needs to conform to the media type specification in RFC6838.
	MediaType string `yaml:"mediaType"`

	// Additional overall documentation for the API.
	// The API definition can include a variety of documents that serve as a
	// user guides and reference documentation for the API. Such documents can
	// clarify how the API works or provide business context.
	// All the sections are in the order in which the documentation is declared.
	Documentation []Documentation `yaml:"documentation"`

	// An alias for the equivalent "types" property for compatibility with RAML 0.8.
	// Deprecated - API definitions should use the "types" property
	// because a future RAML version might remove the "schemas" alias for that property name.
	// The "types" property supports XML and JSON schemas.
	Schemas []map[string]string

	// Declarations of (data) types for use within the API.
	Types map[string]Type `yaml:"types"`

	// Declarations of traits for use within the API.
	Traits []map[string]Trait `yaml:"traits"`

	// Declarations of resource types for use within the API.
	ResourceTypes []map[string]ResourceType `yaml:"resourceTypes"`

	// TODO : annontation types

	// Declarations of security schemes for use within the API.
	SecuritySchemes []map[string]SecurityScheme `yaml:"securitySchemes"`

	// The security schemes that apply to every resource and method in the API.
	SecuredBy []DefinitionChoice `yaml:"securedBy"`

	// Imported external libraries for use within the API.
	Uses string `yaml:"uses"`

	// The resources of the API, identified as relative URIs that begin with a slash (/).
	// A resource property is one that begins with the slash and is either
	// at the root of the API definition or a child of a resource property. For example, /users and /{groupId}.
	Resources map[string]Resource `yaml:",regexp:/.*"`
}

func (apiDef *APIDefinition) postProcess() {
	// traits
	for i, tMap := range apiDef.Traits {
		for name := range tMap {
			t := tMap[name]
			t.postProcess(name)
			tMap[name] = t
			traitsMap[name] = t // add to global traits map
		}
		apiDef.Traits[i] = tMap
	}

	// resource types
	for i, rtMap := range apiDef.ResourceTypes {
		for name := range rtMap {
			rt := rtMap[name]
			rt.postProcess(name)
			rtMap[name] = rt
		}
		apiDef.ResourceTypes[i] = rtMap
	}

	// resources
	for k := range apiDef.Resources {
		r := apiDef.Resources[k]
		r.postProcess(k, nil, apiDef.ResourceTypes)
		apiDef.Resources[k] = r
	}

}

package raml

// Library is used to combine any collection of data type declarations,
// resource type declarations, trait declarations, and security scheme declarations
// into modular, externalized, reusable groups.
// While libraries are intended to define common declarations in external documents,
// which are then included where needed, libraries can also be defined inline.
type Library struct {
	Types           map[string]Type `yaml:"types"`
	Schemas         []map[string]string
	ResourceTypes   []map[string]ResourceType   `yaml:"resourceTypes"`
	Traits          []map[string]Trait          `yaml:"traits"`
	SecuritySchemes []map[string]SecurityScheme `yaml:"securitySchemes"`
	Uses            string
	Usage           string
}

func (l *Library) postProcess() error {
	// traits
	for i, tMap := range l.Traits {
		for name := range tMap {
			t := tMap[name]
			t.postProcess(name)
			tMap[name] = t
		}
		l.Traits[i] = tMap
	}

	// resource types
	for i, rtMap := range l.ResourceTypes {
		for name := range rtMap {
			rt := rtMap[name]
			rt.postProcess(name)
			rtMap[name] = rt
		}
		l.ResourceTypes[i] = rtMap
	}
	return nil
}

// get type by the name
func (l *Library) getType(name string) *Type {
	for k, t := range l.Types {
		if k == name {
			return &t
		}
	}
	return nil
}

// get schema by the name
func (l *Library) getSchema(name string) *string {
	for _, schemas := range l.Schemas {
		for k, schema := range schemas {
			if k == name {
				return &schema
			}
		}
	}
	return nil
}

// get resource type by the name
func (l *Library) getResourceType(name string) *ResourceType {
	for _, rts := range l.ResourceTypes {
		for k, rt := range rts {
			if k == name {
				return &rt
			}
		}
	}
	return nil
}

// get trait by the name
func (l *Library) getTrait(name string) *Trait {
	for _, ts := range l.Traits {
		for k, t := range ts {
			if k == name {
				return &t
			}
		}
	}
	return nil
}

// get security scheme by the name
func (l *Library) getSecurityScheme(name string) *SecurityScheme {
	for _, schemes := range l.SecuritySchemes {
		for k, s := range schemes {
			if k == name {
				return &s
			}
		}
	}
	return nil
}

package codegen

import (
	"strings"
	"time"
)

//Date represent RFC3399 date
type Date time.Time

//MarshalJSON override marshalJSON
func (t *Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*t).Format(`"` + time.RFC3339 + `"`)), nil
}

//MarshalText override marshalText
func (t *Date) MarshalText() ([]byte, error) {
	return []byte(time.Time(*t).Format(`"` + time.RFC3339 + `"`)), nil
}

//UnmarshalJSON override unmarshalJSON
func (t *Date) UnmarshalJSON(b []byte) error {
	ts, err := time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err != nil {
		return err
	}

	*t = Date(ts)
	return nil
}

//UnmarshalText override unmarshalText
func (t *Date) UnmarshalText(b []byte) error {
	ts, err := time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err != nil {
		return err
	}

	*t = Date(ts)
	return nil
}

func (t *Date) String() string {
	return time.Time(*t).String()
}

func convertUnion(strType string) string {
	if strings.Index(strType, "[]") > 0 {
		return "[]interface{}"
	}
	return "interface{}"
}

//ConvertToGoType handle convert from raml to go type
func convertToGoType(source string) string {
	switch source {
	case "string", "file":
		return "string"
	case "number":
		return "float64"
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "date":
		return "Date"
	}

	// other types that need some processing
	switch {
	case strings.HasSuffix(source, "[][]"): // bidimensional array
		return "[][]" + convertToGoType(source[:len(source)-4])
	case strings.HasSuffix(source, "[]"): // array
		return "[]" + convertToGoType(source[:len(source)-2])
	case strings.Index(source, "|") > 0:
		return convertUnion(source)
	}
	return source
}

package codegen

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/templates"
	"github.com/Jumpscale/go-raml/raml"
)

var (
	regex          = regexp.MustCompile("({{1}[\\w\\s]+}{1})")
	regNonAlphanum = regexp.MustCompile("[^A-Za-z0-9]+")
)

// doNormalizeURI removes `{`, `}`, and `/` from an URI
func doNormalizeURI(URI string) string {
	s := strings.Replace(URI, "/", " ", -1)
	s = strings.Replace(s, "{", "", -1)
	return strings.Replace(s, "}", "", -1)
}

// normalizeURI removes `{`, `}`, `/`, and space from an URI
func normalizeURI(URI string) string {
	return strings.Replace(doNormalizeURI(URI), " ", "", -1)
}

func normalizeURITitle(URI string) string {
	s := strings.Title(doNormalizeURI(URI))
	return strings.Replace(s, " ", "", -1)

}

// _getResourceParams is the recursive function of getResourceParams
func _getResourceParams(r *raml.Resource, params []string) []string {
	if r == nil {
		return params
	}

	matches := regex.FindAllString(r.URI, -1)
	for _, v := range matches {
		params = append(params, v[1:len(v)-1])
	}

	return _getResourceParams(r.Parent, params)
}

// get all params of a resource
// examples:
// /users  							  : no params
// /users/{userId}					  : params 1 = userId
// /users/{userId}/address/{addressId : params 1= userId, param 2= addressId
func getResourceParams(r *raml.Resource) []string {
	params := []string{}
	return _getResourceParams(r, params)
}

// create parameterized URI
// Input : raw string, ex : /users/{userId}/address/{addressId}
// Output : "/users/"+userId+"/address/"+addressId
func paramizingURI(URI string) string {
	uri := `"` + URI + `"`
	// replace { with "+
	uri = strings.Replace(uri, "{", `"+`, -1)

	// if ended with }/" or }", remove trailing "
	if strings.HasSuffix(uri, `}/"`) || strings.HasSuffix(uri, `}"`) {
		uri = uri[:len(uri)-1]
	}

	// replace } with +"
	uri = strings.Replace(uri, "}", `+"`, -1)

	// clean trailing +"
	if strings.HasSuffix(uri, `+"`) {
		uri = uri[:len(uri)-2]
	}
	return uri
}

// generate Go file from a template.
// if file already exist and override==false, file won't be regenerated
// funcMap = this parameter is used for passing go function to the template
func generateFile(data interface{}, tmplFile, tmplName, filename string, override bool) error {
	if !override && isFileExist(filename) {
		log.Infof("file %v already exist and override=false, no need to regenerate", filename)
		return nil
	}

	// pass Go function to template
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	// all template files path is relative to current directory (./)
	// while go-bindata files exist in templates directory
	tmplFile = strings.Replace(tmplFile, "./", "", -1)

	byteData, err := templates.Asset(tmplFile)
	if err != nil {
		return err
	}

	t, err := template.New(tmplName).Funcs(funcMap).Parse(string(byteData))
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Infof("generating file %v", filename)
	if err := t.ExecuteTemplate(f, tmplName, data); err != nil {
		return err
	}

	if strings.HasSuffix(filename, ".go") {
		return runGoFmt(filename)
	}
	return nil
}

// create directory if not exist
func checkCreateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}
	return nil
}

// cek if a file exist
func isFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// run `go fmt` command to a file
func runGoFmt(filePath string) error {
	args := []string{"fmt", filePath}

	if out, err := exec.Command("go", args...).CombinedOutput(); err != nil {
		log.Errorf("Error running go fmt on '%s' failed:\n%s", filePath, string(out))
		return errors.New("go fmt failed")
	}
	return nil
}

// convert interface type to string
// example :
// 1. string type, result would be string
// 2. []interface{} type, result would be array of string. ex: a,b,c
// Please add other type as needed
func interfaceToString(data interface{}) string {
	switch data.(type) {
	case string:
		return data.(string)
	case []interface{}:
		interfaceArr := data.([]interface{})
		resultStr := ""
		for _, v := range interfaceArr {
			resultStr += interfaceToString(v) + ","
		}
		return resultStr[:len(resultStr)-1]
	default:
		return ""
	}
}

// create string slice from an RAML description.
// each element is a  description line
func commentBuilder(desc string) []string {
	// we need to trim it because our parser usually give
	// space after last newline
	desc = strings.TrimSpace(desc)

	if desc == "" {
		return []string{}
	}

	return strings.Split(desc, "\n")
}

// replace non alphanumerics with "_"
func replaceNonAlphanumerics(s string) string {
	return strings.Trim(regNonAlphanum.ReplaceAllString(s, "_"), "_")
}

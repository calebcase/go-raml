package codegen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jumpscale/go-raml/raml"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateStructFromRaml(t *testing.T) {
	Convey("generate struct from raml", t, func() {
		apiDef, err := raml.ParseFile("./fixtures/struct/struct.raml")
		So(err, ShouldBeNil)
		targetdir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)

		Convey("Simple struct from raml", func() {
			err = generateStructs(apiDef, targetdir, "main", langGo)
			So(err, ShouldBeNil)

			//first test
			s, err := testLoadFile(filepath.Join(targetdir, "EnumCity.go"))
			So(err, ShouldBeNil)

			tmpl, err := testLoadFile("./fixtures/struct/enumcity.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			//second test
			s, err = testLoadFile(filepath.Join(targetdir, "animal.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/animal.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			//third test, single inheritance
			s, err = testLoadFile(filepath.Join(targetdir, "SingleInheritance.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/singleinheritance.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			//fourth test, multiple inheritance
			s, err = testLoadFile(filepath.Join(targetdir, "MultipleInheritance.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/multipleinheritance.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			//fifth test, array of object
			s, err = testLoadFile(filepath.Join(targetdir, "ArrayOfCats.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/arrayofcats.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// bidimensional array
			s, err = testLoadFile(filepath.Join(targetdir, "BidimensionalArrayOfCats.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/bidimensionalarrayofcats.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// using map type & testing case sensitive type name
			s, err = testLoadFile(filepath.Join(targetdir, "petshop.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/petshop.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// Union
			s, err = testLoadFile(filepath.Join(targetdir, "Pet.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/Pet.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// Array of union
			s, err = testLoadFile(filepath.Join(targetdir, "ArrayOfPets.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/ArrayOfPets.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// Specialization
			s, err = testLoadFile(filepath.Join(targetdir, "Specialization.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/Specialization.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// Enum type
			s, err = testLoadFile(filepath.Join(targetdir, "EnumString.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/enumstring.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

			// With validation
			s, err = testLoadFile(filepath.Join(targetdir, "ValidationString.go"))
			So(err, ShouldBeNil)

			tmpl, err = testLoadFile("./fixtures/struct/ValidationString.txt")
			So(err, ShouldBeNil)

			So(s, ShouldEqual, tmpl)

		})

		Reset(func() {
			os.RemoveAll(targetdir)
		})
	})
}

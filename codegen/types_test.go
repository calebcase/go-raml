package codegen

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeConversion(t *testing.T) {
	Convey("Test Type Conversion", t, func() {
		Convey("Type conversion", func() {
			So(convertToGoType("string"), ShouldEqual, "string")
			So(convertToGoType("number"), ShouldEqual, "float64")
			So(convertToGoType("integer"), ShouldEqual, "int")
			So(convertToGoType("boolean"), ShouldEqual, "bool")
			So(convertToGoType("date"), ShouldEqual, "Date")
			So(convertToGoType("file"), ShouldEqual, "string")
			So(convertToGoType("Object"), ShouldEqual, "Object")
			So(convertToGoType("string[]"), ShouldEqual, "[]string")
			So(convertToGoType("string[][]"), ShouldEqual, "[][]string")
			So(convertToGoType("string | Person"), ShouldEqual, "interface{}")
			So(convertToGoType("(string | Person)[]"), ShouldEqual, "[]interface{}")
		})
	})
}

func TestMarshalUnmarshal(t *testing.T) {
	Convey("Test Marshal Unmarshal", t, func() {
		Convey("Test Marshal Unmarshal", func() {
			dateNow := Date(time.Now())

			j, err := dateNow.MarshalJSON()
			So(err, ShouldBeNil)
			So(len(string(j)), ShouldNotEqual, 0)

			t, err := dateNow.MarshalText()
			So(err, ShouldBeNil)
			So(len(string(t)), ShouldNotEqual, 0)

			err = dateNow.UnmarshalJSON(j)
			So(err, ShouldBeNil)

			err = dateNow.UnmarshalText(t)
			So(err, ShouldBeNil)
		})
	})
}

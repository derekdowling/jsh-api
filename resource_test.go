package jshapi

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResource(t *testing.T) {

	Convey("Resource Tests", t, func() {

		resource := testResource()
		server := httptest.NewServer(resource)
		baseURL := server.URL

		Convey("->Post()", func() {

			object := testObject()
			object.ID = "2"

			resp, err := jsc.Post(baseURL, object)
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			So(err, ShouldBeNil)

			So(obj.ID, ShouldEqual, "1")
		})

		Convey("->Get()", func() {
			resp, err := jsc.Get(baseURL, "test", "")
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			log.Printf("err.String() = %+v\n", err.String())
			So(err, ShouldBeNil)

			testObj := testObject()
			testObj.ID = "1"
			So(obj, ShouldResemble, testObj)
		})

		Convey("->Get(id)", func() {
			resp, err := jsc.Get(baseURL, "test", "3")
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			So(err, ShouldBeNil)

			testObj := testObject()
			testObj.ID = "3"
			So(obj, ShouldResemble, testObj)
		})
	})
}

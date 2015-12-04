package jshapi

import (
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResource(t *testing.T) {

	Convey("Resource Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resource := NewMockResource("", "test", 2, attrs)
		server := httptest.NewServer(resource)
		baseURL := server.URL

		Convey("->Matcher()", func() {
			resource.prefix = "api"
			So(resource.Matcher(), ShouldEqual, "/api/tests")
		})

		Convey("->IDMatcher()", func() {
			resource.prefix = "api"
			So(resource.IDMatcher(), ShouldEqual, "/api/tests/:id")
		})

		Convey("->Post()", func() {

			object := sampleObject("", "test", attrs)
			resp, err := jsc.Post(baseURL, object)
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			So(err, ShouldBeNil)

			So(obj.ID, ShouldEqual, "1")
		})

		Convey("->List()", func() {
			resp, err := jsc.Get(baseURL, "test", "")
			So(err, ShouldBeNil)

			list, err := resp.GetList()
			So(err, ShouldBeNil)

			So(len(list), ShouldEqual, 2)
			So(list[0].ID, ShouldEqual, "1")
		})

		Convey("->Get()", func() {
			resp, err := jsc.Get(baseURL, "test", "3")
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			So(err, ShouldBeNil)
			So(obj.ID, ShouldEqual, "3")
		})
	})
}

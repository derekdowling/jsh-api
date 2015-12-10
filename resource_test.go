package jshapi

import (
	"net/http"
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

		resourceType := "foo"
		resource := NewMockResource("", resourceType, 2, attrs)
		server := httptest.NewServer(resource)
		baseURL := server.URL

		Convey("->Matcher()", func() {
			resource.prefix = "/api"
			So(resource.Matcher(), ShouldEqual, "/api/"+resourceType+"s")
		})

		Convey("->IDMatcher()", func() {
			resource.prefix = "/api"
			So(resource.IDMatcher(), ShouldEqual, "/api/"+resourceType+"s/:id")
		})

		Convey("->Post()", func() {
			object := sampleObject("", resourceType, attrs)
			resp, err := jsc.Post(baseURL, object)
			So(resp.StatusCode, ShouldEqual, http.StatusCreated)
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			So(err, ShouldBeNil)

			So(obj.ID, ShouldEqual, "1")
		})

		Convey("->List()", func() {
			resp, err := jsc.Get(baseURL, resourceType, "")
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)

			list, err := resp.GetList()
			So(err, ShouldBeNil)

			So(len(list), ShouldEqual, 2)
			So(list[0].ID, ShouldEqual, "1")
		})

		Convey("->Get()", func() {
			resp, err := jsc.Get(baseURL, resourceType, "3")
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			So(err, ShouldBeNil)
			So(obj.ID, ShouldEqual, "3")
		})

		Convey("->Patch()", func() {
			object := sampleObject("1", resourceType, attrs)
			resp, err := jsc.Patch(baseURL, object)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)

			obj, err := resp.GetObject()
			So(err, ShouldBeNil)
			So(obj.ID, ShouldEqual, "1")
		})

		Convey("->Delete()", func() {
			resp, err := jsc.Delete(baseURL, resourceType, "1")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
		})

	})
}

func TestSubResource(t *testing.T) {

	Convey("Sub-Resource Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "bar"
		resource := NewMockResource("/foo", resourceType, 2, attrs)

		subResourceType := "baz"
		subStorageMock := &MockStorage{
			ResourceType:       subResourceType,
			ResourceAttributes: attrs,
			ListCount:          2,
		}

		subResource := resource.CreateSubResource(subResourceType)
		subResource.CRUD(subStorageMock)

		// server := httptest.NewServer(resource)
		// baseURL := server.URL + subResource.prefix

		Convey("subResource prefix", func() {
			So(subResource.prefix, ShouldEqual, "/foo/bars/:id")
		})

		Convey("->Matcher()", func() {
			So(subResource.Matcher(), ShouldEqual, "/foo/bars/:id/bazs")
		})

		Convey("->IDMatcher()", func() {
			So(subResource.IDMatcher(), ShouldEqual, "/foo/bars/:id/bazs/:id")
		})

		Convey("->List()", func() {
			// resp, err := jsc.Get(baseURL, subResourceType, "")
			// log.Printf("resp.Response.Request = %+v\n", resp.Response.Request)
			// So(resp.StatusCode, ShouldEqual, http.StatusOK)
			// So(err, ShouldBeNil)

			// list, err := resp.GetList()
			// So(err, ShouldBeNil)

			// So(len(list), ShouldEqual, 2)
			// So(list[0].ID, ShouldEqual, "1")

			// So(1, ShouldBeFalse)
		})

	})
}

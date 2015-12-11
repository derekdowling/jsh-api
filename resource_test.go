package jshapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResource(t *testing.T) {

	Convey("Resource Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "foo"
		resource := NewMockResource("api", resourceType, 2, attrs)
		server := httptest.NewServer(resource)
		baseURL := server.URL

		So(len(resource.Routes), ShouldEqual, 5)

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

func TestCustomAction(t *testing.T) {

	Convey("Resource Action Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}
		resourceType := "bar"
		resource := NewMockResource("/foo", resourceType, 2, attrs)

		Convey("->NewAction()", func() {
			resource.NewAction("mutate", func(ctx context.Context, object *jsh.Object) jsh.SendableError {
				target := map[string]string{}
				object.Unmarshal("bar", target)
				target["mutated"] = "true"
				return nil
			})

			So(len(resource.Routes), ShouldEqual, 6)
			So(resource.Routes[len(resource.Routes)-1], ShouldEqual, "PATCH - /foo/bars/:id/mutate")
		})
	})
}

func TestNestedResource(t *testing.T) {

	Convey("Nested Resource Tests", t, func() {

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

		subResource := resource.NewNestedResource(subResourceType)
		subResource.CRUD(subStorageMock)

		// server := httptest.NewServer(resource)
		// baseURL := server.URL + subResource.prefix

		Convey("Resource", func() {

			Convey("should track sub-resources properly", func() {
				So(len(resource.Subresources), ShouldEqual, 1)
				So(resource.Subresources[subResourceType], ShouldEqual, subResource)
			})
		})

		Convey("Sub-Resource prefix", func() {
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
			// So(resp.StatusCode, ShouldEqual, http.StatusOK)
			// So(err, ShouldBeNil)

			// list, err := resp.GetList()
			// So(err, ShouldBeNil)

			// So(len(list), ShouldEqual, 2)
			// So(list[0].ID, ShouldEqual, "1")
		})

	})
}

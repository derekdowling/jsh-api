package jshapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI(t *testing.T) {

	Convey("API Tests", t, func() {

		resourceType := "foos"
		api := New("api")

		So(api.prefix, ShouldEqual, "/api")

		testAttrs := map[string]string{
			"foo": "bar",
		}

		Convey("->AddResource()", func() {
			resource := NewMockResource(resourceType, 1, testAttrs)
			api.Add(resource)

			So(api.Resources[resourceType], ShouldEqual, resource)

			server := httptest.NewServer(api)
			baseURL := server.URL + api.prefix

			_, resp, err := jsc.Fetch(baseURL, resourceType, "1")

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)

			patchObj, err := jsh.NewObject("1", resourceType, testAttrs)

			_, resp, patchErr := jsc.Patch(baseURL, patchObj)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(patchErr, ShouldBeNil)
		})
	})
}

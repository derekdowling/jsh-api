package jshapi

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"goji.io"
	"goji.io/pat"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func TestResource(t *testing.T) {

	Convey("Resource Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "foos"
		resource := NewMockResource(resourceType, 2, attrs)

		mux := goji.NewMux()
		mux.HandleC(pat.New("/foos"), resource.Mux)
		mux.HandleC(pat.New("/foos/*"), resource.Mux)

		server := httptest.NewServer(mux)
		baseURL := server.URL

		So(len(resource.Routes), ShouldEqual, 5)

		Convey("->NewResource()", func() {

			Convey("should be agnostic to plurality", func() {
				resource := NewResource("users")
				So(resource.Type, ShouldEqual, "users")

				resource2 := NewResource("user")
				So(resource2.Type, ShouldEqual, "user")
			})
		})

		Convey("->Post()", func() {
			object := sampleObject("", resourceType, attrs)
			doc, resp, err := jsc.Post(baseURL, object)

			So(resp.StatusCode, ShouldEqual, http.StatusCreated)
			So(err, ShouldBeNil)
			So(doc.Data[0].ID, ShouldEqual, "1")
		})

		Convey("->List()", func() {
			doc, resp, err := jsc.List(baseURL, resourceType)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(len(doc.Data), ShouldEqual, 2)
			So(doc.Data[0].ID, ShouldEqual, "1")
		})

		Convey("->Fetch()", func() {
			doc, resp, err := jsc.Fetch(baseURL, resourceType, "3")

			log.Printf("resp = %+v\n", resp.Request)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(doc.Data[0].ID, ShouldEqual, "3")
		})

		Convey("->Patch()", func() {
			object := sampleObject("1", resourceType, attrs)
			doc, resp, err := jsc.Patch(baseURL, object)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(doc.Data[0].ID, ShouldEqual, "1")
		})

		Convey("->Delete()", func() {
			resp, err := jsc.Delete(baseURL, resourceType, "1")

			So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
			So(err, ShouldBeNil)
		})
	})
}

func TestMutateHandler(t *testing.T) {

	Convey("Custom Handler Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "bars"
		resource := NewMockResource(resourceType, 2, attrs)

		handler := func(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType) {
			object := sampleObject(id, resourceType, attrs)
			return object, nil
		}

		resource.Mutate("mutate", handler)

		server := httptest.NewServer(resource)
		baseURL := server.URL + patID

		Convey("Resource State", func() {
			So(len(resource.Routes), ShouldEqual, 6)
			So(resource.Routes[len(resource.Routes)-1], ShouldEqual, "PATCH - /bars/:id/mutate")
		})

		Convey("->Custom()", func() {
			doc, response, err := jsc.Get(baseURL + "/mutate")

			So(err, ShouldBeNil)
			So(response.StatusCode, ShouldEqual, http.StatusOK)
			So(doc.Data, ShouldNotBeEmpty)
		})
	})
}

func TestToOne(t *testing.T) {

	Convey("Relationship ToOne Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "bar"
		resource := NewMockResource(resourceType, 2, attrs)

		relationshipHandler := func(ctx context.Context, resourceID string) (*jsh.Object, jsh.ErrorType) {
			return sampleObject("1", "baz", map[string]string{"baz": "ball"}), nil
		}

		subResourceType := "baz"
		resource.ToOne(subResourceType, relationshipHandler)

		server := httptest.NewServer(resource)
		baseURL := server.URL + patID

		Convey("Resource State", func() {

			Convey("should track sub-resources properly", func() {
				So(len(resource.Relationships), ShouldEqual, 1)
				So(len(resource.Routes), ShouldEqual, 7)
			})
		})

		Convey("->ToOne()", func() {

			Convey("/foo/bars/:id/baz", func() {
				doc, resp, err := jsc.Get(baseURL + "/" + subResourceType)

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})

			Convey("/foo/bars/:id/relationships/baz", func() {
				doc, resp, err := jsc.Get(baseURL + "/relationships/" + subResourceType)

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})
		})
	})
}

func TestToMany(t *testing.T) {

	Convey("Relationship ToMany Tests", t, func() {

		attrs := map[string]string{
			"foo": "bar",
		}

		resourceType := "bar"
		resource := NewMockResource(resourceType, 2, attrs)

		relationshipHandler := func(ctx context.Context, resourceID string) (jsh.List, jsh.ErrorType) {
			return jsh.List{
				sampleObject("1", "baz", map[string]string{"baz": "ball"}),
				sampleObject("2", "baz", map[string]string{"baz": "ball2"}),
			}, nil
		}

		subResourceType := "baz"
		resource.ToMany(subResourceType, relationshipHandler)

		server := httptest.NewServer(resource)
		baseURL := server.URL + patID

		Convey("Resource State", func() {

			Convey("should track sub-resources properly", func() {
				So(len(resource.Relationships), ShouldEqual, 1)
				So(len(resource.Routes), ShouldEqual, 7)
			})
		})

		Convey("->ToOne()", func() {

			Convey("/foo/bars/:id/bazs", func() {
				doc, resp, err := jsc.Get(baseURL + "/" + subResourceType + "s")

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(len(doc.Data), ShouldEqual, 2)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})

			Convey("/foo/bars/:id/relationships/bazs", func() {
				doc, resp, err := jsc.Get(baseURL + "/relationships/" + subResourceType + "s")

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)

				So(err, ShouldBeNil)
				So(len(doc.Data), ShouldEqual, 2)
				So(doc.Data[0].ID, ShouldEqual, "1")
			})
		})

	})
}

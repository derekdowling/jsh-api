package jshapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI(t *testing.T) {

	Convey("API Tests", t, func() {

		api := New("")

		Convey("->AddResource()", func() {
			resource := testResource()
			api.AddResource(resource)

			So(api.Resources["test"], ShouldEqual, resource)
		})
	})
}

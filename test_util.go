package jshapi

import (
	"log"

	"github.com/derekdowling/go-json-spec-handler"
)

const testType = "test"

type mockStorage struct{}

func (m *mockStorage) Save(object *jsh.Object) jsh.SendableError {
	object.ID = "1"
	return nil
}

func (m *mockStorage) Get(id string) (*jsh.Object, jsh.SendableError) {
	obj := testObject()
	obj.ID = id
	return obj, nil
}

func (m *mockStorage) List() (jsh.List, jsh.SendableError) {
	return testList(), nil
}

func (m *mockStorage) Patch(object *jsh.Object) jsh.SendableError {
	return nil
}

func (m *mockStorage) Delete(id string) jsh.SendableError {
	return nil
}

func testResource() *Resource {
	mock := &mockStorage{}
	return NewResource("", "test", mock)
}

func testObject() *jsh.Object {
	object, err := jsh.NewObject("1", testType, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	return object
}

func testList() jsh.List {
	test1 := testObject()
	test1.ID = "1"

	test2 := testObject()
	test2.ID = "2"
	return jsh.List{test1, test2}
}

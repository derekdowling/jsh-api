package store

// generator implements a gen(https://github.com/clipperhouse/gen) TypeWriter for
// stubbing a store.CRUD implementation and corresponding tests.

import (
	"io"
	"log"

	"github.com/clipperhouse/typewriter"
)

func init() {
	err := typewriter.Register(NewCRUDWriter())
	if err != nil {
		panic(err)
	}
}

var templates = typewriter.TemplateSlice{
	crud,
	crudTest,
}

// CRUDWriter implements gen/TypeWriter
type CRUDWriter struct{}

// NewCRUDWriter is a simple constructor
func NewCRUDWriter() *CRUDWriter {
	return &CRUDWriter{}
}

// Name sets the gen tag for the CRUDWriter
func (cw *CRUDWriter) Name() string {
	return "CRUDWriter"
}

// Imports represent the dependencies of the generated code.
func (cw *CRUDWriter) Imports(t typewriter.Type) (result []typewriter.ImportSpec) {

	imports := []typewriter.ImportSpec{
		typewriter.ImportSpec{
			Path: "github.com/derekdowling/go-json-spec-handler",
		},
		typewriter.ImportSpec{
			Path: "github.com/derekdowling/jsh-api/store",
		},
	}

	// none
	return result
}

func (cw *CRUDWriter) Write(w io.Writer, t typewriter.Type) error {
	tag, found := t.FindTag(sw)
	if !found {
		log.Printf("Store type '%s' not found\n", t.Name)
		return nil
	}

	license := `// This file implement store.CRUD https://github.com/derekdowling/jsh-api/store`
	if _, err := w.Write([]byte(license)); err != nil {
		return err
	}

	tmpl, err := templates.ByTag(t, tag)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, t); err != nil {
		return err
	}

	return nil
}

package set

import "github.com/clipperhouse/typewriter"

var crud = &typewriter.Template{
	Name: "CRUD",
	Text: `
// {{.Name}}Store holds state for {{.Name}} Storage actions
type {{.Name}}Store struct{}

// New{{.Name}}Set creates and returns a reference to an empty set.
func New{{.Name}}() {{.Pointer}}{{.Name}}Store {
    return &{{.Name}}Store{}
}

// ToSlice returns the elements of the current set as a slice
func (set {{.Name}}Set) ToSlice() []{{.Pointer}}{{.Name}} {
    var s []{{.Pointer}}{{.Name}}
    for v := range set {
        s = append(s, v)
    }
    return s
}

// Save handles POST /{{.Name}}
func (store {{.Pointer}}{{.Name}}Store) Save(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
    {{.Name}} := &{{.Name}}{}
    unmarshalErr := object.Unmarshal("{{.Name}}", {{.Name}})
    if err != nil {
	return nil, unmarshalErr
    }

    // TODO: Save {{.Name}}

    marshalErr := object.Marshal({{.Name}})
    if marshalErr != nil {
	return nil, marshalErr
    }

    return object, nil
}

// Save handles GET /{{.Name}}/:id
func (store {{.Pointer}}{{.Name}}Store) Get(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType) {
    var {{.Name}} *{{.Name}}

    // TODO: Get {{.Name}}

    object, err := jsh.NewObject(id, "{{.Name}}", {{.Name}})
    if err != nil {
	return nil, err
    }

    return object, nil
}

// Save handles GET /{{.Name}}
func (store {{.Pointer}}{{.Name}}Store) List(ctx context.Context) (jsh.List, jsh.ErrorType) {
    list := &jsh.List{}

    // TODO: Get {{.Name}} list

    return list, nil
}

// Save handles PATCH /{{.Name}}/:id
func (store {{.Pointer}}{{.Name}}Store) Update(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
    {{.Name}} := &{{.Name}}{}
    unmarshalErr := object.Unmarshal("{{.Name}}", {{.Name}})
    if unmarshalErr != nil {
	return nil, unmarshalErr
    }

    id := object.ID

    // TODO: Update {{.Name}} object

    marshalErr := object.Marshal({{.Name}})
    if marshalErr != nil {
	return nil, marshalErr
    }

    return object, nil
}

// Save handles DELETE /{{.Name}}/:id
func (store {{.Pointer}}{{.Name}}Store) Delete(ctx context.Context, id string) jsh.ErrorType {
    // TODO: implement Delete logic
    return nil
}
`,
	TypeConstraint: typewriter.Constraint{},
}

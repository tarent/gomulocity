package generic

import (
	"testing"
)

func TestJsonc_ErrorOnNoStruct(t *testing.T) {
	j, err := JsonFromObject("Hallo Welt")

	if err == nil {
		t.Errorf("JsonFromObject - error expected. Instead: %v", j)
	}
}

func TestJsonc_ErrorOnPlainStruct(t *testing.T) {
	type A struct {
		B string
		C int
	}

	a := A{B: "Foo", C: 4711}

	j, err := JsonFromObject(a)

	if err == nil {
		t.Errorf("JsonFromObject - error expected. Instead: %v", j)
	}
}

func TestJsonc_MarshalStandardFields(t *testing.T) {
	type A struct {
		B string
		C int
	}

	a := &A{B: "Foo", C: 4711}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"B":"Foo","C":4711}`
	if string(j) != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_SupportsStandardJsonTags(t *testing.T) {
	type A struct {
		B string `json:"myB"`
		C int    `json:"myC"`
		D string `json:"-"`
		E string `json:"myE,omitempty"`
		F bool   `json:"myF,omitempty"`
		G int    `json:",omitempty"`
		H *A     `json:"myH,omitempty"`
		I string `json:"myI,otherstuff"`
	}

	a := &A{B: "Foo", C: 4711, D: "Bar", E: "", F: false, G: 0, H: nil, I: ""}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"myB":"Foo","myC":4711,"myI":""}`
	if string(j) != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_FlatsTaggedFields(t *testing.T) {
	type A struct {
		B string            `json:"myB"`
		C int               `json:"myC"`
		D map[string]string `jsonc:"flat"`
	}

	a := &A{B: "Foo", C: 4711, D: map[string]string{"Da": "foo", "Db": "bar", "Dc": "baz"}}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"Da":"foo","Db":"bar","Dc":"baz","myB":"Foo","myC":4711}`
	if string(j) != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_FlatStructs(t *testing.T) {
	type B struct {
		Field1 string
		Field2 string
	}

	type A struct {
		Bs map[string]B `jsonc:"flat"`
	}

	a := &A{Bs: map[string]B{
		"foo": {Field1: "#Field2", Field2: "#Field1"},
	}}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"foo":{"Field1":"#Field2","Field2":"#Field1"}}`
	if string(j) != want {
		t.Errorf("JsonFromObject\n json = %v\n want %v", string(j), want)
	}
}

func TestJsonc_FlatsJsonTaggedStructs(t *testing.T) {
	type Sub struct {
		Custom1 string `json:"one"`
		Custom2 string `json:"two"`
		Custom3 int    `json:"-"`
		Custom4 bool   `json:"four,omitempty"`
		Custom5 string `json:",omitempty"`
	}

	type A struct {
		B string                 `json:"myB"`
		C int                    `json:"myC"`
		D map[string]interface{} `jsonc:"flat"`
	}

	m := map[string]interface{}{
		"foo": Sub{
			Custom1: "myCustom1",
			Custom2: "myCustom2",
			Custom3: 4711,
			Custom4: false,
			Custom5: "",
		},
	}

	a := &A{B: "Foo", C: 4711, D: m}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"foo":{"one":"myCustom1","two":"myCustom2"},"myB":"Foo","myC":4711}`
	if string(j) != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_DoesNotFlatUntaggedMaps(t *testing.T) {
	type A struct {
		B string
		C map[string]string
	}

	a := &A{B: "Foo", C: map[string]string{"foo1": "bar", "foo2": "baz"}}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"B":"Foo","C":{"foo1":"bar","foo2":"baz"}}`
	if string(j) != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_JsonFromObject_WrongTags(t *testing.T) {
	type WrongFlat struct {
		A string `jsonc:"flat"`
	}
	type WrongCollection struct {
		A string `jsonc:"collection"`
	}

	_, err := JsonFromObject(&WrongFlat{A: "Hello"})
	if err == nil {
		t.Errorf("JsonFromObject - no error, want error for wrong use of jsonc:flat.")
	}

	_, err = JsonFromObject(&WrongCollection{A: "Hello"})
	if err == nil {
		t.Errorf("JsonFromObject - no error, want error for wrong use of jsonc:collection.")
	}
}

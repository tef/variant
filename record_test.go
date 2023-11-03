package variant

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Header struct {
	Kind string
}

type RecordOne struct {
	Header
	A int
}

type RecordTwo struct {
	Header
	B string
}

func recordNew(kind string) any {
	switch kind {
	case "one":
		return &RecordOne{}
	case "two":
		return &RecordTwo{}
	}
	return nil
}

func recordKind(obj any) string {
	switch obj.(type) {
	case *RecordOne:
		return "one"
	case *RecordTwo:
		return "two"
	}
	return ""
}

func jsonKind(bytes []byte) (string, error) {
	header := Header{}

	err := json.Unmarshal(bytes, &header)
	if err != nil {
		return "", err
	}

	return header.Kind, nil
}

func newJson(kind string, obj any) ([]byte, error) {
	switch o := obj.(type) {
	case *RecordOne:
		o.Header.Kind = kind
	case *RecordTwo:
		o.Header.Kind = kind
	}
	return json.Marshal(obj)
}

func TestRecord(t *testing.T) {

	roundTrip := func(objIn *Record, objOut *Record) error {
		buf, err := json.Marshal(objIn)

		if err != nil {
			return err
		}

		return json.Unmarshal(buf, objOut)
	}

	RecordHandler.New = recordNew
	RecordHandler.Kind = recordKind
	RecordHandler.JsonKind = jsonKind
	RecordHandler.NewJson = newJson

	one := NewRecord(&RecordOne{A: 1})
	two := NewRecord(&RecordTwo{B: "2"})

	if recordKind(one.Value) != "one" {
		t.Error("kind of one is not one")
	}

	if recordKind(two.Value) != "two" {
		t.Error("kind of one is not one")
	}

	var three, four Record

	roundTrip(&one, &three)
	roundTrip(&two, &four)

	if !reflect.DeepEqual(one, three) {
		t.Error("mismatch", one.Value, three.Value)
	}

	if !reflect.DeepEqual(two, four) {
		t.Error("mismatch", two.Value, four.Value)
	}
}

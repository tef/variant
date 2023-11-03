package variant

import (
	"encoding/json"
	"reflect"
	"testing"
)

type One struct {
	A int
}

type Two struct {
	B string
}

func variantNew(kind string) any {
	switch kind {
	case "one":
		return &One{}
	case "two":
		return &Two{}
	}
	return nil
}

func variantKind(obj any) string {
	switch obj.(type) {
	case *One:
		return "one"
	case *Two:
		return "two"
	}
	return ""
}

func TestEntry(t *testing.T) {

	roundTrip := func(objIn *Entry, objOut *Entry) error {
		buf, err := json.Marshal(objIn)

		if err != nil {
			return err
		}

		return json.Unmarshal(buf, objOut)
	}

	EntryHandler.New = variantNew
	EntryHandler.Kind = variantKind

	one := NewEntry(&One{A: 1})
	two := NewEntry(&Two{B: "2"})

	if variantKind(one.Value) != "one" {
		t.Error("kind of one is not one")
	}

	if variantKind(two.Value) != "two" {
		t.Error("kind of one is not one")
	}

	var three, four Entry

	roundTrip(&one, &three)
	roundTrip(&two, &four)

	if !reflect.DeepEqual(one, three) {
		t.Error("mismatch", one.Value, three.Value)
	}

	if !reflect.DeepEqual(two, four) {
		t.Error("mismatch", two.Value, four.Value)
	}
}

func TestPair(t *testing.T) {

	roundTrip := func(objIn *Pair, objOut *Pair) error {
		buf, err := json.Marshal(objIn)

		if err != nil {
			return err
		}

		return json.Unmarshal(buf, objOut)
	}

	PairHandler.New = variantNew
	PairHandler.Kind = variantKind
	PairHandler.KindKey = "kind"
	PairHandler.ValueKey = "value"

	one := NewPair(&One{A: 1})
	two := NewPair(&Two{B: "2"})

	if variantKind(one.Value) != "one" {
		t.Error("kind of one is not one")
	}

	if variantKind(two.Value) != "two" {
		t.Error("kind of one is not one")
	}

	var three, four Pair

	roundTrip(&one, &three)
	roundTrip(&two, &four)

	if !reflect.DeepEqual(one, three) {
		t.Error("mismatch", one.Value, three.Value)
	}

	if !reflect.DeepEqual(two, four) {
		t.Error("mismatch", two.Value, four.Value)
	}
}

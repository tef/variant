package variant

import (
	"encoding/json"
)

var RecordHandler struct {
	JsonKind func([]byte) (string, error)      // the tag in this json
	NewJson  func(string, any) ([]byte, error) // json with this tag
	Kind     func(any) string                  // the tag for this struct
	New      func(string) any                  // new(*struct) for this tag
}

type Record struct {
	Kind  string
	Value any
}

func NewRecord(obj any) Record {
	return Record{
		Kind:  RecordHandler.Kind(obj),
		Value: obj,
	}
}

func (r *Record) UnmarshalJSON(bytes []byte) error {
	kind, err := RecordHandler.JsonKind(bytes)
	if err != nil {
		return err
	}

	r.Kind = kind
	r.Value = RecordHandler.New(kind)
	return json.Unmarshal(bytes, r.Value)
}

func (r *Record) MarshalJSON() ([]byte, error) {
	if RecordHandler.NewJson != nil {
		return RecordHandler.NewJson(r.Kind, r.Value)
	}

	// if none defined, just turn it into json
	// and assume it contains the tag
	return json.Marshal(r.Value)
}

func TagFromHeader(newHeader func() any, headerKind func(any) string) func([]byte) (string, error) {
	return func(bytes []byte) (string, error) {
		header := newHeader()
		// header is a ptr to a struct, so we
		// don't need &header

		err := json.Unmarshal(bytes, header)
		if err != nil {
			return "", err
		}

		return headerKind(header), nil
	}
}

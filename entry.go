package variant

import (
	"encoding/json"
)

var EntryHandler struct {
	Kind func(any) string // the tag for this struct
	New  func(string) any // new(*struct) for this tag
}

type Entry struct {
	// a {"key": value} json obj
	Kind  string
	Value any
}

func NewEntry(obj any) Entry {
	return Entry{
		Kind:  EntryHandler.Kind(obj),
		Value: obj,
	}
}

func (e *Entry) UnmarshalJSON(bytes []byte) error {
	var raw map[string]json.RawMessage

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	for k, v := range raw {
		e.Kind = k
		e.Value = EntryHandler.New(k)
		return json.Unmarshal(v, e.Value)
	}

	return nil
}

func (e *Entry) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(e.Value)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]json.RawMessage, 1)
	raw[e.Kind] = buf

	return json.Marshal(raw)
}

package variant

import (
	"encoding/json"
)

var PairHandler struct {
	KindKey  string           // the key name that holds the tag
	ValueKey string           // the key name that holds the value
	Kind     func(any) string // the tag for this struct
	New      func(string) any // new(*struct) for this tag
}

type Pair struct {
	// a {"kind": kind, "value":value} json obj
	Kind  string
	Value any
}

func NewPair(obj any) Pair {
	return Pair{
		Kind:  PairHandler.Kind(obj),
		Value: obj,
	}
}

func (p *Pair) UnmarshalJSON(bytes []byte) error {
	var raw map[string]json.RawMessage

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	kindRaw := raw[PairHandler.KindKey]
	valueRaw := raw[PairHandler.ValueKey]

	err = json.Unmarshal(kindRaw, &p.Kind)
	if err != nil {
		return nil
	}

	p.Value = PairHandler.New(p.Kind)
	return json.Unmarshal(valueRaw, p.Value)
}

func (p *Pair) MarshalJSON() ([]byte, error) {
	rawKind, err := json.Marshal(p.Kind)
	if err != nil {
		return nil, err
	}

	rawValue, err := json.Marshal(p.Value)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]json.RawMessage, 2)
	raw[PairHandler.KindKey] = rawKind
	raw[PairHandler.ValueKey] = rawValue

	return json.Marshal(raw)
}

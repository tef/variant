# variant: union types for json

when you have a union type in json, there aren't many good options for
handling it in go. really, all you can do is something like this:

```
type Variant struct {
    date *Date `json:",omitempty"`
    complex *Complex `json:",omitempty"`
}
```

if your json doesn't look like `{"date":...}` or `{"complex": ...}`, then
you're kinda out of luck. If you have over 20+ types, it gets a little messy.

this library handles three different styles of tagged json. the first is one
you've just seen.

# entry style, or externally tagged json

does your json look like this?

```
{"date": "1997-03-30"} 
{"complex: [1.0, -1.0]}
```

then your go could look like this:

```

variant.EntryHandler.New = func(kind string) any {
        switch kind {
        case "date":
                return &Date{}
        case "complex":
                return &Complex{}
        }
}

variant.EntryHandler.Kind = func(obj any) string {
        switch obj.(type) {
        case *Date:
                return "date"
        case *Complex:
                return "complex"
        }
```

... and you're done.

# pair style, or adjacently tagged json

ok, what if your json looks like this?

```
{"kind": "date", "value":"1997-03-30"} 
{"kind": "complex", "value": [1.0, -1.0]}
```

then you'll need to set `KindKey` and `ValueKey`, as well as `Kind` and `New`

```
variant.PairHandler.KindKey = "kind"
variant.PairHandler.ValueKey = "value"
variant.PairHandler.New = func(kind string) any {
        switch kind {
        case "date":
                return &Date{}
        case "complex":
                return &Complex{}
        }
        return nil
}

variant.PairHandler.Kind = func(obj any) string {
        switch obj.(type) {
        case *Date:
                return "date"
        case *Complex:
                return "complex"
        }
        return ""
}
```

# record style, or internally tagged json

Your json might look like this:

```
{"kind": "date", "year":1997, "month":03, "day":30} 
{"kind": "complex", "real": 1.0, "imaginary": -1.0}
```

Unfortunately, this is a lot more work.

```

type Header struct {
        Kind string `json:"kind"`
}

struct Date {
        Header
        // ...
}

struct Complex {
        Header
        // ...
}

variant.RecordHandler.JsonKind = func(bytes []byte) (string, error) {
        header := Header{}

        err := json.Unmarshal(bytes, &header)
        if err != nil {
                return "", err
        }

        return header.Kind, nil
}

variant.RecordHandler.NewJson = func(kind string, obj any) ([]byte, error) {
        switch o := obj.(type) {
        case *Date:
                o.Header.Kind = kind
        case *Complex:
                o.Header.Kind = kind
        }
        return json.Marshal(obj)
}

variant.RecordHandler.New = func(kind string) any {
        switch kind {
        case "date":
                return &Date{}
        case "complex":
                return &Complex{}
        }
        return nil
}

variant.RecordHandler.Kind = func(obj any) string {
        switch obj.(type) {
        case *Date:
                return "date"
        case *Complex:
                return "complex"
        }
        return ""
}
```

# licensing, etc

this library is in the public domain. 

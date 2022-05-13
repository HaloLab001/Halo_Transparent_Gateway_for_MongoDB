// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tjson

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// tjsontype is a type that can be marshaled to/from TJSON.
type tjsontype interface {
	tjsontype() // seal for go-sumtype

	Marshal([]byte, map[string]any) error     // tigris to build-in
	Unmarshal(map[string]any) ([]byte, error) // build-in to tigris.
}

//go-sumtype:decl tjsontype

// Unmarshal build-in to tigris.
func Unmarshal(v any, schema map[string]any) ([]byte, error) {
	fieldType, ok := schema["type"]
	if !ok {
		return nil, lazyerrors.Errorf("canont find field type")
	}
	switch v := v.(type) {
	case *types.Document:
		if fieldType != "object" {
			return nil, lazyerrors.Errorf("wrong schema %s for types.Document", fieldType)
		}
		d := documentType(*v)
		return d.Unmarshal(schema)

	// case *types.Array:
	case float64:
		d := doubleType(v)
		return d.Unmarshal(schema)
	case string:
		s := stringType(v)
		return s.Unmarshal(schema)
	case types.Binary:
		b := binaryType(v)
		return b.Unmarshal(schema)
	case types.ObjectID:
		o := objectIDType(v)
		return o.Unmarshal(schema)
	case bool:
		b := boolType(v)
		return b.Unmarshal(schema)
	case time.Time:
		t := dateTimeType(v)
		return t.Unmarshal(schema)
	case types.NullType:
		n := nullType(v)
		return n.Unmarshal(schema)
	case types.Regex:
		r := regexType(v)
		return r.Unmarshal(schema)
		// case int32:

	case types.Timestamp:
		t := timestampType(v)
		return t.Unmarshal(schema)
		// case int64:
	}
	return nil, lazyerrors.Errorf("%T is not supported", v)
}

// Marshal tigris to build-in.
func Marshal(v []byte, schema map[string]any) (any, error) {
	fieldType, ok := schema["type"]
	if !ok {
		return nil, lazyerrors.Errorf("canont find field type")
	}

	var err error
	var res any
	switch fieldType {
	case "object":
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			return nil, lazyerrors.Errorf("tjson.Document.Marshal: missing properties in schema")
		}
		if _, ok := properties["$b"]; ok {
			var o binaryType
			err = o.Marshal(v, schema)
			res = &o
			return res, err
		}
		if _, ok := properties["$o"]; ok {
			var o objectIDType
			err = o.Marshal(v, schema)
			res = &o
			return res, err
		}
		if _, ok := properties["$r"]; ok {
			var o regexType
			err = o.Marshal(v, schema)
			res = &o
			return res, err
		}
		if _, ok := properties["$t"]; ok {
			var o timestampType
			err = o.Marshal(v, schema)
			res = &o
			return res, err
		}
		var o documentType
		err = o.Marshal(v, schema)
		res = &o

	case "array":
		err = lazyerrors.Errorf("arrays not supported yet")

	case "boolean":
		var o boolType
		err = o.Marshal(v, schema)
		res = &o

	case "string":
		if format, ok := schema["format"]; ok {
			if format == "date-time" {
				var o dateTimeType
				err = o.Marshal(v, schema)
				res = &o
				return res, err
			}
		}
		res = string(v)

	default:
		err = lazyerrors.Errorf("tjson.Unmarshal: unhandled map %#v", v)
	}
	return res, err
}

// checkConsumed returns error if decoder or reader have buffered or unread data.
func checkConsumed(dec *json.Decoder, r *bytes.Reader) error {
	if dr := dec.Buffered().(*bytes.Reader); dr.Len() != 0 {
		b, _ := io.ReadAll(dr)
		return lazyerrors.Errorf("%d bytes remains in the decoded: %s", dr.Len(), b)
	}

	if l := r.Len(); l != 0 {
		b, _ := io.ReadAll(r)
		return lazyerrors.Errorf("%d bytes remains in the reader: %s", l, b)
	}

	return nil
}

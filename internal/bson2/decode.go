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

package bson2

import (
	"encoding/binary"
	"strconv"

	"github.com/cristalhq/bson/bsonproto"

	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

// DecodeDocument decodes a single BSON document that takes the whole b slice.
//
// Only first-level fields are decoded;
// nested documents and arrays are converted to RawDocument and RawArray respectively,
// using b subslices without copying.
func DecodeDocument(b []byte) (*Document, error) {
	bl := len(b)
	if bl < 5 {
		return nil, lazyerrors.Errorf("bl = %d: %w", bl, ErrDecodeShortInput)
	}
	if dl := int(binary.LittleEndian.Uint32(b)); bl != dl {
		return nil, lazyerrors.Errorf("bl = %d, dl = %d: %w", bl, dl, ErrDecodeInvalidInput)
	}
	if last := b[bl-1]; last != 0 {
		return nil, lazyerrors.Errorf("last = %d: %w", last, ErrDecodeInvalidInput)
	}

	res := MakeDocument(1)

	offset := 4
	for offset != len(b)-1 {
		t := tag(b[offset])
		offset++

		name, err := bsonproto.DecodeCString(b[offset:])
		offset += len(name) + 1
		if err != nil {
			return nil, lazyerrors.Error(err)
		}

		var v any
		switch t {
		case tagFloat64:
			v, err = bsonproto.DecodeFloat64(b[offset:])
			offset += bsonproto.SizeFloat64

		case tagString:
			var s string
			s, err = bsonproto.DecodeString(b[offset:])
			offset += bsonproto.SizeString(s)
			v = s

		case tagDocument:
			l := int(binary.LittleEndian.Uint32(b[offset:]))
			v = RawDocument(b[offset : offset+l])
			offset += l

		case tagArray:
			l := int(binary.LittleEndian.Uint32(b[offset:]))
			v = RawArray(b[offset : offset+l])
			offset += l

		case tagBinary:
			var s Binary
			s, err = bsonproto.DecodeBinary(b[offset:])
			offset += bsonproto.SizeBinary(s)
			v = s

		case tagObjectID:
			v, err = bsonproto.DecodeObjectID(b[offset:])
			offset += bsonproto.SizeObjectID

		case tagBool:
			v, err = bsonproto.DecodeBool(b[offset:])
			offset += bsonproto.SizeBool

		case tagTime:
			v, err = bsonproto.DecodeTime(b[offset:])
			offset += bsonproto.SizeTime

		case tagNull:
			v = Null

		case tagRegex:
			var s Regex
			s, err = bsonproto.DecodeRegex(b[offset:])
			offset += bsonproto.SizeRegex(s)
			v = s

		case tagInt32:
			v, err = bsonproto.DecodeInt32(b[offset:])
			offset += bsonproto.SizeInt32

		case tagTimestamp:
			v, err = bsonproto.DecodeTimestamp(b[offset:])
			offset += bsonproto.SizeTimestamp

		case tagInt64:
			v, err = bsonproto.DecodeInt64(b[offset:])
			offset += bsonproto.SizeInt64

		default:
			return nil, lazyerrors.Errorf("unsupported tag: %s", t)
		}

		if err != nil {
			return nil, lazyerrors.Error(err)
		}

		must.NoError(res.add(name, v))
	}

	return res, nil
}

// DecodeArray decodes a BSON array.
//
// Only first-level elements are decoded;
// nested documents and arrays are converted to RawDocument and RawArray respectively,
// using b subslices without copying.
func DecodeArray(b []byte) (*Array, error) {
	doc, err := DecodeDocument(b)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	res := &Array{
		elements: make([]any, len(doc.fields)),
	}

	for i, f := range doc.fields {
		if f.name != strconv.Itoa(i) {
			return nil, lazyerrors.Errorf("invalid array index: %q", f.name)
		}

		res.elements[i] = f.value
	}

	return res, nil
}

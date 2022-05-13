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

	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// doubleType represents BSON 64-bit binary floating point type.
type doubleType float64

// tjsontype implements tjsontype interface.
func (d *doubleType) tjsontype() {}

type doubleJSON float64

// Unmarshal build-in to tigris.
func (d *doubleType) Unmarshal(_ map[string]any) ([]byte, error) {
	f := float64(*d)
	res, err := json.Marshal(f)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	return res, nil
}

// Marshal tigris to build-in.
func (d *doubleType) Marshal(data []byte, _ map[string]any) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}
	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var o doubleJSON
	if err := dec.Decode(&o); err != nil {
		return lazyerrors.Error(err)
	}
	if err := checkConsumed(dec, r); err != nil {
		return lazyerrors.Error(err)
	}
	*d = doubleType(o)
	return nil
}

// check interfaces
var (
	_ tjsontype = (*doubleType)(nil)
)

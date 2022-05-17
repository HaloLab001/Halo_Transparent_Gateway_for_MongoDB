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

// boolType represents BSON Boolean type.
type boolType bool

// tjsontype implements tjsontype interface.
func (b *boolType) tjsontype() {}

var boolSchema = map[string]any{"type": "boolean"}

// Marshal build-in to tigris.
func (b *boolType) Marshal(_ map[string]any) ([]byte, error) {
	res, err := json.Marshal(bool(*b))
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	return res, nil
}

// Unmarshal tigris to build-in.
func (b *boolType) Unmarshal(data []byte, _ map[string]any) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var o bool
	if err := dec.Decode(&o); err != nil {
		return lazyerrors.Error(err)
	}
	if err := checkConsumed(dec, r); err != nil {
		return lazyerrors.Error(err)
	}

	*b = boolType(o)
	return nil
}

// check interfaces
var (
	_ tjsontype = (*boolType)(nil)
)

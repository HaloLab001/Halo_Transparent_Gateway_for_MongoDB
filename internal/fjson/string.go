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

package fjson

import (
	"bytes"
	"encoding/json"

	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// String represents BSON String data type.
type String string

func (str *String) fjsontype() {}

// UnmarshalJSON implements fjsontype interface.
func (str *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return lazyerrors.Error(err)
	}

	*str = String(s)
	return nil
}

// MarshalJSON implements fjsontype interface.
func (str String) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(string(str))
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return b, nil
}

// check interfaces
var (
	_ fjsontype = (*String)(nil)
)

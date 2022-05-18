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
	"time"

	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// dateTimeType represents BSON UTC datetime type.
type dateTimeType time.Time

// tjsontype implements tjsontype interface.
func (dt *dateTimeType) tjsontype() {}

// String returns formatted time for debugging.
func (dt *dateTimeType) String() string {
	return time.Time(*dt).Format(time.RFC3339Nano)
}

var dateTimeSchema = map[string]any{
	"type":   "string",
	"format": "date-time",
}

// string in RFC 3339.
type dateTimeJSON string

// Marshal build-in to tigris.
func (dt *dateTimeType) Marshal(_ map[string]any) ([]byte, error) {
	res, err := json.Marshal(dateTimeJSON(
		time.Time(*dt).Format(time.RFC3339Nano),
	))
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	return res, nil
}

// Unmarshal tigris to build-in.
func (dt *dateTimeType) Unmarshal(data []byte, _ map[string]any) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var o dateTimeJSON
	if err := dec.Decode(&o); err != nil {
		return lazyerrors.Error(err)
	}
	if err := checkConsumed(dec, r); err != nil {
		return lazyerrors.Error(err)
	}
	t, err := time.Parse(time.RFC3339Nano, string(o))
	if err != nil {
		return lazyerrors.Error(err)
	}
	*dt = dateTimeType(t)
	return nil
}

// check interfaces
var (
	_ tjsontype = (*dateTimeType)(nil)
)

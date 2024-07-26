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

package bson

import "github.com/FerretDB/FerretDB/internal/util/must"

// RawArray represents a single BSON array in the binary encoded form.
//
// It generally references a part of a larger slice, not a copy.
type RawArray []byte

// Encode returns itself to implement the [AnyArray] interface.
//
// Receiver must not be nil.
func (raw RawArray) Encode() (RawArray, error) {
	must.BeTrue(raw != nil)
	return raw, nil
}

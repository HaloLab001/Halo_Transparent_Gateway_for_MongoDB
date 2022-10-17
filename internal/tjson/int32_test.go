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
	"math"
	"testing"

	"github.com/AlekSi/pointer"
)

var int32TestCases = []testCase{{
	name:   "42",
	v:      pointer.To(int32Type(42)),
	schema: int32Schema,
	j:      `42`,
}, {
	name:   "zero",
	v:      pointer.To(int32Type(0)),
	schema: int32Schema,
	j:      `0`,
}, {
	name:   "max int32",
	v:      pointer.To(int32Type(math.MaxInt32)),
	schema: int32Schema,
	j:      `2147483647`,
}, {
	name:   "min int32",
	v:      pointer.To(int32Type(math.MinInt32)),
	schema: int32Schema,
	j:      `-2147483648`,
}, {
	name:   "schema mismatch",
	schema: boolSchema,
	v:      pointer.To(int32Type(42)),
	j:      `42`,
	sErr:   "json: cannot unmarshal number into Go value of type bool",
}}

func TestInt32(t *testing.T) {
	t.Parallel()
	testJSON(t, int32TestCases, func() tjsontype { return new(int32Type) })
}

func FuzzInt32(f *testing.F) {
	fuzzJSON(f, int32TestCases)
}

func BenchmarkInt32(b *testing.B) {
	benchmark(b, int32TestCases)
}

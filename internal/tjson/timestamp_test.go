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
	"testing"

	"github.com/FerretDB/FerretDB/internal/types"
)

var timestampTestCases = []testCase{{
	name: "one",
	v:    types.Timestamp(1652700697465990022),
	j:    `{"$t":"1652700697465990022"}`,
	s:    timestampSchema,
}, {
	name: "zero",
	v:    types.Timestamp(0),
	j:    `{"$t":"0"}`,
	s:    timestampSchema,
}, {
	name: "EOF",
	j:    `{`,
	jErr: `unexpected EOF`,
	s:    timestampSchema,
}}

func TestTimestamp(t *testing.T) {
	t.Parallel()
	testJSON(t, timestampTestCases, func() tjsontype { return new(timestampType) })
}

func FuzzTimestamp(f *testing.F) {
	fuzzJSON(f, timestampTestCases, func() tjsontype { return new(timestampType) })
}

func BenchmarkTimestamp(b *testing.B) {
	benchmark(b, timestampTestCases, func() tjsontype { return new(timestampType) })
}

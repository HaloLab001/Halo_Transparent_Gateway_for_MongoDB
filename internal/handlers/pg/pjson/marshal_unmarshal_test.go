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

package pjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/util/testutil"
)

func TestMarshalUnmarshal(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		doc  *types.Document
		json string
	}{
		"Empty": {
			json: `{"$s":{}}`,
			doc:  must.NotFail(types.NewDocument()),
		},
		"Filled": {
			json: `{
			"$s": {
				"p": {
					"foo": {"t": "string"}
				},
				"$k": ["foo"]
			}, 
			"foo": "bar"
		}`,
			doc: must.NotFail(types.NewDocument(
				"foo", "bar",
			)),
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			doc, err := Unmarshal([]byte(tc.json))
			require.NoError(t, err)
			assert.Equal(t, tc.doc, doc)

			actualB, err := Marshal(tc.doc)
			require.NoError(t, err)
			actualB = testutil.IndentJSON(t, actualB)

			expectedB := testutil.IndentJSON(t, []byte(tc.json))
			assert.Equal(t, string(expectedB), string(actualB))
		})
	}
}

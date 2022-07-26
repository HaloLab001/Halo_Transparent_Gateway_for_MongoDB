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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FerretDB/FerretDB/internal/util/testutil"
)

func TestSchemaMarshalUnmarshal(t *testing.T) {
	expected := Schema{
		Title:       "users",
		Description: "FerretDB users collection",
		Properties: map[string]*Schema{
			"$k":      {Type: Array, Items: stringSchema},
			"_id":     objectIDSchema,
			"name":    stringSchema,
			"balance": doubleSchema,
			"data":    binarySchema,
		},
		PrimaryKey: []string{"_id"},
	}

	actualB, err := expected.Marshal()
	require.NoError(t, err)
	actualB = testutil.IndentJSON(t, actualB)

	expectedB := testutil.IndentJSON(t, []byte(`{
		"title": "users",
		"description": "FerretDB users collection",
		"properties": {
			"$k": {"type": "array", "items": {"type": "string"}},
			"_id": {"type": "string", "format": "byte"},
			"balance": {"type": "number"},
			"data": {
				"type": "object",
				"properties": {
					"$b": {"type": "string", "format": "byte"},
					"s": {"type": "integer", "format": "int32"}
				}
			},
			"name": {"type": "string"}
		},
		"primary_key": ["_id"]
	}`))
	assert.Equal(t, string(expectedB), string(actualB))

	var actual Schema
	err = actual.Unmarshal(expectedB)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestSchemaEqual(t *testing.T) {
	t.Parallel()

	cInt64Schema := Schema{
		Type:   Integer,
		Format: Int64,
	}
	cIntEmptySchema := Schema{
		Type:   Integer,
		Format: EmptyFormat,
	}
	cDoubleSchema := Schema{
		Type:   Number,
		Format: Double,
	}
	cDoubleEmptySchema := Schema{
		Type:   Number,
		Format: EmptyFormat,
	}
	cObjectSchema := Schema{
		Type: Object,
		Properties: map[string]*Schema{
			"a":  stringSchema,
			"42": &cIntEmptySchema,
		},
	}
	cObjectSchemaEqual := Schema{
		Type: Object,
		Properties: map[string]*Schema{
			"42": &cIntEmptySchema,
			"a":  stringSchema,
		},
	}
	cObjectSchemaNotEqual := Schema{
		Type: Object,
		Properties: map[string]*Schema{
			"42": &cIntEmptySchema,
			"a":  boolSchema,
		},
	}
	cObjectSchemaKeyMissing := Schema{
		Type: Object,
		Properties: map[string]*Schema{
			"42": &cIntEmptySchema,
			"b":  stringSchema,
		},
	}
	cObjectSchemaEmpty := Schema{
		Type:       Object,
		Properties: map[string]*Schema{},
	}
	cArrayDoubleSchema := Schema{
		Type:  Array,
		Items: &cDoubleSchema,
	}
	cArrayDoubleEmptySchema := Schema{
		Type:  Array,
		Items: &cDoubleEmptySchema,
	}
	cArrayObjectsSchema := Schema{
		Type:  Array,
		Items: &cObjectSchema,
	}
	cArrayObjectsSchemaEqual := Schema{
		Type:  Array,
		Items: &cObjectSchemaEqual,
	}
	cArrayObjectsSchemaNotEqual := Schema{
		Type:  Array,
		Items: &cObjectSchemaNotEqual,
	}

	for name, tc := range map[string]struct {
		s        *Schema
		other    *Schema
		expected bool
	}{"StringString": {
		s:        stringSchema,
		other:    stringSchema,
		expected: true,
	}, "StringNumber": {
		s:        stringSchema,
		other:    doubleSchema,
		expected: false,
	}, "NumberString": {
		s:        doubleSchema,
		other:    stringSchema,
		expected: false,
	}, "EmptyInt64": {
		s:        &cIntEmptySchema,
		other:    &cInt64Schema,
		expected: true,
	}, "Int64Empty": {
		s:        &cInt64Schema,
		other:    &cIntEmptySchema,
		expected: true,
	}, "Int64Int32": {
		s:        &cInt64Schema,
		other:    int32Schema,
		expected: false,
	}, "EmptyInt32": {
		s:        &cIntEmptySchema,
		other:    int32Schema,
		expected: false,
	}, "DoubleEmpty": {
		s:        &cDoubleSchema,
		other:    &cDoubleEmptySchema,
		expected: true,
	}, "ObjectsEqual": {
		s:        &cObjectSchema,
		other:    &cObjectSchemaEqual,
		expected: true,
	}, "ObjectsNotEqual": {
		s:        &cObjectSchemaEqual,
		other:    &cObjectSchemaNotEqual,
		expected: false,
	}, "ObjectsKeyMissing": {
		s:        &cObjectSchema,
		other:    &cObjectSchemaKeyMissing,
		expected: false,
	}, "ObjectsEmpty": {
		s:        &cObjectSchema,
		other:    &cObjectSchemaEmpty,
		expected: false,
	}, "ArrayDouble": {
		s:        &cArrayDoubleSchema,
		other:    &cArrayDoubleEmptySchema,
		expected: true,
	}, "ArrayObjects": {
		s:        &cArrayObjectsSchema,
		other:    &cArrayObjectsSchemaEqual,
		expected: true,
	}, "ArrayObjectsNotEqual": {
		s:        &cArrayObjectsSchemaNotEqual,
		other:    &cArrayObjectsSchemaEqual,
		expected: false,
	}, "ArrayObjectsDouble": {
		s:        &cArrayObjectsSchema,
		other:    &cArrayDoubleSchema,
		expected: false,
	}} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.s.Equal(tc.other))
		})
	}
}

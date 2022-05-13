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

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/FerretDB/FerretDB/integration/shareddata"
)

func TestUpdateUpsert(t *testing.T) {
	t.Parallel()
	ctx, collection := setup(t, shareddata.Composites)

	// this upsert inserts document
	filter := bson.D{{"foo", "bar"}}
	update := bson.D{{"$set", bson.D{{"foo", "baz"}}}}
	res, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	require.NoError(t, err)

	id := res.UpsertedID
	assert.NotEmpty(t, id)
	res.UpsertedID = nil
	expected := &mongo.UpdateResult{
		MatchedCount:  0,
		ModifiedCount: 0,
		UpsertedCount: 1,
	}
	require.Equal(t, expected, res)

	// check inserted document
	var doc bson.D
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&doc)
	require.NoError(t, err)
	if !AssertEqualDocuments(t, bson.D{{"_id", id}, {"foo", "baz"}}, doc) {
		t.FailNow()
	}

	// this upsert updates document
	filter = bson.D{{"foo", "baz"}}
	update = bson.D{{"$set", bson.D{{"foo", "qux"}}}}
	res, err = collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	require.NoError(t, err)

	expected = &mongo.UpdateResult{
		MatchedCount:  1,
		ModifiedCount: 1,
		UpsertedCount: 0,
	}
	require.Equal(t, expected, res)

	// check updated document
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&doc)
	require.NoError(t, err)
	AssertEqualDocuments(t, bson.D{{"_id", id}, {"foo", "qux"}}, doc)
}

func TestUpdateIncOperatorErrors(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		filter bson.D
		update bson.D
		err    *mongo.WriteError
	}{
		"BadIncType": {
			filter: bson.D{{"_id", "string"}},
			update: bson.D{{"$inc", bson.D{{"value", "bad value"}}}},
			err: &mongo.WriteError{
				Code:    14,
				Message: `Cannot increment with non-numeric argument: {value: "bad value"}`,
			},
		},
		"IncOnNullValue": {
			filter: bson.D{{"_id", "string"}},
			update: bson.D{{"$inc", bson.D{{"value", int32(1)}}}},
			err: &mongo.WriteError{
				Code:    14,
				Message: `Cannot apply $inc to a value of non-numeric type. {_id: "string"} has the field 'value' of non-numeric type string`,
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, collection := setup(t, shareddata.Scalars)

			_, err := collection.UpdateOne(ctx, tc.filter, tc.update)
			if tc.err != nil {
				AssertEqualWriteError(t, tc.err, err)
				return
			}
			require.NoError(t, err)

			var actual bson.D
			err = collection.FindOne(ctx, tc.filter).Decode(&actual)
			require.NoError(t, err)

			t.Log(actual)
		})
	}
}

func TestUpdateIncOperator(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		filter bson.D
		update bson.D
		result bson.D
	}{
		"DoubleDoubleIncrement": {
			filter: bson.D{{"_id", "double"}},
			update: bson.D{{"$inc", bson.D{{"value", float64(42.13)}}}},
			result: bson.D{{"_id", "double"}, {"value", float64(84.26)}},
		},
		"DoubleDoubleNegativeIncrement": {
			filter: bson.D{{"_id", "double"}},
			update: bson.D{{"$inc", bson.D{{"value", float64(-42.13)}}}},
			result: bson.D{{"_id", "double"}, {"value", float64(0)}},
		},
		"DoubleIncrementIntField": {
			filter: bson.D{{"_id", "int32"}},
			update: bson.D{{"$inc", bson.D{{"value", float64(1.13)}}}},
			result: bson.D{{"_id", "int32"}, {"value", float64(43.13)}},
		},
		"DoubleIntIncrement": {
			filter: bson.D{{"_id", "double"}},
			update: bson.D{{"$inc", bson.D{{"value", int32(1)}}}},
			result: bson.D{{"_id", "double"}, {"value", float64(43.13)}},
		},
		"Int": {
			filter: bson.D{{"_id", "int32"}},
			update: bson.D{{"$inc", bson.D{{"value", int32(1)}}}},
			result: bson.D{{"_id", "int32"}, {"value", int32(43)}},
		},
		"IntNegativeIncrement": {
			filter: bson.D{{"_id", "int32"}},
			update: bson.D{{"$inc", bson.D{{"value", int32(-1)}}}},
			result: bson.D{{"_id", "int32"}, {"value", int32(41)}},
		},
		"Long": {
			filter: bson.D{{"_id", "int64"}},
			update: bson.D{{"$inc", bson.D{{"value", int64(1)}}}},
			result: bson.D{{"_id", "int64"}, {"value", int64(43)}},
		},
		"LongNegativeIncrement": {
			filter: bson.D{{"_id", "int64"}},
			update: bson.D{{"$inc", bson.D{{"value", int64(-1)}}}},
			result: bson.D{{"_id", "int64"}, {"value", int64(41)}},
		},

		"FieldNotExist": {
			filter: bson.D{{"_id", "int32"}},
			update: bson.D{{"$inc", bson.D{{"foo", int32(1)}}}},
			result: bson.D{{"_id", "int32"}, {"value", int32(42)}, {"foo", int32(1)}},
		},

		"DotNotationInt": {
			filter: bson.D{{"_id", "document"}},
			update: bson.D{{"$inc", bson.D{{"value.foo", int32(1)}}}},
			result: bson.D{{"_id", "document"}, {"value", bson.D{{"foo", int32(43)}}}},
		},
		"DotNotationFieldNotExist": {
			filter: bson.D{{"_id", "document"}},
			update: bson.D{{"$inc", bson.D{{"value.no-such-field", int32(1)}}}},
			result: bson.D{
				{"_id", "document"},
				{"value", bson.D{{"foo", int32(42)}, {"no-such-field", int32(1)}}},
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, collection := setup(t, shareddata.Scalars, shareddata.Composites)

			_, err := collection.UpdateOne(ctx, tc.filter, tc.update)
			require.NoError(t, err)

			var actual bson.D
			err = collection.FindOne(ctx, tc.filter).Decode(&actual)
			require.NoError(t, err)

			AssertEqualDocuments(t, tc.result, actual)
		})
	}
}

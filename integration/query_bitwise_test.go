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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/FerretDB/FerretDB/integration/shareddata"
)

func TestQueryBitwiseAllClear(t *testing.T) {
	t.Parallel()
	ctx, collection := setup(t, shareddata.Scalars)

	_, err := collection.InsertMany(ctx, []any{
		bson.D{{"_id", "binary-big"}, {"value", primitive.Binary{Data: []byte{0, 0, 128}}}},
	})
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		filter      any
		expectedIDs []any
		err         mongo.CommandError
	}{
		"Float": {
			filter: 1.2,
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected an integer: $bitsAllClear: 1.2",
			},
		},
		"String": {
			filter: "123",
			err: mongo.CommandError{
				Code:    2,
				Name:    "BadValue",
				Message: "value takes an Array, a number, or a BinData but received: $bitsAllClear: \"123\"",
			},
		},
		"Int32": {
			filter: int32(2),
			expectedIDs: []any{
				"binary-big", "binary-empty",
				"double-negative-zero", "double-zero",
				"int32-min", "int32-zero",
				"int64-min", "int64-zero",
			},
		},
		"Int64": {
			filter:      math.MaxInt64,
			expectedIDs: []any{},
		},
		"Array": {
			filter:      primitive.A{1, 5},
			expectedIDs: []any{},
		},
		"NegativeValue": {
			filter: int32(-1),
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected a positive number in: $bitsAllClear: -1",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filter := bson.D{{"value", bson.D{{"$bitsAllClear", tc.filter}}}}
			cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"_id", 1}}))
			if tc.err.Code != 0 {
				require.Nil(t, tc.expectedIDs)
				assertEqualError(t, tc.err, err)
				return
			}
			require.NoError(t, err)

			var actual []bson.D
			err = cursor.All(ctx, &actual)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedIDs, collectIDs(t, actual))
		})
	}
}

func TestQueryBitwiseAllSet(t *testing.T) {
	t.Parallel()
	ctx, collection := setup(t, shareddata.Scalars)

	_, err := collection.InsertMany(ctx, []any{
		bson.D{{"_id", "binary-big"}, {"value", primitive.Binary{Data: []byte{0, 0, 128}}}},
	})
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		filter      any
		expectedIDs []any
		err         mongo.CommandError
	}{
		"Int32": {
			filter:      int32(2),
			expectedIDs: []any{"binary", "double-whole", "int32", "int32-max", "int64", "int64-max"},
		},
		"String": {
			filter: "123",
			err: mongo.CommandError{
				Code:    2,
				Name:    "BadValue",
				Message: "value takes an Array, a number, or a BinData but received: $bitsAllSet: \"123\"",
			},
		},
		"Float": {
			filter: 1.2,
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected an integer: $bitsAllSet: 1.2",
			},
		},
		"NegativeValue": {
			filter: int32(-1),
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected a positive number in: $bitsAllSet: -1",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filter := bson.D{{"value", bson.D{{"$bitsAllSet", tc.filter}}}}
			cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"_id", 1}}))
			if tc.err.Code != 0 {
				require.Nil(t, tc.expectedIDs)
				assertEqualError(t, tc.err, err)
				return
			}
			require.NoError(t, err)

			var actual []bson.D
			err = cursor.All(ctx, &actual)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedIDs, collectIDs(t, actual))
		})
	}
}

func TestQueryBitwiseAnyClear(t *testing.T) {
	t.Parallel()
	ctx, collection := setup(t, shareddata.Scalars)

	_, err := collection.InsertMany(ctx, []any{
		bson.D{{"_id", "binary-big"}, {"value", primitive.Binary{Data: []byte{0, 0, 128}}}},
	})
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		filter      any
		expectedIDs []any
		err         mongo.CommandError
	}{
		"Int32": {
			filter: int32(2),
			expectedIDs: []any{
				"binary-empty", "double-negative-zero", "double-zero",
				"int32-min", "int32-zero", "int64-min", "int64-zero",
			},
		},

		//"BitsAnyClear": {
		//	filter:      int32(1),
		//	expectedIDs: []any{"int32"},
		//},
		//"BitsAnyClearEmpty": {
		//	filter:      int32(42),
		//	expectedIDs: []any{},
		//},
		"BitsAnyClearBigBinary": {
			filter:      int32(0b1000_0000_0000_0000),
			expectedIDs: []any{"binary-big"},
		},
		//"BitsAnyClearBigBinaryEmptyResult": {
		//	filter:      int32(0b1000_0000_0000_0000_0000_0000),
		//	expectedIDs: []any{},
		//},
		"BitsAllSetString": {
			filter: "123",
			err: mongo.CommandError{
				Code:    2,
				Name:    "BadValue",
				Message: "value takes an Array, a number, or a BinData but received: $bitsAllSet: \"123\"",
			},
		},
		"BitsAllSetPassFloat": {
			filter: 1.2,
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected an integer: $bitsAllSet: 1.2",
			},
		},
		"BitsAllSetPassNegativeValue": {
			filter: int32(-1),
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected a positive number in: $bitsAllSet: -1",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filter := bson.D{{"value", bson.D{{"$bitsAnyClear", tc.filter}}}}
			cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"_id", 1}}))
			if tc.err.Code != 0 {
				require.Nil(t, tc.expectedIDs)
				assertEqualError(t, tc.err, err)
				return
			}
			require.NoError(t, err)

			var actual []bson.D
			err = cursor.All(ctx, &actual)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedIDs, collectIDs(t, actual))
		})
	}
}

func TestQueryBitwiseAnySet(t *testing.T) {
	t.Parallel()
	ctx, collection := setup(t, shareddata.Scalars)

	_, err := collection.InsertMany(ctx, []any{
		bson.D{{"_id", "binary-big"}, {"value", primitive.Binary{Data: []byte{0, 0, 128}}}},
	})
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		filter      any
		expectedIDs []any
		err         mongo.CommandError
	}{
		"Int32": {
			filter:      int32(2),
			expectedIDs: []any{"binary", "double-whole", "int32", "int32-max", "int64", "int64-max"},
		},
		"BitsAnySetBigBinary": {
			filter:      int32(0b1000_0000_0000_0000_0000_0000),
			expectedIDs: []any{"binary-big"},
		},
		"BitsAnySetBigBinaryEmptyResult": {
			filter:      int32(0b1000_0000_0000_0000),
			expectedIDs: []any{},
		},
		"BitsAllSetString": {
			filter: "123",
			err: mongo.CommandError{
				Code:    2,
				Name:    "BadValue",
				Message: "value takes an Array, a number, or a BinData but received: $bitsAllSet: \"123\"",
			},
		},
		"BitsAllSetPassFloat": {
			filter: 1.2,
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected an integer: $bitsAllSet: 1.2",
			},
		},
		"BitsAllSetPassNegativeValue": {
			filter: int32(-1),
			err: mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Expected a positive number in: $bitsAllSet: -1",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filter := bson.D{{"value", bson.D{{"$bitsAnySet", tc.filter}}}}
			cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"_id", 1}}))
			if tc.err.Code != 0 {
				require.Nil(t, tc.expectedIDs)
				assertEqualError(t, tc.err, err)
				return
			}
			require.NoError(t, err)

			var actual []bson.D
			err = cursor.All(ctx, &actual)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedIDs, collectIDs(t, actual))
		})
	}
}

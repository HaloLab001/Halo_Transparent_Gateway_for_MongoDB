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

	"github.com/FerretDB/FerretDB/integration/setup"
	"github.com/FerretDB/FerretDB/integration/shareddata"
)

func TestFindAndModifyEmptyCollectionName(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		err        *mongo.CommandError
		altMessage string
	}{
		"EmptyCollectionName": {
			err: &mongo.CommandError{
				Code:    73,
				Message: "Invalid namespace specified 'testfindandmodifyemptycollectionname_emptycollectionname.'",
				Name:    "InvalidNamespace",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, collection := setup.Setup(t, shareddata.Doubles)

			var actual bson.D
			err := collection.Database().RunCommand(ctx, bson.D{{"findAndModify", ""}}).Decode(&actual)

			AssertEqualError(t, *tc.err, err)
		})
	}
}

func TestFindAndModifyErrors(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		command    bson.D
		err        *mongo.CommandError
		altMessage string
	}{
		"UpsertAndRemove": {
			command: bson.D{
				{"upsert", true},
				{"remove", true},
			},
			err: &mongo.CommandError{
				Code:    9,
				Name:    "FailedToParse",
				Message: "Cannot specify both upsert=true and remove=true ",
			},
			altMessage: "Cannot specify both upsert=true and remove=true",
		},
		"BadSortType": {
			command: bson.D{
				{"update", bson.D{}},
				{"sort", "123"},
			},
			err: &mongo.CommandError{
				Code:    14,
				Name:    "TypeMismatch",
				Message: "BSON field 'findAndModify.sort' is the wrong type 'string', expected type 'object'",
			},
			altMessage: "BSON field 'sort' is the wrong type 'string', expected type 'object'",
		},
		"BadRemoveType": {
			command: bson.D{
				{"query", bson.D{}},
				{"remove", "123"},
			},
			err: &mongo.CommandError{
				Code:    14,
				Name:    "TypeMismatch",
				Message: "BSON field 'findAndModify.remove' is the wrong type 'string', expected types '[bool, long, int, decimal, double']",
			},
			altMessage: "BSON field 'remove' is the wrong type 'string', expected types '[bool, long, int, decimal, double]'",
		},
		"BadNewType": {
			command: bson.D{
				{"query", bson.D{}},
				{"new", "123"},
			},
			err: &mongo.CommandError{
				Code:    14,
				Name:    "TypeMismatch",
				Message: "BSON field 'findAndModify.new' is the wrong type 'string', expected types '[bool, long, int, decimal, double']",
			},
			altMessage: "BSON field 'new' is the wrong type 'string', expected types '[bool, long, int, decimal, double]'",
		},
		"BadUpsertType": {
			command: bson.D{
				{"query", bson.D{}},
				{"upsert", "123"},
			},
			err: &mongo.CommandError{
				Code:    14,
				Name:    "TypeMismatch",
				Message: "BSON field 'findAndModify.upsert' is the wrong type 'string', expected types '[bool, long, int, decimal, double']",
			},
			altMessage: "BSON field 'upsert' is the wrong type 'string', expected types '[bool, long, int, decimal, double]'",
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, collection := setup.Setup(t, shareddata.DocumentsStrings)

			command := bson.D{{"findAndModify", collection.Name()}}
			command = append(command, tc.command...)
			if command.Map()["sort"] == nil {
				command = append(command, bson.D{{"sort", bson.D{{"_id", 1}}}}...)
			}

			var actual bson.D
			err := collection.Database().RunCommand(ctx, command).Decode(&actual)

			AssertEqualAltError(t, *tc.err, tc.altMessage, err)
		})
	}
}

func TestFindAndModifyUpsert(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		command  bson.D
		response bson.D
	}{
		"Upsert": {
			command: bson.D{
				{"query", bson.D{{"_id", "double"}}},
				{"update", bson.D{{"$set", bson.D{{"v", 43.13}}}}},
				{"upsert", true},
			},
			response: bson.D{
				{"lastErrorObject", bson.D{
					{"n", int32(1)},
					{"updatedExisting", true},
				}},
				{"value", bson.D{{"_id", "double"}, {"v", 42.13}}},
				{"ok", float64(1)},
			},
		},
		"UpsertNew": {
			command: bson.D{
				{"query", bson.D{{"_id", "double"}}},
				{"update", bson.D{{"$set", bson.D{{"v", 43.13}}}}},
				{"upsert", true},
				{"new", true},
			},
			response: bson.D{
				{"lastErrorObject", bson.D{
					{"n", int32(1)},
					{"updatedExisting", true},
				}},
				{"value", bson.D{{"_id", "double"}, {"v", 43.13}}},
				{"ok", float64(1)},
			},
		},
		"UpsertNoSuchDocument": {
			command: bson.D{
				{"query", bson.D{{"_id", "no-such-doc"}}},
				{"update", bson.D{{"$set", bson.D{{"v", 43.13}}}}},
				{"upsert", true},
				{"new", true},
			},
			response: bson.D{
				{"lastErrorObject", bson.D{
					{"n", int32(1)},
					{"updatedExisting", false},
					{"upserted", "no-such-doc"},
				}},
				{"value", bson.D{{"_id", "no-such-doc"}, {"v", 43.13}}},
				{"ok", float64(1)},
			},
		},
		"UpsertNoSuchReplaceDocument": {
			command: bson.D{
				{"query", bson.D{{"_id", "no-such-doc"}}},
				{"update", bson.D{{"v", 43.13}}},
				{"upsert", true},
				{"new", true},
			},
			response: bson.D{
				{"lastErrorObject", bson.D{
					{"n", int32(1)},
					{"updatedExisting", false},
					{"upserted", "no-such-doc"},
				}},
				{"value", bson.D{{"_id", "no-such-doc"}, {"v", 43.13}}},
				{"ok", float64(1)},
			},
		},
		"UpsertReplace": {
			command: bson.D{
				{"query", bson.D{{"_id", "double"}}},
				{"update", bson.D{{"v", 43.13}}},
				{"upsert", true},
			},
			response: bson.D{
				{"lastErrorObject", bson.D{
					{"n", int32(1)},
					{"updatedExisting", true},
				}},
				{"value", bson.D{{"_id", "double"}, {"v", 42.13}}},
				{"ok", float64(1)},
			},
		},
		"UpsertReplaceReturnNew": {
			command: bson.D{
				{"query", bson.D{{"_id", "double"}}},
				{"update", bson.D{{"v", 43.13}}},
				{"upsert", true},
				{"new", true},
			},
			response: bson.D{
				{"lastErrorObject", bson.D{
					{"n", int32(1)},
					{"updatedExisting", true},
				}},
				{"value", bson.D{{"_id", "double"}, {"v", 43.13}}},
				{"ok", float64(1)},
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, collection := setup.Setup(t, shareddata.Doubles)

			command := append(bson.D{{"findAndModify", collection.Name()}}, tc.command...)

			var actual bson.D
			err := collection.Database().RunCommand(ctx, command).Decode(&actual)
			require.NoError(t, err)

			m := actual.Map()
			assert.Equal(t, float64(1), m["ok"])

			AssertEqualDocuments(t, tc.response, actual)
		})
	}
}

func TestFindAndModifyUpsertComplex(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		command         bson.D
		lastErrorObject bson.D
	}{
		"UpsertNoSuchDocumentNoIdInQuery": {
			command: bson.D{
				{"query", bson.D{{
					"$and",
					bson.A{
						bson.D{{"v", bson.D{{"$gt", 0}}}},
						bson.D{{"v", bson.D{{"$lt", 0}}}},
					},
				}}},
				{"update", bson.D{{"$set", bson.D{{"v", 43.13}}}}},
				{"upsert", true},
			},
			lastErrorObject: bson.D{
				{"n", int32(1)},
				{"updatedExisting", false},
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, collection := setup.Setup(t, shareddata.Doubles)

			command := append(bson.D{{"findAndModify", collection.Name()}}, tc.command...)

			var actual bson.D
			err := collection.Database().RunCommand(ctx, command).Decode(&actual)
			require.NoError(t, err)

			m := actual.Map()
			assert.Equal(t, float64(1), m["ok"])

			leb, ok := m["lastErrorObject"].(bson.D)
			if !ok {
				t.Fatal(actual)
			}

			// TODO: add document comparison here. Skip _id check as it always would different.
			for _, v := range leb {
				if v.Key == "upserted" {
					continue
				}
				assert.Equal(t, tc.lastErrorObject.Map()[v.Key], v.Value)
			}
		})
	}
}

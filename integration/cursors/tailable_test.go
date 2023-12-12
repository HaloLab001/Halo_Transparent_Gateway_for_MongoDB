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

package cursors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/FerretDB/FerretDB/integration"
	"github.com/FerretDB/FerretDB/integration/setup"
	"github.com/FerretDB/FerretDB/integration/shareddata"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

func TestTailableErrors(t *testing.T) {
	t.Parallel()

	t.Run("NonCapped", func(t *testing.T) {
		t.Parallel()

		ctx, collection := setup.Setup(t, shareddata.Scalars)

		for _, ct := range []options.CursorType{options.Tailable, options.TailableAwait} {
			cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetCursorType(ct))
			expected := mongo.CommandError{
				Code: 2,
				Name: "BadValue",
				Message: "error processing query: " +
					"ns=TestTailable-NonCapped.TestTailable-NonCappedTree: $and\nSort: {}\nProj: {}\n " +
					"tailable cursor requested on non capped collection",
			}
			integration.AssertEqualAltCommandError(t, expected, "tailable cursor requested on non capped collection", err)
			assert.Nil(t, cursor)
		}
	})
}

func TestTailableGetMore(t *testing.T) {
	s := setup.SetupWithOpts(t, &setup.SetupOpts{})

	db, ctx := s.Collection.Database(), s.Ctx

	opts := options.CreateCollection().SetCapped(true).SetSizeInBytes(10000)
	err := db.CreateCollection(s.Ctx, t.Name(), opts)
	require.NoError(t, err)

	collection := db.Collection(t.Name())

	bsonArr, arr := integration.GenerateDocuments(0, 3)

	_, err = collection.InsertMany(ctx, bsonArr)
	require.NoError(t, err)

	cmd := bson.D{
		{"find", collection.Name()},
		{"batchSize", 1},
		{"tailable", true},
	}

	var res bson.D
	err = collection.Database().RunCommand(ctx, cmd).Decode(&res)
	require.NoError(t, err)

	firstBatch, cursorID := GetFirstBatch(t, res)

	expectedFirstBatch := integration.ConvertDocuments(t, arr[:1])
	require.Equal(t, len(expectedFirstBatch), firstBatch.Len())
	require.Equal(t, expectedFirstBatch[0], must.NotFail(firstBatch.Get(0)))

	getMoreCmd := bson.D{
		{"getMore", cursorID},
		{"collection", collection.Name()},
		{"batchSize", 1},
	}

	for i := 0; i < 2; i++ {
		var res bson.D
		err = collection.Database().RunCommand(ctx, getMoreCmd).Decode(&res)
		require.NoError(t, err)

		nextBatch, nextID := GetNextBatch(t, res)
		expectedNextBatch := integration.ConvertDocuments(t, arr[i+1:i+2])

		assert.Equal(t, cursorID, nextID, res)

		require.Equal(t, len(expectedNextBatch), nextBatch.Len())
		require.Equal(t, expectedNextBatch[0], must.NotFail(nextBatch.Get(0)))
	}
}

func GetFirstBatch(t testing.TB, res bson.D) (*types.Array, any) {
	t.Helper()

	doc := integration.ConvertDocument(t, res)

	v, _ := doc.Get("cursor")
	require.NotNil(t, v)

	cursor, ok := v.(*types.Document)
	require.True(t, ok)

	cursorID, _ := cursor.Get("id")
	assert.NotNil(t, cursorID)

	v, _ = cursor.Get("firstBatch")
	require.NotNil(t, v)

	firstBatch, ok := v.(*types.Array)
	require.True(t, ok)

	return firstBatch, cursorID
}

func GetNextBatch(t testing.TB, res bson.D) (*types.Array, any) {
	t.Helper()

	doc := integration.ConvertDocument(t, res)

	v, _ := doc.Get("cursor")
	require.NotNil(t, v)

	cursor, ok := v.(*types.Document)
	require.True(t, ok)

	cursorID, _ := cursor.Get("id")
	assert.NotNil(t, cursorID)

	v, _ = cursor.Get("nextBatch")
	require.NotNil(t, v)

	firstBatch, ok := v.(*types.Array)
	require.True(t, ok)

	return firstBatch, cursorID
}

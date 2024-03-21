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

package driver

import (
	"testing"

	"github.com/cristalhq/bson/bsonproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FerretDB/FerretDB/internal/bson"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/util/testutil"
	"github.com/FerretDB/FerretDB/internal/wire"
)

func TestDriver(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in -short mode")
	}

	ctx := testutil.Ctx(t)

	c, err := Connect(ctx, "mongodb://127.0.0.1:47017/", testutil.SLogger(t))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, c.Close()) })

	dbName := "TestDriver"

	var lsid bson.Binary

	lsid = bsonproto.Binary{
		B: []byte{
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			// 0xa3, 0x19, 0xf2, 0xb4, 0xa1, 0x75, 0x40, 0xc7,
			// 0xb8, 0xe7, 0xa3, 0xa3, 0x2e, 0xc2, 0x56, 0xbe,
		},
		Subtype: bsonproto.BinaryUUID,
	}

	t.Run("Drop", func(t *testing.T) {
		doc := must.NotFail(bson.NewDocument(
			"dropDatabase", int32(1),
			"lsid", must.NotFail(bson.NewDocument("id", lsid)),
			"$db", dbName,
		))

		body, err := wire.NewOpMsg(must.NotFail(doc.Encode()))
		require.NoError(t, err)

		// TODO verify header
		header := wire.MsgHeader{
			RequestID: 13,
			OpCode:    wire.OpCodeMsg,
		}

		_, _, err = c.Request(ctx, &header, body)
		require.NoError(t, err)
	})

	insertDocs := must.NotFail(bson.NewArray(
		must.NotFail(bson.NewDocument("w", int32(2), "v", int32(1), "_id", int32(0))),
		must.NotFail(bson.NewDocument("v", int32(2), "_id", int32(1))),
		must.NotFail(bson.NewDocument("v", int32(3), "_id", int32(2))),
	))

	t.Run("Insert", func(t *testing.T) {
		insertCmd := must.NotFail(bson.NewDocument(
			"insert", "values",
			"documents", insertDocs,
			"ordered", true,
			"lsid", must.NotFail(bson.NewDocument("id", lsid)),
			"$db", dbName,
		))

		body, err := wire.NewOpMsg(must.NotFail(insertCmd.Encode()))
		require.NoError(t, err)

		// TODO verify header
		header := wire.MsgHeader{
			RequestID: 13,
			OpCode:    wire.OpCodeMsg,
		}

		//resHeader, resBody, err := c.Request(ctx, &header, body)
		_, _, err = c.Request(ctx, &header, body)
		require.NoError(t, err)

		//	assert.Equal(t, wire.MsgHeader{
		//		RequestID:     368,
		//		ResponseTo:    13,
		//		OpCode:        2013,
		//		MessageLength: 252,
		//	}, *resHeader)
		//	assert.Equal(t, wire.OpMsg{}, resBody)
	})

	var cursorID int64

	t.Run("Find", func(t *testing.T) {
		findCmd := must.NotFail(bson.NewDocument(
			"find", "values",
			//"filter", must.NotFail(bson.NewDocument("v", int32(1))),
			"filter", must.NotFail(bson.NewDocument()),
			//"sort", // TODO
			"lsid", must.NotFail(bson.NewDocument("id", lsid)),
			"batchSize", int32(1),
			"$db", dbName,
		))

		body, err := wire.NewOpMsg(must.NotFail(findCmd.Encode()))
		require.NoError(t, err)

		// TODO verify header
		header := wire.MsgHeader{
			RequestID: 14,
			OpCode:    wire.OpCodeMsg,
		}

		_, resBody, err := c.Request(ctx, &header, body)
		require.NoError(t, err)

		resMsg, err := must.NotFail(resBody.(*wire.OpMsg).RawDocument()).Decode()
		require.NoError(t, err)

		cursor, err := resMsg.Get("cursor").(bson.RawDocument).Decode()
		require.NoError(t, err)

		firstBatch := cursor.Get("firstBatch").(bson.RawArray)
		cursorID = cursor.Get("id").(int64)

		expectedDocs := must.NotFail(bson.NewArray(
			must.NotFail(bson.NewDocument("_id", int32(0), "w", int32(2), "v", int32(1))),
		))

		testutil.AssertEqual(t, must.NotFail(expectedDocs.Convert()), must.NotFail(firstBatch.Convert()))
		require.NotZero(t, cursorID)
	})

	getMoreCmd := must.NotFail(bson.NewDocument(
		"getMore", cursorID,
		"collection", "values",
		"lsid", must.NotFail(bson.NewDocument("id", lsid)),
		"batchSize", int32(1),
		"$db", dbName,
	))

	t.Run("GetMore", func(t *testing.T) {
		for i := 0; i < 2; i++ {

			body, err := wire.NewOpMsg(must.NotFail(getMoreCmd.Encode()))
			require.NoError(t, err)

			// TODO verify header
			header := wire.MsgHeader{
				RequestID: 14,
				OpCode:    wire.OpCodeMsg,
			}

			_, resBody, err := c.Request(ctx, &header, body)
			require.NoError(t, err)

			resMsg, err := must.NotFail(resBody.(*wire.OpMsg).RawDocument()).Decode()
			require.NoError(t, err)

			cursor, err := resMsg.Get("cursor").(bson.RawDocument).Decode()
			require.NoError(t, err)

			nextBatch := cursor.Get("nextBatch").(bson.RawArray)
			newCursorID := cursor.Get("id").(int64)

			expectedDocs := must.NotFail(bson.NewArray(
				must.NotFail(bson.NewDocument("_id", int32(i+1), "v", int32(i+2))),
			))

			testutil.AssertEqual(t, must.NotFail(expectedDocs.Convert()), must.NotFail(nextBatch.Convert()))
			assert.Equal(t, cursorID, newCursorID)
		}
	})

	t.Run("GetMoreEmpty", func(t *testing.T) {
		body, err := wire.NewOpMsg(must.NotFail(getMoreCmd.Encode()))
		require.NoError(t, err)

		// TODO verify header
		header := wire.MsgHeader{
			RequestID: 14,
			OpCode:    wire.OpCodeMsg,
		}

		_, resBody, err := c.Request(ctx, &header, body)
		require.NoError(t, err)

		resMsg, err := must.NotFail(resBody.(*wire.OpMsg).RawDocument()).Decode()
		require.NoError(t, err)

		cursor, err := resMsg.Get("cursor").(bson.RawDocument).Decode()
		require.NoError(t, err)

		nextBatch := cursor.Get("nextBatch").(bson.RawArray)
		newCursorID := cursor.Get("id").(int64)

		expectedDocs := must.NotFail(bson.NewArray())

		testutil.AssertEqual(t, must.NotFail(expectedDocs.Convert()), must.NotFail(nextBatch.Convert()))
		assert.Zero(t, newCursorID)
	})
}

// TODO: we don't need this function if both checks are made within Request function
func assertEqualHeader(t testing.TB, expected, actual wire.MsgHeader) (ok bool) {
	t.Helper()

	ok = ok || assert.Equal(t, expected.ResponseTo, actual.ResponseTo)
	ok = ok || assert.Equal(t, expected.MessageLength, actual.MessageLength)

	return
}

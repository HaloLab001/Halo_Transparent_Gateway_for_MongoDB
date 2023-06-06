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

	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/FerretDB/FerretDB/integration/setup"
	"github.com/FerretDB/FerretDB/integration/shareddata"
)

func TestIndexesDropInvalidCollection(t *testing.T) {
	setup.SkipForTigrisWithReason(t, "Indexes are not supported for Tigris")

	t.Parallel()

	for name, tc := range map[string]struct {
		collectionName any
		indexName      any
		err            *mongo.CommandError
		altMessage     string
	}{
		"NonExistentCollection": {
			collectionName: "non-existent",
			indexName:      "index",
			err: &mongo.CommandError{
				Code:    26,
				Name:    "NamespaceNotFound",
				Message: "ns not found TestIndexesDropInvalidCollection-NonExistentCollection.non-existent",
			},
		},
		"InvalidTypeCollection": {
			collectionName: 42,
			indexName:      "index",
			err: &mongo.CommandError{
				Code:    2,
				Name:    "BadValue",
				Message: "collection name has invalid type int",
			},
		},
		"NilCollection": {
			collectionName: nil,
			indexName:      "index",
			err: &mongo.CommandError{
				Code:    2,
				Name:    "BadValue",
				Message: "collection name has invalid type null",
			},
		},
		"EmptyCollection": {
			collectionName: "",
			indexName:      "index",
			err: &mongo.CommandError{
				Code:    73,
				Name:    "InvalidNamespace",
				Message: "Invalid namespace specified 'TestIndexesDropInvalidCollection-EmptyCollection.'",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			provider := shareddata.ArrayDocuments // one provider is enough to check for errors
			ctx, collection := setup.Setup(t, provider)

			command := bson.D{
				{"dropIndexes", tc.collectionName},
				{"index", tc.indexName},
			}

			var res bson.D
			err := collection.Database().RunCommand(ctx, command).Decode(&res)

			require.Nil(t, res)
			AssertEqualAltCommandError(t, *tc.err, tc.altMessage, err)
		})
	}
}

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
	"fmt"
	"testing"

	"github.com/FerretDB/FerretDB/integration/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type insertTestCase struct {
	insert bson.D

	skip          string // skips test if non-empty
	skipForTigris string // skips test for Tigris if non-empty
}

// testInsertCompat tests insert compatibility test cases.
func testInsertCompat(t *testing.T, testCases map[string]insertTestCase) {
	t.Helper()

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Helper()

			if tc.skip != "" {
				t.Skip(tc.skip)
			}
			if tc.skipForTigris != "" {
				setup.SkipForTigrisWithReason(t, tc.skipForTigris)
			}

			t.Parallel()

			ctx, targetCollections, compatCollections := setup.SetupCompat(t)

			insert := tc.insert
			require.NotNil(t, insert)

			var nonEmptyResults bool
			for i := range targetCollections {
				targetCollection := targetCollections[i]
				compatCollection := compatCollections[i]
				t.Run(targetCollection.Name(), func(t *testing.T) {
					t.Helper()

					allDocs := FindAll(t, ctx, targetCollection)

					for _, doc := range allDocs {
						id, ok := doc.Map()["_id"]
						require.True(t, ok)

						t.Run(fmt.Sprint(id), func(t *testing.T) {
							t.Helper()

							filter := bson.D{{"_id", id}}
							var targetInsertRes, compatInsertRes *mongo.InsertOneResult
							var targetErr, compatErr error

							targetInsertRes, targetErr = targetCollection.InsertOne(ctx, insert)
							compatInsertRes, compatErr = compatCollection.InsertOne(ctx, insert)

							if targetErr != nil {
								t.Logf("Target error: %v", targetErr)
								targetErr = UnsetRaw(t, targetErr)
								compatErr = UnsetRaw(t, compatErr)

								// Skip inserts that could not be performed due to Tigris schema validation.
								if e, ok := targetErr.(mongo.CommandError); ok && e.Name == "DocumentValidationFailure" {
									if e.HasErrorCodeWithMessage(121, "json schema validation failed for field") {
										setup.SkipForTigrisWithReason(t, targetErr.Error())
									}
								}

								assert.Equal(t, compatErr, targetErr)
							} else {
								require.NoError(t, compatErr, "compat error; target returned no error")
							}

						})
					}
				})
			}
		})
	}
}

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

// Use _test package to avoid import cycle with testutil.
package pgdb_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/net/context"

	"github.com/FerretDB/FerretDB/internal/handlers/pg/pgdb"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/util/testutil"
)

func TestQueryDocuments(t *testing.T) {
	t.Parallel()

	ctx := testutil.Ctx(t)

	pool := testutil.Pool(ctx, t, nil, zaptest.NewLogger(t))
	dbName := testutil.Schema(ctx, t, pool)
	collectionName := testutil.Table(ctx, t, pool, dbName)

	cases := []struct {
		name             string
		collection       string
		documents        []*types.Document
		docsPerIteration []int // how many documents should be fetched per each iteration
		// len(docsPerIteration) is the amount of fetch iterations.
	}{
		{
			name:             "empty",
			collection:       collectionName,
			documents:        []*types.Document{},
			docsPerIteration: []int{},
		},
		{
			name:             "one",
			collection:       collectionName + "_one",
			documents:        []*types.Document{must.NotFail(types.NewDocument("id", "1"))},
			docsPerIteration: []int{1},
		},
		{
			name:       "two",
			collection: collectionName + "_two",
			documents: []*types.Document{
				must.NotFail(types.NewDocument("id", "1")),
				must.NotFail(types.NewDocument("id", "2")),
			},
			docsPerIteration: []int{2},
		},
		{
			name:       "three",
			collection: collectionName + "_three",
			documents: []*types.Document{
				must.NotFail(types.NewDocument("id", "1")),
				must.NotFail(types.NewDocument("id", "2")),
				must.NotFail(types.NewDocument("id", "3")),
			},
			docsPerIteration: []int{2, 1},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			for _, doc := range tc.documents {
				require.NoError(t, pool.InsertDocument(ctx, dbName, tc.collection, doc))
			}

			fetchedChan, err := pool.QueryDocuments(context.Background(), dbName, tc.collection, "")
			require.NoError(t, err)

			iter := 0
			for {
				fetched, ok := <-fetchedChan
				if !ok {
					break
				}

				require.NoError(t, fetched.Err)
				require.Len(t, fetched.Docs, tc.docsPerIteration[iter], "iteration %d", iter)
				iter++
			}
			require.Equal(t, len(tc.docsPerIteration), iter)
		})
	}

	// Special case: cancel context before reading from channel.
	t.Run("cancel_context", func(t *testing.T) {
		for i := 1; i <= pgdb.FetchedChannelBufSize*pgdb.FetchedSliceCapacity+1; i++ {
			require.NoError(t, pool.InsertDocument(
				ctx, dbName, collectionName+"_cancel",
				must.NotFail(types.NewDocument("id", fmt.Sprintf("%d", i))),
			))
		}

		ctx, cancel := context.WithCancel(context.Background())
		fetchedChan, err := pool.QueryDocuments(ctx, dbName, collectionName+"_cancel", "")
		cancel()
		require.NoError(t, err)

		<-ctx.Done()
		countDocs := 0
		for {
			fetched, ok := <-fetchedChan
			if !ok {
				break
			}

			countDocs += len(fetched.Docs)
		}
		require.Less(t, countDocs, pgdb.FetchedChannelBufSize*pgdb.FetchedSliceCapacity+1)
	})

	// Special case: querying a non-existing collection.
	t.Run("non-existing_collection", func(t *testing.T) {
		fetchedChan, err := pool.QueryDocuments(context.Background(), dbName, collectionName+"_non-existing", "")
		require.NoError(t, err)
		doc, ok := <-fetchedChan
		require.False(t, ok)
		require.Nil(t, doc.Docs)
		require.Nil(t, doc.Err)
	})
}

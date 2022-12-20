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

package pgdb

import (
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FerretDB/FerretDB/internal/util/testutil"
)

func TestSettings(t *testing.T) {
	t.Parallel()

	ctx := testutil.Ctx(t)

	pool := getPool(ctx, t)
	databaseName := testutil.DatabaseName(t)
	collectionName := testutil.CollectionName(t)
	setupDatabase(ctx, t, pool, databaseName)

	err := pool.InTransaction(ctx, func(tx pgx.Tx) error {
		created, err := addSettingsIfNotExists(ctx, tx, databaseName, collectionName)
		require.NoError(t, err)

		var found string

		found, err = getSettings(ctx, tx, databaseName, collectionName)
		require.NoError(t, err)

		assert.Equal(t, created, found)

		// adding settings that already exist should not fail
		_, err = addSettingsIfNotExists(ctx, tx, databaseName, collectionName)
		require.NoError(t, err)

		err = removeSettings(ctx, tx, databaseName, collectionName)
		require.NoError(t, err)

		return nil
	})
	require.NoError(t, err)
}

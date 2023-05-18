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

package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/FerretDB/FerretDB/internal/backends"
	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/handlers/sjson"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
)

// collection implements backends.Collection interface.
type collection struct {
	db   *database
	name string
}

// newDatabase creates a new Collection.
func newCollection(db *database, name string) backends.Collection {
	return backends.CollectionContract(&collection{
		db:   db,
		name: name,
	})
}

// Insert implements backends.Collection interface.
func (c *collection) Insert(ctx context.Context, params *backends.InsertParams) (*backends.InsertResult, error) {
	err := c.db.CreateCollection(ctx, &backends.CreateCollectionParams{Name: c.name})
	if err != nil {
		return nil, err
	}

	conn, err := c.db.b.pool.DB(c.db.name)
	if err != nil {
		return nil, err
	}

	tableName, err := c.db.b.metadataStorage.collectionInfo(ctx, c.db.name, c.name)
	if err != nil {
		return nil, err
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// TODO: check error
	defer tx.Rollback()

	var inserted int64

	for {
		_, doc, err := params.Docs.Next()
		if errors.Is(err, iterator.ErrIteratorDone) {
			break
		}

		if err != nil {
			return nil, err
		}

		query := fmt.Sprintf(`INSERT INTO %s (sjson) VALUES (?)`, tableName)

		bytes, err := sjson.Marshal(doc)
		if err != nil {
			return nil, err
		}

		_, err = tx.ExecContext(ctx, query, bytes)
		if err != nil {
			return nil, err
		}

		inserted++
	}

	params.Docs.Close()

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &backends.InsertResult{
		InsertedCount: inserted,
		Errors:        &commonerrors.WriteErrors{},
	}, nil
}

// check interfaces
var (
	_ backends.Collection = (*collection)(nil)
)

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

package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/handlers/pg/pgdb"
	"github.com/FerretDB/FerretDB/internal/handlers/pg/pjson"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgDropIndexes implements HandlerInterface.
func (h *Handler) MsgDropIndexes(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	dbPool, err := h.DBPool(ctx)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	common.Ignored(document, h.L, "writeConcern", "comment")

	command := document.Command()

	db, err := common.GetRequiredParam[string](document, "$db")
	if err != nil {
		return nil, err
	}

	collection, err := common.GetRequiredParam[string](document, command)
	if err != nil {
		return nil, err
	}

	err = dbPool.InTransactionRetry(ctx, func(tx pgx.Tx) error {
		return processIndexDrop(ctx, tx, db, collection, document, command)
	})

	switch {
	case err == nil:
		// nothing
	case errors.Is(err, pgdb.ErrTableNotExist):
		return nil, commonerrors.NewCommandErrorMsgWithArgument(
			commonerrors.ErrNamespaceNotFound,
			fmt.Sprintf("ns not found %s.%s", db, collection),
			command,
		)
	case errors.Is(err, pgdb.ErrIndexNotExist):
		return nil, commonerrors.NewCommandErrorMsgWithArgument(
			commonerrors.ErrNamespaceNotFound,
			"index not found",
			command,
		)
	case errors.Is(err, pgdb.ErrIndexCannotDelete):
		return nil, commonerrors.NewCommandErrorMsgWithArgument(
			commonerrors.ErrInvalidOptions,
			"cannot drop _id index",
			command,
		)
	default:
		return nil, lazyerrors.Error(err)
	}

	var reply wire.OpMsg
	must.NoError(reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{must.NotFail(types.NewDocument(
			"ok", float64(1),
		))},
	}))

	return &reply, nil
}

// processIndexDrop parses index doc and processes index deletion based on the provided params.
func processIndexDrop(ctx context.Context, tx pgx.Tx, db, collection string, doc *types.Document, command string) error {
	v, err := doc.Get("index")
	if err != nil {
		return lazyerrors.Error(err)
	}

	switch v := v.(type) {
	case *types.Document:
		// Index specification (key) is provided to drop a specific index.
		// TODO: call processIndexKey() from createIndex:
		// indexKey, err := processIndexKey(v)
		// if err != nil {
		// 	return lazyerrors.Error(err)
		// }
		//
		// return pgdb.DropIndex(ctx, tx, db, collection, &pgdb.Index{Key: indexKey})
	case *types.Array:
		// List of index names is provided to drop multiple indexes.
		for {
			iter := v.Iterator()

			defer iter.Close() // It's safe to defer here as the iterators reads everything.

			_, val, err := iter.Next()

			switch {
			case err == nil:
				// nothing
			case errors.Is(err, iterator.ErrIteratorDone):
				return nil
			default:
				return lazyerrors.Error(err)
			}

			index, ok := val.(string)
			if !ok {
				return commonerrors.NewCommandErrorMsgWithArgument(
					commonerrors.ErrTypeMismatch,
					fmt.Sprintf(
						"BSON field 'dropIndexes.index' is the wrong type '%s', expected types '[string, object']",
						pjson.GetTypeOfValue(v),
					),
					command,
				)
			}

			if err := pgdb.DropIndex(ctx, tx, db, collection, &pgdb.Index{Name: index}); err != nil {
				return lazyerrors.Error(err)
			}
		}
	case string:
		if v == "*" {
			// Drop all indexes except the _id index.
			return pgdb.DropAllIndexes(ctx, tx, db, collection)
		}

		// Index name is provided to drop a specific index.
		return pgdb.DropIndex(ctx, tx, db, collection, &pgdb.Index{Name: v})
	}

	return commonerrors.NewCommandErrorMsgWithArgument(
		commonerrors.ErrTypeMismatch,
		fmt.Sprintf(
			"BSON field 'dropIndexes.index' is the wrong type '%s', expected types '[string, object']",
			pjson.GetTypeOfValue(v),
		),
		command,
	)
}

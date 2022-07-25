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
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/handlers/pg/pgdb"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgCount implements HandlerInterface.
func (h *Handler) MsgCount(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	cp, err := h.validateCountParams(ctx, document)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	resDocs := make([]*types.Document, 0, 16)
	err = h.pgPool.InTransaction(ctx, func(tx pgx.Tx) error {
		fetchedChan, err := h.pgPool.QueryDocuments(ctx, tx, cp.sp)
		if err != nil {
			return err
		}
		defer func() {
			// Drain the channel to prevent leaking goroutines.
			// TODO Offer a better design instead of channels: https://github.com/FerretDB/FerretDB/issues/898.
			for range fetchedChan {
			}
		}()

		for fetchedItem := range fetchedChan {
			if fetchedItem.Err != nil {
				return fetchedItem.Err
			}

			for _, doc := range fetchedItem.Docs {
				matches, err := common.FilterDocument(doc, cp.filter)
				if err != nil {
					return err
				}

				if !matches {
					continue
				}

				resDocs = append(resDocs, doc)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if resDocs, err = common.LimitDocuments(resDocs, cp.limit); err != nil {
		return nil, err
	}

	var reply wire.OpMsg
	err = reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{must.NotFail(types.NewDocument(
			"n", int32(len(resDocs)),
			"ok", float64(1),
		))},
	})
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return &reply, nil
}

// findParams are validated find query parameters.
type countParams struct {
	filter *types.Document
	limit  int64
	sp     pgdb.SQLParam
}

// validateCountParams validates document and returns countParams.
func (h *Handler) validateCountParams(ctx context.Context, document *types.Document) (*countParams, error) {
	var cp countParams
	unimplementedFields := []string{
		"skip",
		"collation",
	}
	var err error
	if err = common.Unimplemented(document, unimplementedFields...); err != nil {
		return nil, err
	}
	ignoredFields := []string{
		"hint",
		"readConcern",
		"comment",
	}
	common.Ignored(document, h.l, ignoredFields...)

	if cp.filter, err = common.GetOptionalParam(document, "query", cp.filter); err != nil {
		return nil, err
	}

	if l, _ := document.Get("limit"); l != nil {
		if cp.limit, err = common.GetWholeNumberParam(l); err != nil {
			return nil, err
		}
	}

	if cp.sp.DB, err = common.GetRequiredParam[string](document, "$db"); err != nil {
		return nil, err
	}
	collectionParam, err := document.Get(document.Command())
	if err != nil {
		return nil, err
	}
	var ok bool
	if cp.sp.Collection, ok = collectionParam.(string); !ok {
		return nil, common.NewErrorMsg(
			common.ErrBadValue,
			fmt.Sprintf("collection name has invalid type %s", common.AliasFromType(collectionParam)),
		)
	}
	return &cp, nil
}

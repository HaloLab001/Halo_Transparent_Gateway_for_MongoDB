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
	"os"

	"github.com/FerretDB/FerretDB/build/version"
	"github.com/FerretDB/FerretDB/internal/backends"
	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgExplain implements HandlerInterface.
func (h *Handler) MsgExplain(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	params, err := common.GetExplainParams(document, h.L)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	serverInfo := must.NotFail(types.NewDocument(
		"host", hostname,
		"version", version.Get().MongoDBVersion,
		"gitVersion", version.Get().Commit,

		// our extensions
		"ferretdbVersion", version.Get().Version,
	))

	cmd := params.Command
	cmd.Set("$db", params.DB)

	db, err := h.b.Database(params.DB)
	if err != nil {
		if backends.ErrorCodeIs(err, backends.ErrorCodeDatabaseNameIsInvalid) {
			msg := fmt.Sprintf("Invalid namespace specified '%s.%s'", params.DB, params.Collection)
			return nil, commonerrors.NewCommandErrorMsgWithArgument(commonerrors.ErrInvalidNamespace, msg, document.Command())
		}

		return nil, lazyerrors.Error(err)
	}

	coll, err := db.Collection(params.Collection)
	if err != nil {
		if backends.ErrorCodeIs(err, backends.ErrorCodeCollectionNameIsInvalid) {
			msg := fmt.Sprintf("Invalid collection name: %s", params.Collection)
			return nil, commonerrors.NewCommandErrorMsgWithArgument(commonerrors.ErrInvalidNamespace, msg, document.Command())
		}

		return nil, lazyerrors.Error(err)
	}

	var qp backends.ExplainParams
	if !h.DisableFilterPushdown {
		qp.Filter = params.Filter
	}

	if h.EnableSortPushdown {
		iter := params.Sort.Iterator()
		defer iter.Close()

		var key string
		var order types.SortType
		var count int

		for ; ; count++ {
			var k string
			var v any

			k, v, err = iter.Next()
			if err != nil {
				if errors.Is(err, iterator.ErrIteratorDone) {
					break
				}

				return nil, lazyerrors.Error(err)
			}

			order, err = common.GetSortType(k, v)
			if err != nil {
				return nil, err
			}

			key = k
		}

		// Skip sorting if there are more than one sort parameters
		if count == 1 {
			qp.Sort = &backends.SortField{
				Key:        key,
				Descending: order == types.Descending,
			}
		}
	}

	res, err := coll.Explain(ctx, &qp)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	var reply wire.OpMsg
	must.NoError(reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{must.NotFail(types.NewDocument(
			"queryPlanner", res.QueryPlanner,
			"explainVersion", "1",
			"command", cmd,
			"serverInfo", serverInfo,

			// our extensions
			// TODO https://github.com/FerretDB/FerretDB/issues/3235
			"pushdown", res.QueryPushdown,
			"sortingPushdown", res.SortPushdown,
			"limitPushdown", false,

			"ok", float64(1),
		))},
	}))

	return &reply, nil
}

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

package tigris

import (
	"context"

	api "github.com/tigrisdata/tigrisdb-client-go/api/server/v1"
	"google.golang.org/grpc/codes"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgDropDatabase removes the current database.
func (h *Handler) MsgDropDatabase(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	common.Ignored(document, h.l, "writeConcern", "comment")

	var db string
	if db, err = common.GetRequiredParam[string](document, "$db"); err != nil {
		return nil, err
	}

	res := must.NotFail(types.NewDocument())
	err = h.client.conn.DropDatabase(ctx, db)
	if err != nil {
		switch err := err.(type) {
		case *api.TigrisDBError:
			// TODO: database not found DatabaseNotFound error
			// is hidden in codes.InvalidArgument due to same gRPC status codes
			if err.Code == codes.NotFound {
				break
			}
			return nil, lazyerrors.Error(err)
		default:
			return nil, lazyerrors.Error(err)
		}
	}
	res.Set("dropped", db)
	res.Set("ok", float64(1))

	var reply wire.OpMsg
	err = reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{res},
	})
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return &reply, nil
}

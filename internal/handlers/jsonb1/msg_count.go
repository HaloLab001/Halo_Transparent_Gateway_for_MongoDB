// Copyright 2021 Baltoro OÜ.
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

package jsonb1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/MangoDB-io/MangoDB/internal/handlers/common"
	"github.com/MangoDB-io/MangoDB/internal/pg"
	"github.com/MangoDB-io/MangoDB/internal/types"
	"github.com/MangoDB-io/MangoDB/internal/util/lazyerrors"
	"github.com/MangoDB-io/MangoDB/internal/wire"
)

// MsgCount counts the number of documents matching the query conditions.
func (h *storage) MsgCount(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	m := document.Map()
	collection := m["count"].(string)
	db := m["$db"].(string)

	projection, ok := m["projection"].(types.Document)
	if ok && len(projection.Map()) != 0 {
		return nil, common.NewErrorMessage(common.ErrNotImplemented, "MsgFind: projection is not supported")
	}

	// in count query, key of filter valyes if is "query"
	filter, _ := m["query"].(types.Document)
	sort, _ := m["sort"].(types.Document)
	limit, _ := m["limit"].(int32)

	sql := fmt.Sprintf(`SELECT count(_jsonb) FROM %s`, pgx.Identifier{db, collection}.Sanitize())
	var args []interface{}
	var placeholder pg.Placeholder

	whereSQL, args, err := where(filter, &placeholder)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	sql += whereSQL

	sortMap := sort.Map()
	if len(sortMap) > 0 {
		sql += " ORDER BY"

		for i, k := range sort.Keys() {
			if i != 0 {
				sql += ","
			}

			sql += " _jsonb->" + placeholder.Next()
			args = append(args, k)

			order := sortMap[k].(int32)
			if order > 0 {
				sql += " ASC"
			} else {
				sql += " DESC"
			}
		}
	}

	switch {
	case limit == 0:
		// undefined or zero - no limit
	case limit > 0:
		sql += " LIMIT " + placeholder.Next()
		args = append(args, limit)
	default:
		// TODO https://github.com/MangoDB-io/MangoDB/issues/79
		return nil, common.NewErrorMessage(common.ErrNotImplemented, "MsgFind: negative limit values are not supported")
	}

	rows, err := h.pgPool.Query(ctx, sql, args...)
	var count int32
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return nil, lazyerrors.Error(err)
		}
	}
	// in psql, the SELECT * FROM table limit `x` ignores the value of the limit,
	// so, we need this `if` statement to support this kind of query `db.actor.find().limit(10).count()`
	if count > limit && limit != 0 {
		count = limit
	}
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	defer rows.Close()

	var reply wire.OpMsg
	err = reply.SetSections(wire.OpMsgSection{
		Documents: []types.Document{types.MustMakeDocument(
			"n",
			count,
			"ok",
			float64(1),
		)},
	})
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return &reply, nil
}

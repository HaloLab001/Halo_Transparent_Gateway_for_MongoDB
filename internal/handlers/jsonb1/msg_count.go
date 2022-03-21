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

package jsonb1

import (
	"context"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgFindOrCount finds documents in a collection or view and returns a cursor to the selected documents
// or count the number of documents that matches the query filter.
func (s *storage) MsgCount(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	unimplementedFields := []string{
		"skip",
		"returnKey",
		"showRecordId",
		"tailable",
		"oplogReplay",
		"noCursorTimeout",
		"awaitData",
		"allowPartialResults",
		"collation",
		"allowDiskUse",
		"let",
	}
	if err := common.Unimplemented(document, unimplementedFields...); err != nil {
		return nil, lazyerrors.Error(err)
	}
	ignoredFields := []string{
		"hint",
		"batchSize",
		"singleBatch",
		"comment",
		"maxTimeMS",
		"readConcern",
		"max",
		"min",
	}
	common.Ignored(document, s.l, ignoredFields...)

	command := document.Command()

	var db, collection string
	if db, err = common.GetRequiredParam[string](document, "$db"); err != nil {
		return nil, err
	}
	if collection, err = common.GetRequiredParam[string](document, command); err != nil {
		return nil, err
	}

	var filter, sort *types.Document
	var limit int32
	if filter, err = common.GetOptionalParam(document, "query", filter); err != nil {
		return nil, err
	}
	if sort, err = common.GetOptionalParam(document, "sort", sort); err != nil {
		return nil, err
	}
	if limit, err = common.GetOptionalParam(document, "limit", limit); err != nil {
		return nil, err
	}

	fetchedDocs, err := s.fetch(ctx, db, collection)
	if err != nil {
		return nil, err
	}

	resDocs := make([]*types.Document, 0, 16)
	for _, doc := range fetchedDocs {
		matches, err := common.FilterDocument(doc, filter)
		if err != nil {
			return nil, err
		}

		if !matches {
			continue
		}

		resDocs = append(resDocs, doc)
	}

	if resDocs, err = common.LimitDocuments(resDocs, limit); err != nil {
		return nil, err
	}

	var reply wire.OpMsg
	err = reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{types.MustNewDocument(
			"n", int32(len(resDocs)),
			"ok", float64(1),
		)},
	})
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return &reply, nil
}

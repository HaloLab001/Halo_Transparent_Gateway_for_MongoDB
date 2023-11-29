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

package handler

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/FerretDB/FerretDB/internal/backends"
	"github.com/FerretDB/FerretDB/internal/handler/common"
	"github.com/FerretDB/FerretDB/internal/handler/commonerrors"
	"github.com/FerretDB/FerretDB/internal/handler/commonparams"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgUpdate implements `update` command.
func (h *Handler) MsgUpdate(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	params, err := GetUpdateParams(document, h.L)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	// TODO https://github.com/FerretDB/FerretDB/issues/2612
	_ = params.Ordered

	var we *mongo.WriteError

	matched, modified, upserted, err := h.updateDocument(ctx, params)
	if err != nil {
		switch {
		case backends.ErrorCodeIs(err, backends.ErrorCodeInsertDuplicateID):
			// TODO https://github.com/FerretDB/FerretDB/issues/3263
			we = &mongo.WriteError{
				Index:   0,
				Code:    int(commonerrors.ErrDuplicateKeyInsert),
				Message: fmt.Sprintf(`E11000 duplicate key error collection: %s.%s`, params.DB, params.Collection),
			}

		default:
			if we, err = handleValidationError(err); err != nil {
				return nil, lazyerrors.Error(err)
			}
		}
	}

	res := must.NotFail(types.NewDocument(
		"n", matched,
	))

	if we != nil {
		res.Set("writeErrors", must.NotFail(types.NewArray(WriteErrorDocument(we))))
	}

	if upserted.Len() != 0 {
		res.Set("upserted", upserted)
	}

	res.Set("nModified", modified)
	res.Set("ok", float64(1))

	var reply wire.OpMsg
	must.NoError(reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{res},
	}))

	return &reply, nil
}

// updateDocument iterate through all documents in collection and update them.
func (h *Handler) updateDocument(ctx context.Context, params *common.UpdateParams) (int32, int32, *types.Array, error) {
	var matched, modified int32
	var upserted types.Array

	db, err := h.b.Database(params.DB)
	if err != nil {
		if backends.ErrorCodeIs(err, backends.ErrorCodeDatabaseNameIsInvalid) {
			msg := fmt.Sprintf("Invalid namespace specified '%s.%s'", params.DB, params.Collection)
			return 0, 0, nil, commonerrors.NewCommandErrorMsgWithArgument(commonerrors.ErrInvalidNamespace, msg, "update")
		}

		return 0, 0, nil, lazyerrors.Error(err)
	}

	err = db.CreateCollection(ctx, &backends.CreateCollectionParams{Name: params.Collection})

	switch {
	case err == nil:
		// nothing
	case backends.ErrorCodeIs(err, backends.ErrorCodeCollectionAlreadyExists):
		// nothing
	case backends.ErrorCodeIs(err, backends.ErrorCodeCollectionNameIsInvalid):
		msg := fmt.Sprintf("Invalid collection name: %s", params.Collection)
		return 0, 0, nil, commonerrors.NewCommandErrorMsgWithArgument(commonerrors.ErrInvalidNamespace, msg, "insert")
	default:
		return 0, 0, nil, lazyerrors.Error(err)
	}

	for _, u := range params.Updates {
		c, err := db.Collection(params.Collection)
		if err != nil {
			if backends.ErrorCodeIs(err, backends.ErrorCodeCollectionNameIsInvalid) {
				msg := fmt.Sprintf("Invalid collection name: %s", params.Collection)
				return 0, 0, nil, commonerrors.NewCommandErrorMsgWithArgument(commonerrors.ErrInvalidNamespace, msg, "insert")
			}

			return 0, 0, nil, lazyerrors.Error(err)
		}

		var qp backends.QueryParams
		if !h.DisableFilterPushdown {
			qp.Filter = u.Filter
		}

		res, err := c.Query(ctx, &qp)
		if err != nil {
			return 0, 0, nil, lazyerrors.Error(err)
		}

		var resDocs []*types.Document

		defer res.Iter.Close()

		for {
			var doc *types.Document

			_, doc, err = res.Iter.Next()
			if err != nil {
				if errors.Is(err, iterator.ErrIteratorDone) {
					break
				}

				return 0, 0, nil, lazyerrors.Error(err)
			}

			var matches bool

			matches, err = common.FilterDocument(doc, u.Filter)
			if err != nil {
				return 0, 0, nil, lazyerrors.Error(err)
			}

			if !matches {
				continue
			}

			resDocs = append(resDocs, doc)
		}

		res.Iter.Close()

		if len(resDocs) == 0 {
			if !u.Upsert {
				// nothing to do, continue to the next update operation
				continue
			}

			// TODO https://github.com/FerretDB/FerretDB/issues/3040
			hasQueryOperators, err := common.HasQueryOperator(u.Filter)
			if err != nil {
				return 0, 0, nil, lazyerrors.Error(err)
			}

			var doc *types.Document
			if hasQueryOperators {
				doc = must.NotFail(types.NewDocument())
			} else {
				doc = u.Filter
			}

			hasUpdateOperators, err := common.HasSupportedUpdateModifiers("update", u.Update)
			if err != nil {
				return 0, 0, nil, err
			}

			if hasUpdateOperators {
				// TODO https://github.com/FerretDB/FerretDB/issues/3044
				if _, err = common.UpdateDocument("update", doc, u.Update); err != nil {
					return 0, 0, nil, err
				}
			} else {
				doc = u.Update
			}

			if !doc.Has("_id") {
				doc.Set("_id", types.NewObjectID())
			}
			upserted.Append(must.NotFail(types.NewDocument(
				"index", int32(upserted.Len()),
				"_id", must.NotFail(doc.Get("_id")),
			)))

			// TODO https://github.com/FerretDB/FerretDB/issues/3454
			if err = doc.ValidateData(); err != nil {
				return 0, 0, nil, err
			}

			// TODO https://github.com/FerretDB/FerretDB/issues/2612

			_, err = c.InsertAll(ctx, &backends.InsertAllParams{
				Docs: []*types.Document{doc},
			})
			if err != nil {
				return 0, 0, nil, err
			}

			matched++

			continue
		}

		if len(resDocs) > 1 && !u.Multi {
			resDocs = resDocs[:1]
		}

		matched += int32(len(resDocs))

		for _, doc := range resDocs {
			changed, err := common.UpdateDocument("update", doc, u.Update)
			if err != nil {
				return 0, 0, nil, lazyerrors.Error(err)
			}

			if !changed {
				continue
			}

			// TODO https://github.com/FerretDB/FerretDB/issues/3454
			if err = doc.ValidateData(); err != nil {
				return 0, 0, nil, err
			}

			updateRes, err := c.UpdateAll(ctx, &backends.UpdateAllParams{Docs: []*types.Document{doc}})
			if err != nil {
				return 0, 0, nil, lazyerrors.Error(err)
			}

			modified += int32(updateRes.Updated)
		}
	}

	return matched, modified, &upserted, nil
}

// UpdateParams represents parameters for the update command.
//
//nolint:vet // for readability
type UpdateParams struct {
	DB         string `ferretdb:"$db"`
	Collection string `ferretdb:"update,collection"`

	Updates []Update `ferretdb:"updates"`

	Comment string `ferretdb:"comment,opt"`

	Let *types.Document `ferretdb:"let,unimplemented"`

	Ordered                  bool            `ferretdb:"ordered,ignored"`
	BypassDocumentValidation bool            `ferretdb:"bypassDocumentValidation,ignored"`
	WriteConcern             *types.Document `ferretdb:"writeConcern,ignored"`
	LSID                     any             `ferretdb:"lsid,ignored"`
}

// Update represents a single update operation parameters.
//
//nolint:vet // for readability
type Update struct {
	Filter *types.Document `ferretdb:"q,opt"`
	Update *types.Document `ferretdb:"u,opt"` // TODO https://github.com/FerretDB/FerretDB/issues/2742
	Multi  bool            `ferretdb:"multi,opt"`
	Upsert bool            `ferretdb:"upsert,opt,numericBool"`

	C            *types.Document `ferretdb:"c,unimplemented"`
	Collation    *types.Document `ferretdb:"collation,unimplemented"`
	ArrayFilters *types.Array    `ferretdb:"arrayFilters,unimplemented"`

	Hint string `ferretdb:"hint,ignored"`
}

// GetUpdateParams returns parameters for update command.
func GetUpdateParams(document *types.Document, l *zap.Logger) (*UpdateParams, error) {
	var params UpdateParams

	err := commonparams.ExtractParams(document, "update", &params, l)
	if err != nil {
		return nil, err
	}

	if len(params.Updates) > 0 {
		for _, update := range params.Updates {
			if update.Update == nil {
				continue
			}

			if err := common.ValidateUpdateOperators(document.Command(), update.Update); err != nil {
				return nil, err
			}
		}
	}

	return &params, nil
}

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

package common

import (
	"errors"

	"github.com/FerretDB/FerretDB/internal/handlers/common/aggregations/operators"
	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

// AddFieldsIterator returns an iterator that adds a new field to the underlying iterator.
// It will be added to the given closer.
//
// Next method returns the next document after adding the new field to the document.
//
// Close method closes the underlying iterator.
func AddFieldsIterator(iter types.DocumentsIterator, closer *iterator.MultiCloser, newField *types.Document) types.DocumentsIterator { //nolint:lll // for readability
	res := &addFieldsIterator{
		iter:     iter,
		newField: newField,
	}
	closer.Add(res)

	return res
}

// addFieldsIterator is returned by AddFieldsIterator.
type addFieldsIterator struct {
	iter     types.DocumentsIterator
	newField *types.Document
}

// Next implements iterator.Interface. See addFieldsIterator for details.
func (iter *addFieldsIterator) Next() (struct{}, *types.Document, error) {
	var unused struct{}

	_, doc, err := iter.iter.Next()
	if err != nil {
		return unused, nil, lazyerrors.Error(err)
	}

	for _, key := range iter.newField.Keys() {
		val := must.NotFail(iter.newField.Get(key))
		switch v := val.(type) {
		case *types.Document:
			op, err := operators.NewOperator(v)
			if err != nil {
				var opErr operators.OperatorError

				if !errors.As(err, &opErr) {
					return unused, nil, err
				}

				if opErr.Code() == operators.ErrNoOperator {
					break
				}

				switch opErr.Code() {
				case operators.ErrTooManyFields:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrAddFieldsExpressionWrongAmountOfArgs,
						"Invalid $addFields :: caused by :: FieldPath field names may not start with '$'."+
							" Consider using $getField or $setField.",
						"$addFields (stage)",
					)
				case operators.ErrNotImplemented:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrNotImplemented,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrArgsInvalidLen:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrOperatorWrongLenOfArgs,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrInvalidExpression:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrAggregateInvalidExpression,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrInvalidNestedExpression:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrInvalidPipelineOperator,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrWrongType:
					fallthrough
				default:
					return unused, nil, lazyerrors.Error(err)
				}
			}

			val, err = op.Process(doc)
			if err != nil {
				var opErr operators.OperatorError

				if !errors.As(err, &opErr) {
					return unused, nil, err
				}

				if opErr.Code() == operators.ErrNoOperator {
					break
				}

				switch opErr.Code() {
				case operators.ErrTooManyFields:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrFieldPathInvalidName,
						"Invalid $addFields :: caused by :: FieldPath field names may not start with '$'."+
							" Consider using $getField or $setField.",
						"$addFields (stage)",
					)
				case operators.ErrNotImplemented:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrNotImplemented,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrArgsInvalidLen:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrOperatorWrongLenOfArgs,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrInvalidExpression:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrAggregateInvalidExpression,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrInvalidNestedExpression:
					return unused, nil, commonerrors.NewCommandErrorMsgWithArgument(
						commonerrors.ErrInvalidPipelineOperator,
						"Invalid $addFields :: caused by :: "+opErr.Error(),
						"$addFields (stage)",
					)
				case operators.ErrWrongType:
					fallthrough
				default:
					return unused, nil, lazyerrors.Error(err)
				}
			}
		}

		doc.Set(key, val)
	}

	return unused, doc, nil
}

// Close implements iterator.Interface. See AddFieldsIterator for details.
func (iter *addFieldsIterator) Close() {
	iter.iter.Close()
}

// check interfaces
var (
	_ types.DocumentsIterator = (*addFieldsIterator)(nil)
)

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

package aggregations

import (
	"fmt"
	"strings"

	"github.com/FerretDB/FerretDB/internal/handlers/commonpath"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

//go:generate ../../../../bin/stringer -linecomment -type ExpressionErrorCode

// ExpressionErrorCode represents Expression error code.
type ExpressionErrorCode int

const (
	_ ExpressionErrorCode = iota

	// ErrNotExpression indicates that field is not an expression.
	ErrNotExpression

	// ErrInvalidExpression indicates that expression is invalid.
	ErrInvalidExpression

	// ErrEmptyFieldPath indicates that field path expression is empty.
	ErrEmptyFieldPath

	// ErrUndefinedVariable indicates that variable name is not defined.
	ErrUndefinedVariable

	// ErrEmptyVariable indicates that variable name is empty.
	ErrEmptyVariable
)

// ExpressionError describes an error that occurs while evaluating expression.
type ExpressionError struct {
	code ExpressionErrorCode
}

// newExpressionError creates a new ExpressionError.
func newExpressionError(code ExpressionErrorCode) error {
	return &ExpressionError{code: code}
}

// Error implements the error interface.
func (e *ExpressionError) Error() string {
	return e.code.String()
}

// Code returns the ExpressionError code.
func (e *ExpressionError) Code() ExpressionErrorCode {
	return e.code
}

// Expression represents a value that needs evaluation, the output varies depending on the evaluation.
//
// Expression used to access fields in document should be prefixed with a dollar sign `$` followed by field key.
// For accessing embedded document or array, dollar sign `$` should be followed by dot notation.
// Options can be provided to specify how to access fields in embedded array.
//
// Expression for variables should be prefixed with double dollar sign `$$` followed by variable name.
type Expression struct {
	opts commonpath.FindValuesOpts
	path types.Path
}

// NewExpression returns Expression from dollar sign `$` prefixed string.
// It can take additional options to specify how to access fields in embedded array.
//
// It returns error if invalid Expression is provided.
func NewExpression(expression string, opts *commonpath.FindValuesOpts) (*Expression, error) {
	if opts == nil {
		opts = &commonpath.FindValuesOpts{
			FindArrayIndex:     false,
			FindArrayDocuments: true,
		}
	}

	var val string

	switch {
	case strings.HasPrefix(expression, "$$"):
		// `$$` indicates field is a variable.
		v := strings.TrimPrefix(expression, "$$")
		if v == "" {
			return nil, newExpressionError(ErrEmptyVariable)
		}

		if strings.HasPrefix(v, "$") {
			return nil, newExpressionError(ErrInvalidExpression)
		}

		// TODO https://github.com/FerretDB/FerretDB/issues/2275
		return nil, newExpressionError(ErrUndefinedVariable)
	case strings.HasPrefix(expression, "$"):
		// `$` indicates field is a path.
		val = strings.TrimPrefix(expression, "$")

		if val == "" {
			return nil, newExpressionError(ErrEmptyFieldPath)
		}
	default:
		return nil, newExpressionError(ErrNotExpression)
	}

	var err error

	path, err := types.NewPathFromString(val)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return &Expression{
		path: path,
		opts: *opts,
	}, nil
}

// Evaluate returns evaluated expression of document.
// It returns error if evaluation found no value, null and empty array are valid values.
func (e *Expression) Evaluate(doc *types.Document) (any, error) {
	path := e.path

	if path.Len() == 1 {
		val, err := doc.Get(path.String())
		if err != nil {
			return nil, err
		}

		return val, nil
	}

	var isPrefixArray bool
	prefix := path.Prefix()

	if v, err := doc.Get(prefix); err == nil {
		if _, isArray := v.(*types.Array); isArray {
			isPrefixArray = true
		}
	}

	vals, err := commonpath.FindValues(doc, path, &e.opts)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	if len(vals) == 0 {
		if isPrefixArray {
			// when the prefix is array, return empty array.
			return must.NotFail(types.NewArray()), nil
		}

		return nil, fmt.Errorf("no document found under %s path", path)
	}

	if len(vals) == 1 && !isPrefixArray {
		// when the prefix is not array, return the value
		return vals[0], nil
	}

	// when the prefix is array, return an array of value.
	arr := types.MakeArray(len(vals))
	for _, v := range vals {
		arr.Append(v)
	}

	return arr, nil
}

// GetExpressionSuffix returns suffix of pathExpression.
func (e *Expression) GetExpressionSuffix() string {
	return e.path.Suffix()
}

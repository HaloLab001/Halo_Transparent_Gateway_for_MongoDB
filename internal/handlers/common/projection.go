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
	"fmt"
	"strings"
	"time"

	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

var errProjectionEmpty = errors.New("projection is empty")

// ValidateProjection check projection document.
// Document fields could be either included or excluded but not both.
// Exception is for the _id field that could be included or excluded.
func ValidateProjection(projection *types.Document) (*types.Document, bool, error) {
	validated := types.MakeDocument(0)

	if projection.Len() == 0 {
		return nil, false, errProjectionEmpty
	}

	var projectionVal *bool

	iter := projection.Iterator()
	defer iter.Close()

	for {
		key, value, err := iter.Next()
		if errors.Is(err, iterator.ErrIteratorDone) {
			break
		}

		if err != nil {
			return nil, false, lazyerrors.Error(err)
		}

		if strings.Contains(key, "$") {
			return nil, false, commonerrors.NewCommandErrorMsg(
				commonerrors.ErrNotImplemented,
				fmt.Sprintf("projection operator $ is not supported in %s", key),
			)
		}

		var result bool

		switch value := value.(type) {
		case *types.Document:
			return nil, false, commonerrors.NewCommandErrorMsg(
				commonerrors.ErrNotImplemented,
				fmt.Sprintf("projection expression %s is not supported", types.FormatAnyValue(value)),
			)
		case *types.Array, string, types.Binary, types.ObjectID,
			time.Time, types.NullType, types.Regex, types.Timestamp: // all this types are treated as new fields value
			result = true

			validated.Set(key, value)
		case float64, int32, int64:
			// projection treats 0 as false and any other value as true
			comparison := types.Compare(value, int32(0))

			if comparison != types.Equal {
				result = true
			}

			// set the value with boolean result to omit type assertion when we will apply projection
			validated.Set(key, result)
		case bool:
			result = value

			// set the value with boolean result to omit type assertion when we will apply projection
			validated.Set(key, result)
		default:
			return nil, false, lazyerrors.Errorf("unsupported operation %s %value (%T)", key, value, value)
		}

		if projection.Len() == 1 && key == "_id" {
			return validated, result, nil
		}

		// if projectionVal is nil we are processing the first field
		if projectionVal == nil {
			if key == "_id" {
				continue
			}

			projectionVal = &result

			continue
		}

		if *projectionVal != result {
			if *projectionVal {
				return nil, false, commonerrors.NewCommandErrorMsgWithArgument(
					commonerrors.ErrProjectionExIn,
					fmt.Sprintf("Cannot do exclusion on field %s in inclusion projection", key),
					"projection",
				)
			}

			return nil, false, commonerrors.NewCommandErrorMsgWithArgument(
				commonerrors.ErrProjectionInEx,
				fmt.Sprintf("Cannot do inclusion on field %s in exclusion projection", key),
				"projection",
			)
		}
	}

	return validated, *projectionVal, nil
}

// ProjectDocument applies projection to the copy of the document.
func ProjectDocument(doc, projection *types.Document, inclusion bool) (*types.Document, error) {
	projected, err := types.NewDocument("_id", must.NotFail(doc.Get("_id")))
	if err != nil {
		return nil, err
	}

	if projection.Has("_id") {
		idValue := must.NotFail(projection.Get("_id"))

		var set bool

		switch idValue := idValue.(type) {
		case *types.Document: // field: { $elemMatch: { field2: value }}
			return nil, commonerrors.NewCommandErrorMsg(
				commonerrors.ErrCommandNotFound,
				fmt.Sprintf("projection %s is not supported",
					types.FormatAnyValue(idValue),
				),
			)

		case *types.Array, string, types.Binary, types.ObjectID,
			time.Time, types.NullType, types.Regex, types.Timestamp: // all this types are treated as new fields value
			projected.Set("_id", idValue)

			set = true
		case bool:
			set = idValue

		default:
			return nil, lazyerrors.Errorf("unsupported operation %s %v (%T)", "_id", idValue, idValue)
		}

		if !set {
			projected.Remove("_id")
		}
	}

	projectedWithoutID, err := projectDocumentWithoutID(doc, projection, inclusion)
	if err != nil {
		return nil, err
	}

	for _, key := range projectedWithoutID.Keys() {
		projected.Set(key, must.NotFail(projectedWithoutID.Get(key)))
	}

	return projected, nil
}

// projectDocumentWithoutID applies projection to the copy of the document and returns projected document.
// It ignores _id field in the projection.
func projectDocumentWithoutID(doc *types.Document, projection *types.Document, inclusion bool) (*types.Document, error) {
	projectionWithoutID := projection.DeepCopy()
	projectionWithoutID.Remove("_id")

	docWithoutID := doc.DeepCopy()
	docWithoutID.Remove("_id")

	projected := types.MakeDocument(0)

	if !inclusion {
		projected = docWithoutID.DeepCopy()
	}

	iter := projectionWithoutID.Iterator()
	defer iter.Close()

	for {
		key, value, err := iter.Next()
		if errors.Is(err, iterator.ErrIteratorDone) {
			break
		}

		if err != nil {
			return nil, lazyerrors.Error(err)
		}

		path, err := types.NewPathFromString(key)
		if err != nil {
			return nil, lazyerrors.Error(err)
		}

		switch value := value.(type) { // found in the projection
		case *types.Document: // field: { $elemMatch: { field2: value }}
			return nil, commonerrors.NewCommandErrorMsg(
				commonerrors.ErrCommandNotFound,
				fmt.Sprintf("projection %s is not supported",
					types.FormatAnyValue(value),
				),
			)

		case *types.Array, string, types.Binary, types.ObjectID,
			time.Time, types.NullType, types.Regex, types.Timestamp: // all these types are treated as new fields value
			projected.Set(key, value)

		case bool: // field: bool
			if inclusion {
				var v any

				// inclusion projection with existing path sets the value.
				v, err = docWithoutID.GetByPath(path)
				if err == nil {
					if err = projected.SetByPath(path, v); err == nil {
						continue
					}
				}

				if _, err = includeProjection(path, docWithoutID, projected); err != nil {
					return nil, err
				}

				continue
			}

			// exclusion projection removes the value from the path.
			if err = excludeProjection(path, projected); err != nil {
				return nil, err
			}
		default:
			return nil, lazyerrors.Errorf("unsupported operation %s %v (%T)", key, value, value)
		}
	}

	return projected, nil
}

var noValueFound = errors.New("no value found")

// includeProjection copies value found at source at path to the projected.
// Inclusion projection with non-existent path creates an empty document
// or an empty array on the path.
func includeProjection(path types.Path, source any, projected *types.Document) (any, error) {
	key := path.Prefix()

	switch val := source.(type) {
	case *types.Document:
		doc := new(types.Document)

		embeddedSource, err := val.Get(key)
		if err != nil {
			// key does not exist, return empty document.
			return doc, nil
		}

		if path.Len() <= 1 {
			// path reached suffix, return the document.
			projected.Set(key, embeddedSource)
			return val, nil
		}

		// recursively set embedded value of the document.
		embedded, err := includeProjection(path.TrimPrefix(), embeddedSource, doc)
		if err != nil && !errors.Is(err, noValueFound) {
			return nil, err
		}

		if err == nil {
			if !projected.Has(key) {
				// only set if projected does not yet have the key.
				// If projected is `{v: {foo: 1}}` and embedded is `{}`,
				// do not overwrite existing projected.
				projected.Set(key, embedded)
			}
		}

		return doc, nil
	case *types.Array:
		iter := val.Iterator()
		defer iter.Close()

		arr := new(types.Array)

		for {
			_, arrElem, err := iter.Next()
			if err != nil {
				if errors.Is(err, iterator.ErrIteratorDone) {
					break
				}

				return nil, lazyerrors.Error(err)
			}

			if _, ok := arrElem.(*types.Document); !ok {
				continue
			}

			doc := new(types.Document)

			embedded, err := includeProjection(path, arrElem, doc)
			if err != nil && !errors.Is(err, noValueFound) {
				return nil, err
			}

			if err == nil {
				arr.Append(embedded)
			}
		}

		projected.Set(key, arr)

		return arr, nil
	default:
		return nil, noValueFound
	}
}

// excludeProjection removes path from projected value.
// When an array is on the path, it checks if the array contains any document
// with the key to remove that document. This is not the case in document.Remove(key).
//
//	Examples: "v.foo" exclusion projection:
//	{v: {foo: 1}                        -> {v: {}}
//	{v: {foo: 1, bar: 1}                -> {v: {bar: 1}}
//	{v: [{foo: 1}, {foo: 2}]}           -> {v: [{}, {}]}
//	{v: [{foo: 1}, {foo: 2}, {bar: 1}]} -> {v: [{}, {}, {bar: 1}]}
func excludeProjection(path types.Path, projected any) error {
	key := path.Prefix()

	switch projectedVal := projected.(type) {
	case *types.Document:
		embeddedSource, err := projectedVal.Get(key)
		if err != nil {
			return nil
		}

		if path.Len() <= 1 {
			// path reached suffix, remove key from the document.
			projectedVal.Remove(key)
			return nil
		}

		// recursively remove embedded value of the document.
		err = excludeProjection(path.TrimPrefix(), embeddedSource)
		if err != nil && !errors.Is(err, noValueFound) {
			return err
		}

		return nil
	case *types.Array:
		for i := projectedVal.Len() - 1; i >= 0; i-- {
			arrElem := must.NotFail(projectedVal.Get(i))

			if _, ok := arrElem.(*types.Document); !ok {
				// not a document, cannot possibly be part of path, do nothing.
				continue
			}

			err := excludeProjection(path, arrElem)

			if errors.Is(err, noValueFound) {
				projectedVal.Remove(i)
				i--
			}

			if err != nil {
				return err
			}
		}

		return nil
	default:
		return noValueFound
	}
}

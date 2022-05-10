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
	"fmt"
	"math"
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

const (
	projectionElemMatch = "$elemMatch"
	projectionSlice     = "$slice"
)

// validateProjectionExpression: projection can be only inclusion or exclusion. Validate and return true if inclusion.
// Exception for the _id field.
func validateProjectionExpression(projection *types.Document) (bool, error) {
	inclusion, _, err := validateExpression(projection, 0, false, false)
	return inclusion, err
}

func validateExpression(projection *types.Document, depth int, inclusion, exclusion bool) (bool, bool, error) {
	var err error
	for _, k := range projection.Keys() {
		if k == "_id" { // _id is a special case and can be both
			continue
		}

		v := must.NotFail(projection.Get(k))
		switch v := v.(type) {
		case *types.Document:

			inclusion, exclusion, err = validateDocProjectionExpression(v, depth, inclusion, exclusion)

		case *types.Array:
			err = NewErrorMsg(ErrElemMatchObjectRequired, "elemMatch: Invalid argument, object required, but got array")
			return false, false, err

		default: // scalar
			if k == projectionElemMatch {
				err = NewError(ErrElemMatchObjectRequired,
					fmt.Errorf("elemMatch: Invalid argument, object required, but got %T", v),
				)
				return false, false, err
			}

			inclusion, exclusion, err = validateScalarProjectionExpression(v, k, inclusion, exclusion)
		}
	}
	return inclusion, exclusion, err
}

func validateScalarProjectionExpression(v any, field string, inclusion, exclusion bool) (bool, bool, error) {
	var err error
	switch v := v.(type) {
	case float64, int32, int64:
		if types.Compare(v, int32(0)) == types.Equal {
			if inclusion {
				err = NewError(ErrElemMatchExclusionInInclusion,
					fmt.Errorf("Cannot do exclusion on field %s in inclusion projection", field),
				)
				return false, false, err
			}
			exclusion = true
		} else {
			if exclusion {
				err = NewError(ErrElemMatchInclusionInExclusion,
					fmt.Errorf("Cannot do inclusion on field %s in exclusion projection", field),
				)
				return false, false, err
			}
			inclusion = true
		}

	case bool:
		if v {
			if exclusion {
				err = NewError(ErrElemMatchInclusionInExclusion,
					fmt.Errorf("Cannot do inclusion on field %s in exclusion projection", field),
				)
				return false, false, err
			}
			inclusion = true
		} else {
			if inclusion {
				err = NewError(ErrElemMatchExclusionInInclusion,
					fmt.Errorf("Cannot do exclusion on field %s in inclusion projection", field),
				)
				return false, false, err
			}
			exclusion = true
		}
	default:
		err = NewError(ErrNotImplemented, fmt.Errorf("%v of (%T) is not supported", v, v))
		return false, false, err
	}
	return inclusion, exclusion, err
}

func validateDocProjectionExpression(v *types.Document, depth int, inclusion, exclusion bool) (bool, bool, error) {
	var err error
	for _, key := range v.Keys() {
		if key == projectionSlice {
			return false, false, nil
		}

		val := must.NotFail(v.Get(key))
		switch val := val.(type) {
		case *types.Document:

			if key == projectionElemMatch && depth >= 1 {
				err = NewErrorMsg(ErrElemMatchNestedField,
					"Cannot use $elemMatch projection on a nested field.",
				)
				return false, false, err
			}
			inclusion, exclusion, err = validateExpression(val, depth+1, inclusion, exclusion)

		case *types.Array:

			if key == projectionElemMatch {
				err = NewErrorMsg(ErrElemMatchObjectRequired, "elemMatch: Invalid argument, object required, but got array")
				return false, false, err
			}

		default:
			switch key {
			case "$eq",
				"$ne",
				"$gt", "$gte",
				"$lt", "$lte":
				inclusion = true

			case "$in":
				switch val.(type) {
				case *types.Array:
					// ok
				default:
					err = NewErrorMsg(ErrBadValue, "$in needs an array")
					return false, false, err
				}
			case "$nin", "$not", projectionSlice:
				exclusion = true

			default: // $mod, etc
				err = NewErrorMsg(ErrNotImplemented, key+" is not supported")
				return inclusion, exclusion, err
			}
		}
	}
	return inclusion, exclusion, err
}

// ProjectDocuments modifies given documents in places according to the given projection.
func ProjectDocuments(docs []*types.Document, projection *types.Document) error {
	if projection.Len() == 0 {
		return nil
	}

	inclusion, err := validateProjectionExpression(projection)
	if err != nil {
		return err
	}

	for i := 0; i < len(docs); i++ {
		err = projectDocument(inclusion, docs[i], projection)
		if err != nil {
			return err
		}
	}
	return nil
}

func projectDocument(inclusion bool, doc *types.Document, projection *types.Document) error {
	projectionMap := projection.Map()

	for fieldLevel1, k1Val := range doc.Map() {
		k1Projection, ok := projectionMap[fieldLevel1]
		if !ok {
			if fieldLevel1 == "_id" { // if _id is not in projection map, do not do anything with it
				continue
			}
			if inclusion { // k1 from doc is absent in projection, remove from doc only if projection type inclusion
				doc.Remove(fieldLevel1)
			}
			continue
		}

		switch k1Projection := k1Projection.(type) { // found in the projection
		case *types.Document: // in projection doc: k1: { k2: value }}, k1Projection == { k2: value }}
			if err := applyDocProjection(fieldLevel1, doc, k1Projection); err != nil {
				return err
			}

		case *types.Array: // in projection doc: { k1: [value1, value2... ], k1Projection = [ value1, value2.. ]
			// it's a switch over elemMatch projection
			return NewError(ErrElemMatchObjectRequired,
				fmt.Errorf("elemMatch: Invalid argument, object required, but got %T", k1Projection),
			)

		case float64, // in projection doc: { k1: k1Projection } where k1Projection is a number
			int32,
			int64:
			if types.Compare(k1Projection, int32(0)) == types.Equal {
				doc.Remove(fieldLevel1)
			}

		case bool: // in projection doc: { k1: k1Projection }
			if !k1Projection {
				doc.Remove(fieldLevel1)
			}

		default:
			return lazyerrors.Errorf("unsupported projection operation %s %v (%T)", fieldLevel1, k1Val, k1Val)
		}
	}
	return nil
}

func applyDocProjection(k1 string, doc *types.Document, k1Projection *types.Document) error {
	var err error

	for _, projectionName := range k1Projection.Keys() {
		switch projectionName {
		case projectionElemMatch:
			conditions := must.NotFail(k1Projection.Get(projectionName)).(*types.Document)

			var found bool
			found, err = findDocElemMatch(k1, doc, conditions)
			if err != nil {
				return err
			}

			if !found {
				doc.Remove(k1)
				return nil
			}

		case projectionSlice:
			var docValue any
			docValue, err = doc.Get(k1)
			if err != nil { // the field can't be obtained, so there is nothing to do
				return err
			}
			// $slice works only for arrays, so docValue must be an array
			arr, ok := docValue.(*types.Array)
			if !ok {
				return err
			}
			projectionVal := must.NotFail(k1Projection.Get(projectionName))
			res, err := filterFieldArraySlice(arr, projectionVal)
			if err != nil {
				return err
			}
			if res == nil {
				must.NoError(doc.Set(k1, types.Null))
				return nil
			}
			must.NoError(doc.Set(k1, res))

		default:
			return NewErrorMsg(ErrCommandNotFound, projectionName+" not supported")
		}
	}
	return err
}

func findInArray(k1, k2 string, value any, doc *types.Document, compareRes []types.CompareResult) bool {
	docValueArray := must.NotFail(doc.GetByPath(k1)).(*types.Array)

	found := -1
fieldArray:
	for j := 0; j < docValueArray.Len(); j++ {
		e, err := docValueArray.Get(j)
		if err != nil {
			continue
		}

		if found >= 0 {
			doc.RemoveByPath(k1, strconv.Itoa(j))
			j -= 1
			continue
		}
		switch e := e.(type) {
		case *types.Document:
			var d any
			d, err = e.Get(k2)
			if err != nil {
				doc.RemoveByPath(k1, strconv.Itoa(j))
				j -= 1
				continue
			}
			switch value := value.(type) {
			// TODO https://github.com/FerretDB/FerretDB/issues/439 resolve in nested
			case *types.Document:

			case *types.Array:
				// value in expression is an array
				// and it is the target array of the doc
				for i := 0; i < e.Len(); i++ {
					arrV := must.NotFail(value.Get(i))
					switch d := d.(type) {
					case *types.Array:
						cmp := types.Compare(d, arrV)
						if slices.Contains(compareRes, cmp) {
							found = j
							continue fieldArray
						}
					default:
						continue
					}
				}

			default:
				cmp := types.Compare(d, value)
				if slices.Contains(compareRes, cmp) {
					found = j
					continue
				}
			}

			doc.RemoveByPath(k1, strconv.Itoa(j))
			j -= 1

		default:
			doc.RemoveByPath(k1, strconv.Itoa(j))
			j -= 1
			continue
		}
	}

	if found < 0 {
		doc.RemoveByPath(k1)
	}
	return found >= 0
}

// findDocElemMatch is for elemMatch conditions.
func findDocElemMatch(k1 string, doc, conditions *types.Document) (bool, error) {
	found := false

	// for sure it's here - see code above
	docValueA := must.NotFail(doc.GetByPath(k1))

	// $elemMatch works only for arrays, it must be an array
	docValueArray, ok := docValueA.(*types.Array)
	if !ok {
		doc.Remove(k1)
		return found, nil
	}

	for k2, condition := range conditions.Map() {
		switch condition := condition.(type) {
		// in condition: { $eq: 42 }
		case *types.Document:
			for operand, value := range condition.Map() {
				var err error
				switch operand {
				case "$eq":
					found = findInArray(k1, k2, value, doc, []types.CompareResult{types.Equal})

				case "$ne":
					found = findInArray(k1, k2, value, doc, []types.CompareResult{types.Less, types.Greater})

				case "$gt":
					found = findInArray(k1, k2, value, doc, []types.CompareResult{types.Greater})

				case "$gte":
					found = findInArray(k1, k2, value, doc, []types.CompareResult{types.Greater, types.Equal})

				case "$lt":
					found = findInArray(k1, k2, value, doc, []types.CompareResult{types.Less})

				case "$lte":
					found = findInArray(k1, k2, value, doc, []types.CompareResult{types.Less, types.Equal})

				case "$nin":
					switch inValue := value.(type) {
					case *types.Array:
						for i := 0; i < inValue.Len(); i++ {
							x := must.NotFail(inValue.Get(i))
							found = findInArray(k1, k2, x, doc,
								[]types.CompareResult{types.Less, types.Greater, types.NotEqual},
							)
							if found {
								return found, err
							}
						}
					default:
						err = NewErrorMsg(ErrBadValue, "$nin needs an array")
						return found, err
					}
					if !found {
						return found, err
					}

				case "$in":
					switch inValue := value.(type) {
					case *types.Array:
						for i := 0; i < inValue.Len(); i++ {
							x := must.NotFail(inValue.Get(i))
							found = findInArray(k1, k2, x, doc, []types.CompareResult{types.Equal})
							if found {
								return found, err
							}
						}
					default:
						err = NewErrorMsg(ErrBadValue, "array values supported for $in only")
						return found, err
					}
					if !found {
						return found, err
					}

					// operand is not an operand possible: <scalar value> OR field: {nested projection}
				default:

					for j := 0; j < docValueArray.Len(); j++ {
						e := must.NotFail(docValueArray.Get(j))

						switch e := e.(type) {
						case *types.Document:
							docVal, err := e.Get(k2)
							if err != nil {
								doc.RemoveByPath(k1, strconv.Itoa(j))
								continue
							}
							if types.Compare(docVal, value) == types.Equal {
								found = true
								break
							}
						default: // field2: value
							if types.Compare(e, value) == types.Equal {
								found = true
								break
							}
						}
						doc.RemoveByPath(k1, strconv.Itoa(j))
						j = j - 1
					}
					err = NewErrorMsg(ErrBadValue, k2+" not supported")
					return found, err
				}
				return found, err
			}
		}
	}
	return found, nil
}

// filterFieldArraySlice implements $slice projection query.
func filterFieldArraySlice(docValue *types.Array, projectionValue any) (*types.Array, error) {
	switch projectionValue := projectionValue.(type) {
	case int32, int64, float64:
		return projectionSliceSingleArg(docValue, projectionValue), nil
	case *types.Array:
		if projectionValue.Len() < 2 || projectionValue.Len() > 3 {
			return nil, NewErrorMsg(
				ErrInvalidArg,
				fmt.Sprintf(
					"Invalid $slice syntax. The given syntax "+
						"did not match the find() syntax because :: Location31272: "+
						"$slice array argument should be of form [skip, limit] :: "+
						"The given syntax did not match the expression "+
						"$slice syntax. :: caused by :: "+
						"Expression $slice takes at least 2 arguments, and at most 3, but %d were passed in.",
					projectionValue.Len(),
				),
			)
		}

		if projectionValue.Len() == 3 {
			// this is the error MongoDB 5.0 is returning in this case
			return nil, NewErrorMsg(
				ErrSliceFirstArg,
				fmt.Sprintf(
					"First argument to $slice must be an array, but is of type: %s",
					AliasFromType(must.NotFail(projectionValue.Get(0))),
				),
			)
		}

		return projectionSliceMultiArgs(docValue, projectionValue)

	default:
		return nil, NewErrorMsg(
			ErrInvalidArg,
			"Invalid $slice syntax. The given syntax "+
				"did not match the find() syntax because :: Location31273: "+
				"$slice only supports numbers and [skip, limit] arrays :: "+
				"The given syntax did not match the expression $slice syntax. :: caused by :: "+
				"Expression $slice takes at least 2 arguments, and at most 3, but 1 were passed in.",
		)
	}
}

func projectionSliceSingleArg(arr *types.Array, arg any) *types.Array {
	var n int
	switch v := arg.(type) {
	case float64:
		if math.IsNaN(v) {
			break // because n == 0 already
		}
		if math.IsInf(v, -1) || v < math.MinInt {
			n = math.MinInt
			break
		}
		if math.IsInf(v, +1) || v > math.MaxInt {
			n = math.MaxInt
			break
		}
		n = int(v)
	case int64:
		if v > math.MaxInt {
			n = math.MaxInt
			break
		}
		if v < math.MinInt {
			n = math.MinInt
			break
		}
		n = int(v)
	case int32:
		n = int(v)
	}

	// negative n is OK in case of a single argument
	var skip, limit int
	if n < 0 {
		skip, limit = arr.Len()+n, arr.Len()
		n = -n
	} else {
		skip, limit = 0, n
	}
	if n < arr.Len() {
		res := types.MakeArray(limit)
		for i := skip; i < limit; i++ {
			must.NoError(res.Append(must.NotFail(arr.Get(i))))
		}
		return res
	}
	// otherwise return arr as is
	return arr
}

func projectionSliceMultiArgs(arr, args *types.Array) (*types.Array, error) {
	var skip, limit int
	pair := [2]int{}
	for i := range pair {
		switch v := must.NotFail(args.Get(i)).(type) {
		case types.NullType:
			return nil, nil //nolint:nilnil // nil is a valid value
		case float64:
			if math.IsNaN(v) {
				break // because pair[i] == 0 already
			}
			if math.IsInf(v, -1) || v < math.MinInt {
				pair[i] = math.MinInt
				break
			}
			if math.IsInf(v, +1) || v > math.MaxInt {
				pair[i] = math.MaxInt
				break
			}
			pair[i] = int(v)
		case int64:
			if v > math.MaxInt {
				pair[i] = math.MaxInt
				break
			}
			if v < math.MinInt {
				pair[i] = math.MinInt
				break
			}
			pair[i] = int(v)
		case int32:
			pair[i] = int(v)
		default:
			return nil, NewErrorMsg(
				ErrSliceFirstArg,
				fmt.Sprintf(
					"First argument to $slice must be an array, but is of type: %s",
					AliasFromType(must.NotFail(args.Get(0))),
				),
			)
		}

		if i == 1 && pair[i] < 0 { // limit can't be negative in case of 2 arguments
			return nil, NewErrorMsg(
				ErrSliceFirstArg,
				fmt.Sprintf(
					"First argument to $slice must be an array, but is of type: %s",
					AliasFromType(must.NotFail(args.Get(0))),
				),
			)
		}
	}

	skip, limit = pair[0], pair[1]

	if skip < 0 {
		if -skip >= arr.Len() {
			skip = 0
		} else {
			skip = arr.Len() + skip
		}
	} else {
		if skip > arr.Len() {
			return types.MakeArray(0), nil
		}
	}
	limit += skip
	if limit >= arr.Len() {
		limit = arr.Len()
	}
	res := types.MakeArray(limit)
	for i := skip; i < limit; i++ {
		must.NoError(res.Append(must.NotFail(arr.Get(i))))
	}
	return res, nil
}

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

package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/FerretDB/FerretDB/internal/util/must"
)

// DocumentPathErrorCode represents DocumentPathError error code.
type DocumentPathErrorCode int

const (
	// ErrDocumentPathKeyNotFound indicates that key was not found in document.
	ErrDocumentPathKeyNotFound = iota + 1
	// ErrDocumentPathCannotAccess indicates that path couldn't be accessed.
	ErrDocumentPathCannotAccess
	// ErrDocumentPathArrayInvalidIndex indicates that provided array index is invalid.
	ErrDocumentPathArrayInvalidIndex
	// ErrDocumentPathIndexOutOfBound indicates that provided array index is out of bound.
	ErrDocumentPathIndexOutOfBound
	// ErrDocumentPathCannotCreateField indicates that it's impossible to create a specific field.
	ErrDocumentPathCannotCreateField
)

// DocumentPathError describes an error that could occur on document path related operations.
type DocumentPathError struct {
	reason error
	code   DocumentPathErrorCode
}

// Error implements the error interface.
func (e *DocumentPathError) Error() string {
	return e.reason.Error()
}

// Code returns the DocumentPathError code.
func (e *DocumentPathError) Code() DocumentPathErrorCode {
	return e.code
}

// newDocumentPathError creates a new DocumentPathError.
func newDocumentPathError(code DocumentPathErrorCode, reason error) error {
	return &DocumentPathError{reason: reason, code: code}
}

// Path represents the field path type. It should be used wherever we work with paths or dot notation.
// Path should be stored and passed as a value. Its methods return new values, not modifying the receiver's state.
type Path struct {
	s []string
}

// NewPath returns Path from a strings slice.
func NewPath(path ...string) Path {
	if len(path) == 0 {
		panic("empty path")
	}
	for _, s := range path {
		if s == "" {
			panic("path element must not be empty")
		}
	}
	p := Path{s: make([]string, len(path))}
	copy(p.s, path)
	return p
}

// NewPathFromString returns Path from path string. Path string should contain fields separated with '.'.
func NewPathFromString(s string) Path {
	path := strings.Split(s, ".")

	return NewPath(path...)
}

// String returns dot-separated path value.
func (p Path) String() string {
	return strings.Join(p.s, ".")
}

// Len returns path length.
func (p Path) Len() int {
	return len(p.s)
}

// Slice returns path values array.
func (p Path) Slice() []string {
	path := make([]string, p.Len())
	copy(path, p.s)
	return path
}

// Suffix returns the last path element.
func (p Path) Suffix() string {
	if len(p.s) <= 1 {
		panic("path should have more than 1 element")
	}
	return p.s[p.Len()-1]
}

// Prefix returns the first path element.
func (p Path) Prefix() string {
	if p.Len() <= 1 {
		panic("path should have more than 1 element")
	}
	return p.s[0]
}

// TrimSuffix returns a path without the last element.
func (p Path) TrimSuffix() Path {
	if p.Len() <= 1 {
		panic("path should have more than 1 element")
	}

	return NewPath(p.s[:p.Len()-1]...)
}

// TrimPrefix returns a copy of path without the first element.
func (p Path) TrimPrefix() Path {
	if p.Len() <= 1 {
		panic("path should have more than 1 element")
	}

	return NewPath(p.s[1:]...)
}

// Append returns new Path constructed from the current path and given element.
func (p Path) Append(elem string) Path {
	elems := p.Slice()

	elems = append(elems, elem)

	return NewPath(elems...)
}

// RemoveByPath removes document by path, doing nothing if the key does not exist.
func RemoveByPath[T CompositeTypeInterface](comp T, path Path) {
	removeByPath(comp, path)
}

// hasByPath checks if the given component contains the given path.
// It returns nil if the path exists and newDocumentPathError otherwise.
// In case of error, its code could be used to determine the exact reason why the path is not found.
//
// It checks the exact path, i.e. to check array element, you need to mention its index in the path.
func hasByPath[T CompositeTypeInterface](comp T, path Path) error {
	var next any = comp
	for _, p := range path.Slice() {
		switch s := next.(type) {
		case *Document:
			var err error
			next, err = s.Get(p)
			if err != nil {
				return newDocumentPathError(ErrDocumentPathKeyNotFound, fmt.Errorf("types.hasByPath: %w", err))
			}

		case *Array:
			index, err := strconv.Atoi(p)
			if err != nil {
				return newDocumentPathError(ErrDocumentPathArrayInvalidIndex, fmt.Errorf("types.hasByPath: %w", err))
			}

			next, err = s.Get(index)
			if err != nil {
				return newDocumentPathError(ErrDocumentPathIndexOutOfBound, fmt.Errorf("types.hasByPath: %w", err))
			}

		default:
			return newDocumentPathError(
				ErrDocumentPathCannotAccess,
				fmt.Errorf("types.hasByPath: can't access %T by path %q", next, p),
			)
		}
	}

	return nil
}

// getAllByPath returns all matching values by path - a sequence of indexes and keys.
func getAllByPath[T CompositeTypeInterface](comp T, path Path, wildcard bool) ([]any, error) {
	next := []any{comp}

	for _, p := range path.Slice() {
		var newNext []any

		for _, n := range next {
			switch s := n.(type) {
			case *Document:
				v, err := s.Get(p)
				if err != nil {
					return nil, newDocumentPathError(ErrDocumentPathKeyNotFound, fmt.Errorf("types.getAllByPath: %w", err))
				}

				newNext = append(newNext, v)

			case *Array:
				index, err := strconv.Atoi(p)

				switch err {
				case nil:
					var v any
					v, err = s.Get(index)

					if err != nil {
						return nil, newDocumentPathError(
							ErrDocumentPathIndexOutOfBound,
							fmt.Errorf("types.getAllByPath: %w", err),
						)
					}
					newNext = append(newNext, v)

				default:
					if !wildcard {
						return nil, newDocumentPathError(
							ErrDocumentPathArrayInvalidIndex,
							fmt.Errorf("types.getAllByPath: %w", err),
						)
					}

					// If p is not an array index, it could be a key of documents inside the array.
					docs := getObjectValuesFromArray(s, p)

					if len(docs) == 0 {
						return nil, newDocumentPathError(
							ErrDocumentPathArrayInvalidIndex,
							fmt.Errorf("types.getAllByPath: %w", err),
						)
					}

					newNext = append(newNext, docs...)
				}

			default:
				return nil, newDocumentPathError(
					ErrDocumentPathCannotAccess,
					fmt.Errorf("types.getAllByPath: can't access %T by path %q", next, p),
				)
			}
		}

		next = newNext
	}

	return next, nil
}

// getObjectValuesFromArray returns a list of documents in the array that have a key with given name.
// If no documents are found, an empty list is returned.
func getObjectValuesFromArray(a *Array, key string) []any {
	var docs []any

	for i := 0; i < a.Len(); i++ {
		v := must.NotFail(a.Get(i))

		if d, ok := v.(*Document); ok {
			if next, err := d.Get(key); err == nil {
				docs = append(docs, next)
			}
		}
	}

	return docs
}

// removeByPath removes path elements for given value, which could be *Document or *Array.
func removeByPath(v any, path Path) {
	if path.Len() == 0 {
		return
	}

	var key string
	if path.Len() == 1 {
		key = path.String()
	} else {
		key = path.Prefix()
	}

	switch v := v.(type) {
	case *Document:
		index := slices.Index(v.Keys(), key)
		if index == -1 {
			return
		}

		if path.Len() == 1 {
			v.Remove(key)
			return
		}

		removeByPath(v.fields[index].value, path.TrimPrefix())

	case *Array:
		i, err := strconv.Atoi(key)
		if err != nil {
			return // no such path
		}
		if i > len(v.s)-1 {
			return // no such path
		}
		if path.Len() == 1 {
			v.s = append(v.s[:i], v.s[i+1:]...)
			return
		}
		removeByPath(v.s[i], path.TrimPrefix())
	default:
		// no such path: scalar value
	}
}

// insertByPath inserts missing parts of the path into Document.
func insertByPath(doc *Document, path Path) error {
	var next any = doc

	var insertedPath Path
	for _, pathElem := range path.TrimSuffix().Slice() {
		insertedPath = insertedPath.Append(pathElem)

		var v any

		if err := doc.HasByPath(insertedPath); err != nil {
			suffix := len(insertedPath.Slice()) - 1
			if suffix < 0 {
				panic("invalid path")
			}

			switch v := next.(type) {
			case *Document:
				v.Set(insertedPath.Slice()[suffix], must.NotFail(NewDocument()))

			case *Array:
				ind, err := strconv.Atoi(insertedPath.Slice()[suffix])
				if err != nil {
					return newDocumentPathError(
						ErrDocumentPathCannotCreateField,
						fmt.Errorf(
							"Cannot create field '%s' in element {%s: %s}",
							pathElem,
							insertedPath.Slice()[suffix-1],
							FormatAnyValue(v),
						),
					)
				}

				// If path needs to be reserved in the middle of the array, we should fill the gap with Null
				for j := v.Len(); j < ind; j++ {
					v.Append(Null)
				}

				v.Append(must.NotFail(NewDocument()))

			default:
				return newDocumentPathError(
					ErrDocumentPathCannotCreateField,
					fmt.Errorf(
						"Cannot create field '%s' in element {%s: %s}",
						pathElem,
						path.Prefix(),
						FormatAnyValue(must.NotFail(doc.Get(path.Prefix()))),
					),
				)
			}

			next = must.NotFail(doc.GetAllByPath(insertedPath, false))[0].(*Document)

			continue
		}

		next = v
	}

	return nil
}

// FormatAnyValue formats value for error message output.
func FormatAnyValue(v any) string {
	switch v := v.(type) {
	case *Document:
		return formatDocument(v)
	case *Array:
		return formatArray(v)
	case float64:
		switch {
		case math.IsNaN(v):
			return "nan.0"

		case math.IsInf(v, -1):
			return "-inf.0"
		case math.IsInf(v, +1):
			return "inf.0"
		case v == 0 && math.Signbit(v):
			return "-0.0"
		case v == 0.0:
			return "0.0"
		case v > 1000 || v < -1000 || v == math.SmallestNonzeroFloat64:
			return fmt.Sprintf("%.15e", v)
		case math.Trunc(v) == v:
			return fmt.Sprintf("%d.0", int64(v))
		default:
			res := fmt.Sprintf("%.2f", v)

			return strings.TrimSuffix(res, "0")
		}

	case string:
		return fmt.Sprintf(`"%v"`, v)
	case Binary:
		return fmt.Sprintf("BinData(%d, %X)", v.Subtype, v.B)
	case ObjectID:
		return fmt.Sprintf("ObjectId('%x')", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case time.Time:
		return fmt.Sprintf("new Date(%d)", v.UnixMilli())
	case NullType:
		return "null"
	case Regex:
		return fmt.Sprintf("/%s/%s", v.Pattern, v.Options)
	case int32:
		return fmt.Sprintf("%d", v)
	case Timestamp:
		return fmt.Sprintf("Timestamp(%v, %v)", int64(v)>>32, int32(v))
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		panic(fmt.Sprintf("unknown type %T", v))
	}
}

// formatDocument formats Document for error output.
func formatDocument(doc *Document) string {
	result := "{ "

	for i, f := range doc.fields {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("%s: %s", f.key, FormatAnyValue(f.value))
	}

	return result + " }"
}

// formatArray formats Array for error output.
func formatArray(array *Array) string {
	if len(array.s) == 0 {
		return "[]"
	}

	result := "[ "

	for _, elem := range array.s {
		result += fmt.Sprintf("%s, ", FormatAnyValue(elem))
	}

	result = strings.TrimSuffix(result, ", ")

	return result + " ]"
}

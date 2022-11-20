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

// Package iterator describes Iterator interface to be used to fetch documents.
package iterator

import "errors"

// ErrIteratorDone  is returned when the iterator is read to the end.
var ErrIteratorDone = errors.New("iterator is read to the end")

// Interface is an iterator interface.
type Interface[E1, E2 any] interface {
	// Next returns a pair of values for containers, such as maps,
	// for which elements have two values.
	// If the iterator is at the end, it returns ErrEndOfIterator.
	Next() (E1, E2, error)
}

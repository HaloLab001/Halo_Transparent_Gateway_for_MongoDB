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

package shareddata

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/maps"
	_ "syreclabs.com/go/faker"
)

// BenchmarkSmallDocuments provides documents that look like:
//
//	{_id: int32(0), id: int32(0), v: "foo"}
//	{_id: int32(1), id: int32(1), v: int32(42)}
//	{_id: int32(2), id: int32(2), v: "42"}
//	{_id: int32(3), id: int32(3), v: {"foo": int32(42)}}
//	...
//
// `_id` is an int32 primary key that starts from 0.
// `id` has the same value as `_id`, but is not indexed by default.
// `v` has one of the four values shown above.
var BenchmarkSmallDocuments = newGeneratorBenchmarkProvider("SmallDocuments", func(docs int) generatorFunc {
	values := []any{
		"foo", int32(42), "42", bson.D{{"foo", int32(42)}},
	}
	l := len(values)

	var i int

	return func() bson.D {
		if i >= docs {
			return nil
		}
		doc := bson.D{
			{"_id", int32(i)},
			{"id", int32(i)},
			{"v", values[i%l]},
		}
		i++

		return doc
	}
})

/*
// BenchmarkLargeDocuments provides a single large document with 123 fields of various types that consists
// of different long keys and simple values.
var BenchmarkLargeDocuments = newGeneratorBenchmarkProvider("LargeDocuments", func(docs int) generatorFunc {
	values := []any{
		true, "foo", false, true, 512, true, 42, false, -42, "42", false, false, true,
		bson.D{{"42", "hello"}},
		primitive.NewDateTimeFromTime(time.Date(2021, 11, 1, 10, 18, 42, 123000000, time.UTC)),
		false, true, 42.13, 501,
	}
	vLen := len(values)

	elements := make([]bson.E, docs)
	elements[0] = bson.E{"_id", 0}

	for i := 1; i < len(elements); i++ {
		elements[i] = bson.E{
			Key:   genLetters(i, 50),
			Value: values[i%vLen],
		}
	}

	doc := bson.D(elements)

	var done bool

	return func() bson.D {
		if done {
			return nil
		}

		done = true

		return doc
	}
})
*/

// AllBenchmarkProviders returns all benchmark providers in random order.
func AllBenchmarkProviders() []BenchmarkProvider {
	providers := []BenchmarkProvider{
		BenchmarkSmallDocuments,
		// FIXME BenchmarkLargeDocuments,
	}

	// check that bse names are unique and randomize order
	res := make(map[string]BenchmarkProvider, len(providers))

	for _, p := range providers {
		n := p.BaseName()
		if _, ok := res[n]; ok {
			panic("duplicate benchmark provider base name: " + n)
		}

		res[n] = p
	}

	return maps.Values(res)
}

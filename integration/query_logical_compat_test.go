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

package integration

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestQueryLogicalCompatAnd(t *testing.T) {
	t.Parallel()

	testCases := map[string]queryCompatTestCase{
		"And": {
			filter: bson.D{{
				"$and", bson.A{
					bson.D{{"value", bson.D{{"$gt", int32(0)}}}},
					bson.D{{"value", bson.D{{"$lt", int64(42)}}}},
				},
			}},
		},
		"BadInput": {
			filter: bson.D{{"$and", nil}},
		},
		"BadExpressionValue": {
			filter: bson.D{{
				"$and", bson.A{
					bson.D{{"value", bson.D{{"$gt", int32(0)}}}},
					nil,
				},
			}},
		},
		"AndOr": {
			filter: bson.D{{
				"$and", bson.A{
					bson.D{{"value", bson.D{{"$gt", int32(0)}}}},
					bson.D{{"$or", bson.A{
						bson.D{{"value", bson.D{{"$lt", int64(42)}}}},
						bson.D{{"value", bson.D{{"$lte", 42.13}}}},
					}}},
				},
			}},
		},
		"AndAnd": {
			filter: bson.D{{
				"$and", bson.A{
					bson.D{{"$and", bson.A{
						bson.D{{"value", bson.D{{"$gt", int32(0)}}}},
						bson.D{{"value", bson.D{{"$lte", 42.13}}}},
					}}},
					bson.D{{"value", bson.D{{"$type", "int"}}}},
				},
			}},
		},
	}

	testQueryCompat(t, testCases)
}

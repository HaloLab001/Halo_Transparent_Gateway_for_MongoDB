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

package tigris

import (
	"context"

	"github.com/tigrisdata/tigris-client-go/driver"

	"github.com/FerretDB/FerretDB/internal/tjson"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// describe gets the collection schema and presents it in built-in map[string]any.
func (h *Handler) describe(ctx context.Context, db, collection string) (map[string]any, error) {
	var res *driver.DescribeCollectionResponse
	var err error
	tigrisDB := h.client.conn.UseDatabase(db)
	res, err = tigrisDB.DescribeCollection(ctx, collection, new(driver.CollectionOptions))
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	return tjson.ParseSchema(res.Schema)
}

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
	"context"
	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

type projectStage struct {
	projection *types.Document
}

func newProject(stage *types.Document) (Stage, error) {
	fields, err := common.GetRequiredParam[*types.Document](stage, "$project")
	if err != nil {
		return nil, err
	}

	return &projectStage{
		projection: fields,
	}, nil
}

func (p *projectStage) Process(_ context.Context, in []*types.Document) ([]*types.Document, error) {
	iter, err := common.ProjectionIterator(iterator.Values(iterator.ForSlice(in)), p.projection)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	res, err := iterator.ConsumeValues(iterator.Interface[struct{}, *types.Document](iter))
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return res, nil
}

func (p *projectStage) Type() StageType {
	return StageTypeDocuments
}

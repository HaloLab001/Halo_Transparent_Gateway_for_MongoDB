// Copyright 2021 Baltoro OÜ.
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
	"github.com/MangoDB-io/MangoDB/internal/pg"
	"github.com/MangoDB-io/MangoDB/internal/types"
	"github.com/MangoDB-io/MangoDB/internal/util/lazyerrors"
)

type wherePair func(key string, value interface{}, p *pg.Placeholder) (sql string, args []interface{}, err error)

//nolint:goconst // $op is fine
func LogicExpr(op string, exprs types.Array, p *pg.Placeholder, wherePair wherePair) (sql string, args []interface{}, err error) {
	if op == "$nor" {
		sql = "NOT ("
	}

	switch op {
	case "$or", "$and", "$nor":
		// {$or: [{expr1}, {expr2}, ...]}
		// {$and: [{expr1}, {expr2}, ...]}
		// {$nor: [{expr1}, {expr2}, ...]}
		for i, expr := range exprs {
			if i != 0 {
				switch op {
				case "$or", "$nor":
					sql += " OR"
				case "$and":
					sql += " AND"
				}
			}

			expr := expr.(types.Document)
			m := expr.Map()
			for j, key := range expr.Keys() {
				if j != 0 {
					sql += " AND"
				}

				var exprSQL string
				var exprArgs []interface{}
				exprSQL, exprArgs, err = wherePair(key, m[key], p)
				if err != nil {
					err = lazyerrors.Errorf("logicExpr: %w", err)
					return
				}

				if sql != "" {
					sql += " "
				}
				sql += "(" + exprSQL + ")"
				args = append(args, exprArgs...)
			}
		}

	default:
		err = lazyerrors.Errorf("logicExpr: unhandled op %q", op)
	}

	if op == "$nor" {
		sql += ")"
	}

	return
}

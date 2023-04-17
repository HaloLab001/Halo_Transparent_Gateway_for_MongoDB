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

package pgdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// DBStats describes statistics for a FerretDB database (PostgreSQL schema).
//
// TODO Include more data https://github.com/FerretDB/FerretDB/issues/2447.
type DBStats struct {
	CountCollections int32
	CountObjects     int32
	CountIndexes     int32
	SizeTotal        int64
	SizeIndexes      int64
	SizeCollections  int64
}

// CollStats describes statistics for a FerretDB collection (PostgreSQL table).
//
// TODO Include more data https://github.com/FerretDB/FerretDB/issues/2447.
type CollStats struct {
	CountObjects   int32
	CountIndexes   int32
	SizeTotal      int64
	SizeIndexes    int64
	SizeCollection int64
}

// CalculateDBStats returns statistics for the given FerretDB database.
//
// If the collection does not exist, it returns an object filled with the given db name and zeros for all the other fields.
func CalculateDBStats(ctx context.Context, tx pgx.Tx, db string) (*DBStats, error) {
	var res DBStats

	// Call ANALYZE to update statistics, the actual statistics are needed to estimate the number of rows in all tables,
	// see https://wiki.postgresql.org/wiki/Count_estimate.
	sql := `ANALYZE`
	if _, err := tx.Exec(ctx, sql); err != nil {
		return nil, err
	}

	// Total size is the disk space used by all the relations in the given schema, including tables, indexes and TOAST data.
	// It also includes the size of FerretDB metadata relations.
	//  See also https://www.postgresql.org/docs/15/functions-admin.html#FUNCTIONS-ADMIN-DBOBJECT
	sql = `
		SELECT 
		    SUM(pg_total_relation_size(schemaname || '.' || tablename)) 
		FROM pg_tables 
		WHERE schemaname = $1`
	args := []any{db}
	row := tx.QueryRow(ctx, sql, args...)

	var schemaSize *int64
	if err := row.Scan(&schemaSize); err != nil {
		return nil, err
	}

	// If the query gave nil, it means the schema does not exist or empty, no need to check other stats.
	if schemaSize == nil {
		return &res, nil
	}

	res.SizeTotal = *schemaSize

	// In this query we select all the tables in the given schema, but we exclude FerretDB metadata table (by reserved prefix).
	sql = `
	SELECT 
	    COUNT(t.tablename)                       AS CountTables,
		COUNT(i.indexname)                       AS CountIndexes,
		COALESCE(SUM(c.reltuples), 0)            AS CountRows,
		COALESCE(SUM(pg_table_size(c.oid)), 0) 	 AS SizeTables,
		COALESCE(SUM(pg_indexes_size(c.oid)), 0) AS SizeIndexes
	FROM pg_tables AS t
		LEFT JOIN pg_class AS c ON c.relname = t.tablename AND c.relnamespace = t.schemaname::regclass
		LEFT JOIN pg_indexes AS i ON i.schemaname = t.schemaname AND i.tablename = t.tablename
	WHERE t.schemaname = $1 AND t.tablename NOT LIKE $2`
	args = []any{db, reservedPrefix + "%"}

	row = tx.QueryRow(ctx, sql, args...)
	if err := row.Scan(&res.CountCollections, &res.CountIndexes, &res.CountObjects, &res.SizeCollections, &res.SizeIndexes); err != nil {
		return nil, err
	}

	return &res, nil
}

// CalculateCollStats returns statistics for the given FerretDB collection in the given database.
//
// If the collection does not exist, it returns an object filled with the given db name and zeros for all the other fields.
func CalculateCollStats(ctx context.Context, tx pgx.Tx, db, collection string) (*CollStats, error) {
	metadata, err := newMetadataStorage(tx, db, collection).get(ctx, false)
	if err != nil {
		return nil, err
	}

	// Call ANALYZE to update statistics, the actual statistics are needed to estimate the number of rows in all tables,
	// see https://wiki.postgresql.org/wiki/Count_estimate.
	sql := `ANALYZE ` + pgx.Identifier{db, metadata.table}.Sanitize()
	if _, err := tx.Exec(ctx, sql); err != nil {
		return nil, lazyerrors.Error(err)
	}

	var res CollStats

	sql = fmt.Sprintf(`
	SELECT
		COALESCE(reltuples, 0)                   AS CountRows,
		COALESCE(pg_total_relation_size(oid), 0) AS SizeTotal,
		COALESCE(pg_table_size(oid), 0) 	     AS SizeTable,
		COALESCE(pg_indexes_size(oid), 0)        AS SizeIndexes
	FROM pg_class 
	WHERE oid = %s::regclass`, quoteString(db+"."+metadata.table))
	row := tx.QueryRow(ctx, sql)
	if err := row.Scan(&res.CountObjects, &res.SizeTotal, &res.SizeCollection, &res.SizeIndexes); err != nil {
		return nil, lazyerrors.Error(err)
	}

	sql = `
	SELECT COUNT(indexname)
	FROM pg_indexes
	WHERE schemaname = $1 AND tablename = $2`
	args := []any{db, metadata.table}
	row = tx.QueryRow(ctx, sql, args...)
	if err := row.Scan(&res.CountIndexes); err != nil {
		return nil, lazyerrors.Error(err)
	}

	return &res, nil
}

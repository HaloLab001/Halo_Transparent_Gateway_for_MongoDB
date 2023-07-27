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

package metadata

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"

	"github.com/FerretDB/FerretDB/internal/backends/sqlite/metadata/pool"
	"github.com/FerretDB/FerretDB/internal/util/fsql"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

const (
	// This prefix is reserved by SQLite for internal use,
	// see https://www.sqlite.org/lang_createtable.html.
	reservedTablePrefix = "sqlite_"

	// SQLite table name where FerretDB metadata is stored.
	metadataTableName = "_ferretdb_collections"
)

const (
	namespace = "ferretdb"
	subsystem = "sqlite_metadata"
)

// Registry provides access to SQLite databases and collections information.
type Registry struct {
	p *pool.Pool
	l *zap.Logger

	rw    sync.RWMutex
	colls map[string]map[string]*Collection // database name -> collection name -> collection
}

// NewRegistry creates a registry for SQLite databases in the directory specified by SQLite URI.
func NewRegistry(u string, l *zap.Logger) (*Registry, error) {
	p, err := pool.New(u, l)
	if err != nil {
		return nil, err
	}

	r := &Registry{
		p:     p,
		l:     l,
		colls: make(map[string]map[string]*Collection),
	}

	for name, db := range p.DBS() {
		r.refreshCollections(context.Background(), name, db)
	}

	return r, nil
}

// Close closes the registry.
func (r *Registry) Close() {
	r.p.Close()
}

func (r *Registry) refreshCollections(ctx context.Context, dbName string, db *fsql.DB) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT name, table_name, settings FROM %q", metadataTableName))
	if err != nil {
		r.l.DPanic("Failed to query metadata table", zap.String("db", dbName), zap.Error(err))
		return
	}
	defer rows.Close()

	colls := make(map[string]*Collection)

	for rows.Next() {
		var c Collection
		if err = rows.Scan(&c.Name, &c.TableName, &c.Settings); err != nil {
			r.l.DPanic("Failed to scan metadata table", zap.String("db", dbName), zap.Error(err))
			return
		}

		colls[c.Name] = &c
	}

	if err = rows.Err(); err != nil {
		r.l.DPanic("Failed to read metadata table", zap.String("db", dbName), zap.Error(err))
	}

	r.rw.Lock()
	r.colls[dbName] = colls
	r.rw.Unlock()
}

func (r *Registry) getCollections(ctx context.Context, dbName string, db *fsql.DB) map[string]*Collection {
	r.rw.RLock()
	colls := r.colls[dbName]
	r.rw.RUnlock()

	// FIXME
	if colls == nil {
		colls = make(map[string]*Collection)
	}

	return colls
}

// DatabaseList returns a sorted list of existing databases.
func (r *Registry) DatabaseList(ctx context.Context) []string {
	return r.p.List(ctx)
}

// DatabaseGetExisting returns a connection to existing database or nil if it doesn't exist.
func (r *Registry) DatabaseGetExisting(ctx context.Context, dbName string) *fsql.DB {
	return r.p.GetExisting(ctx, dbName)
}

// DatabaseGetOrCreate returns a connection to existing database or newly created database.
func (r *Registry) DatabaseGetOrCreate(ctx context.Context, dbName string) (*fsql.DB, error) {
	db, created, err := r.p.GetOrCreate(ctx, dbName)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	if !created {
		return db, nil
	}

	// use transaction
	// TODO https://github.com/FerretDB/FerretDB/issues/2747

	q := fmt.Sprintf("CREATE TABLE %q (name, table_name, settings TEXT)", metadataTableName)
	if _, err = db.ExecContext(ctx, q); err != nil {
		r.DatabaseDrop(ctx, dbName)
		return nil, lazyerrors.Error(err)
	}

	for _, column := range []string{"name", "table_name"} {
		q = fmt.Sprintf("CREATE UNIQUE INDEX %q ON %q (%s)", metadataTableName+"_"+column, metadataTableName, column)
		if _, err = db.ExecContext(ctx, q); err != nil {
			r.DatabaseDrop(ctx, dbName)
			return nil, lazyerrors.Error(err)
		}
	}

	r.refreshCollections(ctx, dbName, db)

	return db, nil
}

// DatabaseDrop drops the database.
//
// Returned boolean value indicates whether the database was dropped.
func (r *Registry) DatabaseDrop(ctx context.Context, dbName string) bool {
	r.rw.Lock()
	delete(r.colls, dbName)
	r.rw.Unlock()

	return r.p.Drop(ctx, dbName)
}

// CollectionList returns a sorted list of collections in the database.
//
// If database does not exist, no error is returned.
func (r *Registry) CollectionList(ctx context.Context, dbName string) ([]string, error) {
	db := r.p.GetExisting(ctx, dbName)
	if db == nil {
		return nil, nil
	}

	colls := r.getCollections(ctx, dbName, db)
	res := maps.Keys(colls)
	sort.Strings(res)
	return res, nil
}

// CollectionCreate creates a collection in the database.
//
// Returned boolean value indicates whether the collection was created.
// If collection already exists, (false, nil) is returned.
func (r *Registry) CollectionCreate(ctx context.Context, dbName string, collectionName string) (bool, error) {
	db, err := r.DatabaseGetOrCreate(ctx, dbName)
	if err != nil {
		return false, lazyerrors.Error(err)
	}

	colls := r.getCollections(ctx, dbName, db)
	if colls[collectionName] != nil {
		return false, nil
	}

	h := fnv.New32a()
	must.NotFail(h.Write([]byte(collectionName)))

	tableName := strings.ToLower(collectionName) + "_" + hex.EncodeToString(h.Sum(nil))
	if strings.HasPrefix(tableName, reservedTablePrefix) {
		tableName = "_" + tableName
	}

	// use transaction
	// TODO https://github.com/FerretDB/FerretDB/issues/2747

	q := fmt.Sprintf("CREATE TABLE %q (%s TEXT)", tableName, DefaultColumn)
	if _, err = db.ExecContext(ctx, q); err != nil {
		var e *sqlite.Error
		if errors.As(err, &e) && e.Code() == sqlitelib.SQLITE_ERROR {
			return false, nil
		}

		return false, lazyerrors.Error(err)
	}

	q = fmt.Sprintf("CREATE UNIQUE INDEX %q ON %q (%s)", tableName+"_id", tableName, IDColumn)
	if _, err = db.ExecContext(ctx, q); err != nil {
		_, _ = db.ExecContext(ctx, fmt.Sprintf("DROP TABLE %q", tableName))
		return false, lazyerrors.Error(err)
	}

	q = fmt.Sprintf("INSERT INTO %q (name, table_name, settings) VALUES (?, ?, '{}')", metadataTableName)
	if _, err = db.ExecContext(ctx, q, collectionName, tableName); err != nil {
		_, _ = db.ExecContext(ctx, fmt.Sprintf("DROP TABLE %q", tableName))
		return false, lazyerrors.Error(err)
	}

	r.refreshCollections(ctx, dbName, db)

	return true, nil
}

// CollectionGet returns collection metadata.
//
// If database or collection does not exist, (nil, nil) is returned.
func (r *Registry) CollectionGet(ctx context.Context, dbName string, collectionName string) (*Collection, error) {
	db := r.p.GetExisting(ctx, dbName)
	if db == nil {
		return nil, nil
	}

	colls := r.getCollections(ctx, dbName, db)
	return colls[collectionName], nil
}

// CollectionDrop drops a collection in the database.
//
// Returned boolean value indicates whether the collection was dropped.
// If database or collection did not exist, (false, nil) is returned.
func (r *Registry) CollectionDrop(ctx context.Context, dbName string, collectionName string) (bool, error) {
	db := r.p.GetExisting(ctx, dbName)
	if db == nil {
		return false, nil
	}

	colls := r.getCollections(ctx, dbName, db)
	if colls[collectionName] == nil {
		return false, nil
	}

	info, err := r.CollectionGet(ctx, dbName, collectionName)
	if err != nil {
		return false, lazyerrors.Error(err)
	}

	if info == nil {
		return false, nil
	}

	// use transaction
	// TODO https://github.com/FerretDB/FerretDB/issues/2747

	query := fmt.Sprintf("DELETE FROM %q WHERE name = ?", metadataTableName)
	if _, err := db.ExecContext(ctx, query, collectionName); err != nil {
		return false, lazyerrors.Error(err)
	}

	query = fmt.Sprintf("DROP TABLE %q", info.TableName)
	if _, err := db.ExecContext(ctx, query); err != nil {
		return false, lazyerrors.Error(err)
	}

	r.refreshCollections(ctx, dbName, db)

	return true, nil
}

// Describe implements prometheus.Collector.
func (r *Registry) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(r, ch)
}

// Collect implements prometheus.Collector.
func (r *Registry) Collect(ch chan<- prometheus.Metric) {
	r.p.Collect(ch)

	r.rw.RLock()
	defer r.rw.RLock()

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "databases"),
			"The current number of database in the registry.",
			nil, nil,
		),
		prometheus.GaugeValue,
		float64(len(r.colls)),
	)

	for db, colls := range r.colls {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "collections"),
				"The current number of collections in the registry.",
				[]string{"db"}, nil,
			),
			prometheus.GaugeValue,
			float64(len(colls)),
			db,
		)
	}
}

// check interfaces
var (
	_ prometheus.Collector = (*Registry)(nil)
)

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
	"math"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/FerretDB/FerretDB/integration/setup"
	"github.com/FerretDB/FerretDB/integration/shareddata"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/util/testutil"
)

func TestCommandsDiagnosticGetLog(t *testing.T) {
	t.Parallel()
	res := setup.SetupWithOpts(t, &setup.SetupOpts{
		DatabaseName: "admin",
	})

	ctx, collection := res.Ctx, res.TargetCollection

	for name, tc := range map[string]struct {
		command  bson.D
		expected map[string]any
		err      *mongo.CommandError
		alt      string
	}{
		"Asterisk": {
			command: bson.D{{"getLog", "*"}},
			expected: map[string]any{
				"names": bson.A(bson.A{"global", "startupWarnings"}),
				"ok":    float64(1),
			},
		},
		"Global": {
			command: bson.D{{"getLog", "global"}},
			expected: map[string]any{
				"totalLinesWritten": int64(1024),
				"log":               bson.A{},
				"ok":                float64(1),
			},
		},
		"StartupWarnings": {
			command: bson.D{{"getLog", "startupWarnings"}},
			expected: map[string]any{
				"totalLinesWritten": int64(1024),
				"log":               bson.A{},
				"ok":                float64(1),
			},
		},
		"NonExistentName": {
			command: bson.D{{"getLog", "nonExistentName"}},
			err: &mongo.CommandError{
				Code:    0,
				Message: `no RamLog named: nonExistentName`,
			},
			alt: `no RecentEntries named: nonExistentName`,
		},
		"Nil": {
			command: bson.D{{"getLog", nil}},
			err: &mongo.CommandError{
				Code:    14,
				Name:    "TypeMismatch",
				Message: `Argument to getLog must be of type String; found null of type null`,
			},
			alt: "Argument to getLog must be of type String",
		},
		"NaN": {
			command: bson.D{{"getLog", math.NaN()}},
			err: &mongo.CommandError{
				Code:    14,
				Name:    "TypeMismatch",
				Message: `Argument to getLog must be of type String; found nan.0 of type double`,
			},
			alt: "Argument to getLog must be of type String",
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var actual bson.D
			err := collection.Database().RunCommand(ctx, tc.command).Decode(&actual)
			if err != nil {
				AssertEqualAltError(t, *tc.err, tc.alt, err)
				return
			}
			require.NoError(t, err)

			m := actual.Map()
			k := CollectKeys(t, actual)

			for key, item := range tc.expected {
				assert.Contains(t, k, key)
				if key != "log" && key != "totalLinesWritten" {
					assert.Equal(t, m[key], item)
				}
			}
		})
	}
}

func TestCommandsDiagnosticHostInfo(t *testing.T) {
	t.Parallel()
	ctx, collection := setup.Setup(t)

	var actual bson.D
	err := collection.Database().RunCommand(ctx, bson.D{{"hostInfo", 42}}).Decode(&actual)
	require.NoError(t, err)

	m := actual.Map()
	t.Log(m)

	assert.Equal(t, float64(1), m["ok"])
	assert.Equal(t, []string{"system", "os", "extra", "ok"}, CollectKeys(t, actual))

	os := m["os"].(bson.D)
	assert.Equal(t, []string{"type", "name", "version"}, CollectKeys(t, os))

	if runtime.GOOS == "linux" {
		require.NotEmpty(t, os.Map()["name"], "os name should not be empty")
		require.NotEmpty(t, os.Map()["version"], "os version should not be empty")
	}

	system := m["system"].(bson.D)
	keys := CollectKeys(t, system)
	assert.Contains(t, keys, "currentTime")
	assert.Contains(t, keys, "hostname")
	assert.Contains(t, keys, "cpuAddrSize")
	assert.Contains(t, keys, "numCores")
	assert.Contains(t, keys, "cpuArch")
}

func TestCommandsDiagnosticListCommands(t *testing.T) {
	t.Parallel()
	ctx, collection := setup.Setup(t)

	var actual bson.D
	err := collection.Database().RunCommand(ctx, bson.D{{"listCommands", 42}}).Decode(&actual)
	require.NoError(t, err)

	m := actual.Map()
	t.Log(m)

	assert.Equal(t, float64(1), m["ok"])
	assert.Equal(t, []string{"commands", "ok"}, CollectKeys(t, actual))

	commands := m["commands"].(bson.D)
	listCommands := commands.Map()["listCommands"].(bson.D)
	assert.NotEmpty(t, listCommands.Map()["help"].(string))
}

func TestCommandsDiagnosticConnectionStatus(t *testing.T) {
	t.Parallel()
	ctx, collection := setup.Setup(t)

	var actual bson.D
	err := collection.Database().RunCommand(ctx, bson.D{{"connectionStatus", "*"}}).Decode(&actual)
	require.NoError(t, err)

	ok := actual.Map()["ok"]

	assert.Equal(t, float64(1), ok)
}

func TestCommandsDiagnosticExplain(t *testing.T) {
	t.Parallel()
	ctx, collection := setup.Setup(t, shareddata.Scalars)
	collectionName := testutil.CollectionName(t)

	for name, tc := range map[string]struct {
		command             bson.D
		expectedCommandKeys []string
		err                 *mongo.CommandError
	}{
		"Count": {
			command: bson.D{
				{
					"explain", bson.D{
						{"count", collectionName},
						{"query", bson.D{{"value", bson.D{{"$type", "array"}}}}},
					},
				},
				{"verbosity", "queryPlanner"},
			},
			expectedCommandKeys: []string{"count", "query", "$db"},
		},
		"Find": {
			command: bson.D{
				{
					"explain", bson.D{
						{"find", collectionName},
						{"filter", bson.D{{"value", bson.D{{"$type", "array"}}}}},
					},
				},
				{"verbosity", "queryPlanner"},
			},
			expectedCommandKeys: []string{"find", "filter", "$db"},
		},
		"FindAndModify": {
			command: bson.D{
				{
					"explain", bson.D{
						{"findAndModify", collectionName},
						{"query", bson.D{{
							"$and",
							bson.A{
								bson.D{{"value", bson.D{{"$gt", int32(0)}}}},
								bson.D{{"value", bson.D{{"$lt", int32(0)}}}},
							},
						}}},
						{"update", bson.D{{"$set", bson.D{{"v", 43.13}}}}},
						{"upsert", true},
					},
				},
				{"verbosity", "queryPlanner"},
			},
			expectedCommandKeys: []string{"findAndModify", "query", "update", "upsert"},
		},
		"FindFieldAbsent": {
			command: bson.D{
				{
					"explain", bson.D{
						{"find", collectionName},
						{"query", bson.D{{"value", bson.D{{"$type", "array"}}}}},
					},
				},
				{"verbosity", "queryPlanner"},
			},
			err: &mongo.CommandError{
				Code:    40415,
				Message: "BSON field 'FindCommandRequest.query' is an unknown field.",
				Name:    "Location40415",
			},
		},
		"FindAndModifyFieldAbsent": {
			command: bson.D{
				{
					"explain", bson.D{
						{"findAndModify", collectionName},
						{"query", bson.D{{
							"$and",
							bson.A{
								bson.D{{"value", bson.D{{"$gt", int32(0)}}}},
								bson.D{{"value", bson.D{{"$lt", int32(0)}}}},
							},
						}}},
						{"upsert", true},
					},
				},
				{"verbosity", "queryPlanner"},
			},
			err: &mongo.CommandError{
				Code:    9,
				Message: "Either an update or remove=true must be specified",
				Name:    "FailedToParse",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var actual bson.D
			err := collection.Database().RunCommand(ctx, tc.command).Decode(&actual)

			if tc.err != nil {
				AssertEqualError(t, *tc.err, err)
				return
			}

			require.NoError(t, err)
			actualD := ConvertDocument(t, actual)

			require.NotEmpty(t, must.NotFail(actualD.Get("queryPlanner")))
			require.NotEmpty(t, must.NotFail(actualD.Get("serverInfo")))

			t.Logf("actual %#v", must.NotFail(actualD.Get("command")))
			commandD := must.NotFail(actualD.Get("command")).(*types.Document)
			for _, key := range tc.expectedCommandKeys {
				assert.NotEmpty(t, commandD.Remove(key), key)
			}
		})
	}
}

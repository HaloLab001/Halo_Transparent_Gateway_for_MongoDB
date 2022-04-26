package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/FerretDB/FerretDB/integration/shareddata"
)

func TestQueryProjection(t *testing.T) {
	t.Parallel()
	providers := []shareddata.Provider{shareddata.Composites}
	ctx, collection := setup(t, providers...)

	_, err := collection.InsertMany(ctx, []any{
		bson.D{
			{"_id", "document-composite-2"},
			{"value", bson.A{
				bson.D{{"field", int32(42)}},
				bson.D{{"field", int32(44)}},
			}},
		},
	})
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		projection  any
		filter      any
		expectedIDs []any
	}{
		"FindProjectionInclusions": {
			projection: bson.D{{"last_name", int32(1)}, {"last_update", true}},
			filter:     bson.D{{"actor_id", int32(28)}},
		},
		"FindProjectionExclusions": {
			projection: bson.D{{"first_name", int32(0)}, {"actor_id", false}},
			filter:     bson.D{{"actor_id", int32(28)}},
		},
		"FindProjectionIDInclusion": {
			projection: bson.D{{"_id", false}, {"actor_id", int32(1)}},
			filter:     bson.D{{"actor_id", int32(28)}},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cursor, err := collection.Find(
				ctx,
				tc.filter,
				options.Find().SetProjection(tc.projection),
				options.Find().SetSort(bson.D{{"_id", 1}}),
			)
			require.NoError(t, err)

			var actual []bson.D
			err = cursor.All(ctx, &actual)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedIDs, CollectIDs(t, actual))
		})
	}
}

func TestQueryProjectionElemMatch(t *testing.T) {
	t.Parallel()
	providers := []shareddata.Provider{shareddata.Composites}
	ctx, collection := setup(t, providers...)

	_, err := collection.InsertMany(ctx, []any{
		bson.D{
			{"_id", "document-composite-2"},
			{"value", bson.A{
				bson.D{{"field", int32(42)}},
				bson.D{{"field", int32(44)}},
			}},
		},
	})
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		projection  any
		expectedIDs []any
	}{
		"ElemMatch": {
			projection: bson.D{{
				"value",
				bson.D{{"$elemMatch", bson.D{{"field", bson.D{{"$eq", 42}}}}}},
			}},
			expectedIDs: []any{
				"document-composite-2",
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cursor, err := collection.Find(
				ctx,
				bson.D{{"_id", "document-composite-2"}},
				options.Find().SetProjection(tc.projection),
				options.Find().SetSort(bson.D{{"_id", 1}}),
			)
			require.NoError(t, err)

			var actual []bson.D
			err = cursor.All(ctx, &actual)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedIDs, CollectIDs(t, actual))
		})
	}
}

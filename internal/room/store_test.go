package room_test

import (
	"context"
	"testing"

	"websockets/internal/room"
	"websockets/internal/test"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestStore_AddToRoom(t *testing.T) {
	testDB := test.NewDynamo(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		prefix := uuid.New().String()
		cleanup := testDB.CreateTables(t, prefix)
		defer cleanup()

		store := room.NewStore(testDB.DB, prefix+"goformation-stack-table", "foo")

		err := store.AddToRoom(ctx, "bar")
		assert.NoError(t, err)

		r, err := store.GetRoom(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []string{"bar"}, r.Connections)
	})

}

func TestStore_RemoveFromRoom(t *testing.T) {
	testDB := test.NewDynamo(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		prefix := uuid.New().String()
		cleanup := testDB.CreateTables(t, prefix)
		defer cleanup()

		store := room.NewStore(testDB.DB, prefix+"goformation-stack-table", "foo")

		err := store.AddToRoom(ctx, "bar")
		assert.NoError(t, err)

		err = store.RemoveFromRoom(ctx, "bar")
		assert.NoError(t, err)

		r, err := store.GetRoom(ctx)
		assert.NoError(t, err)
		assert.Len(t, r.Connections, 0)
	})
}

package memorystorage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	t.Run("create event", func(t *testing.T) {
		event := storage.Event{
			ID:           uuid.New(),
			Title:        "Some title",
			DateTime:     time.Now(),
			EndTime:      time.Now().Add(time.Duration(60)),
			Description:  "Some task description",
			UserID:       1,
			NotifyBefore: time.Now().Add(time.Duration(3600)),
		}
		store := New()
		err := store.Create(&event)
		require.NoError(t, err)
	})
}

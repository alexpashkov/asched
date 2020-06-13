package schedule

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	t.Run("start time cannot be after end", func(t *testing.T) {
		t.Parallel()
		s := NewSchedule()
		now := time.Now()
		require.NoError(t, s.Book(now, now))
		err := s.Book(now, now.Add(-time.Nanosecond))
		require.Error(t, err)
		require.Equal(t, ErrInvalidRange, err)
	})
}

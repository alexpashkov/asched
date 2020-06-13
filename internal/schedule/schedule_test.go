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
		require.NoError(t, s.Book(time.Now(), 1))
		err := s.Book(time.Now(), 0)
		require.Error(t, err)
		require.Equal(t, ErrInvalidDuration, err)
		err = s.Book(time.Now(), -time.Second)
		require.Error(t, err)
		require.Equal(t, ErrInvalidDuration, err)
	})
	t.Run("errors when slot for a given time is booked already", func(t *testing.T) {
		t.Parallel()
		s := NewSchedule()
		bookingTime := time.Now()
		require.NoError(t, s.Book(bookingTime, time.Hour))
		err := s.Book(bookingTime, time.Hour)
		require.Error(t, err)
		require.Equal(t, ErrBooked, err)
	})
	t.Run("errors when trying to book time that's covered by duration of another booking", func(t *testing.T) {
		t.Parallel()
		s := NewSchedule()
		bookingTime := time.Now()
		require.NoError(t, s.Book(bookingTime, time.Hour))
		err := s.Book(bookingTime.Add(time.Hour-1), time.Hour)
		require.Error(t, err)
		require.Equal(t, ErrBooked, err)
	})
	t.Run("allows too book right after another booking", func(t *testing.T) {
		t.Parallel()
		s := NewSchedule()
		bookingTime := time.Now()
		require.NoError(t, s.Book(bookingTime, time.Hour))
		err := s.Book(bookingTime.Add(time.Hour), time.Hour)
		require.NoError(t, err)
	})
}

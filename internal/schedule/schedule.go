package schedule

import (
	"github.com/umpc/go-sortedmap"
	"time"
)

func NewSchedule() *Schedule {
	var s Schedule
	s.bookings = sortedmap.New(0, func(i, j interface{}) bool {
		a, b := i.(Booking), j.(Booking)
		return a.Start.Before(b.Start) &&
			a.Start.Add(a.Duration).Before(b.Start) ||
			a.Start.Add(a.Duration).Equal(b.Start)
	})
	return &s
}

type Schedule struct {
	bookings   *sortedmap.SortedMap
	Validators []Validator
}

func (s *Schedule) Book(start time.Time, duration time.Duration) error {
	if duration <= 0 {
		return ErrInvalidDuration
	}
	booking := Booking{
		Start:    start,
		Duration: duration,
	}
	keys, err := s.bookings.BoundedKeys(booking, Booking{
		Start: booking.Start.Add(booking.Duration),
	})
	if err != nil && err.Error() != "No values found that were equal to or within the given bounds." {
		return err
	}
	if len(keys) > 0 || !s.bookings.Insert(booking.Start, booking) {
		return ErrBooked
	}
	return nil
}

type Validator func(Booking) error

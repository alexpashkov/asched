package schedule

import (
	"github.com/umpc/go-sortedmap"
	"log"
	"time"
)

func NewSchedule() *Schedule {
	var s Schedule
	s.bookings = sortedmap.New(0, func(i, j interface{}) bool {
		a, b := i.(Booking), j.(Booking)
		return a.Start < b.Start && a.Start+a.Duration <= b.Start
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
		Start:    start.UnixNano(),
		Duration: duration.Nanoseconds(),
	}
	keys, err := s.bookings.BoundedKeys(booking, Booking{
		Start: booking.Start + booking.Duration,
	})
	if err != nil && err.Error() != "No values found that were equal to or within the given bounds." {
		return err
	}
	if len(keys) > 0 || !s.bookings.Insert(booking.Start, booking) {
		log.Println(keys, "false")
		return ErrBooked
	}
	return nil
}

type Validator func(Booking) error

type Booking struct {
	Start, Duration int64
}

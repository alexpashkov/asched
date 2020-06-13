package schedule

import "time"

// as a user I can book a slot if its empty

func NewSchedule() *Schedule {
	return new(Schedule)
}

type Schedule struct {

}

func (s *Schedule) Book(start, end time.Time) error {
	if start.After(end) {
		return ErrInvalidRange
	}
	return nil
}


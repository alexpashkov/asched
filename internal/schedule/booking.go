package schedule

import (
	"fmt"
	"time"
)

type Booking struct {
	Start    time.Time
	Duration time.Duration
}

func (b Booking) ID() string {
	return fmt.Sprintf("%s %v", b.Start.Format(time.UnixDate), b.Duration)
}

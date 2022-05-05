package domain

import "time"

type Activity struct {
	CourseName      string
	WeekDay         string
	Times           string
	DaysStr         string
	Days            time.Time
	ComplexName     string
	AvailableSpaces int
}

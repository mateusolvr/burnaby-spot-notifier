package domain

import "time"

type Activity struct {
	CourseName      string
	WeekDay         string
	Times           string
	DateStr         string
	Date            time.Time
	ComplexName     string
	AvailableSpaces int
	ActKeyCache     string
}

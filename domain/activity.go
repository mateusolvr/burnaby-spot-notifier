package domain

import "time"

type Activity struct {
	Name             string
	StartDate        time.Time
	EndDate          time.Time
	RegistrationTime time.Time
	WeekDay          string
	Url              string
	AvailableSpots   string
	ActKeyCache      string
}

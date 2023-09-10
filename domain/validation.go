package domain

import (
	"time"
)

type ValidationService interface {
	ParseTime(timeStr string) time.Time
}

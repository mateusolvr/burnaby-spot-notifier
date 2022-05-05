package domain

import (
	"context"
	"time"
)

type ValidationService interface {
	ValidateActivity(ctx context.Context, activity string) bool
	CleanString(str string) (newStr string)
	CleanFields(courseName, weekDay, times, days, complexName, availableSpaces string) (courseNameCleaned, weekDayCleaned, timesCleaned, daysCleaned, complexNameCleaned, availableSpacesCleaned string)
	ParseDate(dates string) time.Time
}

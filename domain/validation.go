package domain

import (
	"context"
	"time"
)

type ValidationService interface {
	ValidateActivity(ctx context.Context, activity string, actNameConfig string) bool
	CleanString(str string) (newStr string)
	CleanFields(courseName, weekDay, times, date, complexName, availableSpaces string) (courseNameCleaned, weekDayCleaned, timesCleaned, daysCleaned, complexNameCleaned, availableSpacesCleaned string)
	ParseDate(dates string) time.Time
}

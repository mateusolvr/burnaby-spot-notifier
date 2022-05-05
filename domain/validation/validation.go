package validation

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) ValidateActivity(ctx context.Context, activity string) bool {
	r := regexp.MustCompile(`(Volleyball Bonsor)`)
	str := r.FindString(activity)
	if str != "" {
		log.Println(activity)
		return true
	}

	return false
}

func (s *service) CleanString(str string) (newStr string) {

	newStr = strings.ReplaceAll(str, "\n", " ")
	newStr = strings.ReplaceAll(newStr, "View Details", "")
	newStr = strings.TrimSpace(newStr)

	return
}

func (s *service) CleanFields(courseName, weekDay, times, days, complexName, availableSpaces string) (courseNameCleaned, weekDayCleaned, timesCleaned, daysCleaned, complexNameCleaned, availableSpacesCleaned string) {

	courseNameCleaned = s.CleanString(courseName)
	weekDayCleaned = s.CleanString(weekDay)
	timesCleaned = s.CleanString(times)
	daysCleaned = s.CleanString(days)
	complexNameCleaned = s.CleanString(complexName)
	availableSpacesCleaned = s.CleanString(availableSpaces)

	return
}

func (s *service) ParseDate(dateStr string) time.Time {
	date, err := time.Parse("Jan-02-2006", dateStr)
	if err != nil {
		log.Fatal(err)
	}

	return date
}

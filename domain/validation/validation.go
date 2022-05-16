package validation

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/mateusolvr/web-scraper-go/domain"
)

type service struct {
	emailService domain.EmailService
}

func NewService(emailService domain.EmailService) *service {
	return &service{
		emailService: emailService,
	}
}

func (s *service) ValidateActivity(ctx context.Context, activity string, actNameConfig string) bool {
	regexStr := "(" + actNameConfig + ")"
	r := regexp.MustCompile(regexStr)
	str := r.FindString(activity)
	return str != ""
}

func (s *service) CleanString(str string) (newStr string) {

	newStr = strings.ReplaceAll(str, "\n", " ")
	newStr = strings.ReplaceAll(newStr, "View Details", "")
	newStr = strings.TrimSpace(newStr)

	return
}

func (s *service) CleanFields(courseName, weekDay, times, date, complexName, availableSpaces string) (courseNameCleaned, weekDayCleaned, timesCleaned, daysCleaned, complexNameCleaned, availableSpacesCleaned string) {

	courseNameCleaned = s.CleanString(courseName)
	weekDayCleaned = s.CleanString(weekDay)
	timesCleaned = s.CleanString(times)
	daysCleaned = s.CleanString(date)
	complexNameCleaned = s.CleanString(complexName)
	availableSpacesCleaned = s.CleanString(availableSpaces)

	return
}

func (s *service) ParseDate(dateStr string) time.Time {
	date, err := time.Parse("Jan-02-2006", dateStr)
	if err != nil {
		s.emailService.SendErrorEmail(err)
		log.Fatal(err)
	}

	return date
}

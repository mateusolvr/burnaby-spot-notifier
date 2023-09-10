package validation

import (
	"log"
	"time"

	"github.com/mateusolvr/burnaby-spot-notifier/domain"
)

type service struct {
	emailService domain.EmailService
}

func NewService(emailService domain.EmailService) *service {
	return &service{
		emailService: emailService,
	}
}

func (s *service) ParseTime(timeStr string) time.Time {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		s.emailService.SendErrorEmailCache(err)
		log.Fatal(err)
	}

	return parsedTime
}

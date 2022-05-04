package validation

import (
	"context"
	"log"
	"regexp"
	"strings"
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

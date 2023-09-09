package readapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mateusolvr/burnaby-spot-notifier/domain"
)

type service struct {
	validationService domain.ValidationService
	emailService      domain.EmailService
	cacheService      domain.CacheService
	cfg               domain.Config
	activities        []domain.Activity
	postResponseBody  domain.PostResponseBody
}

func NewService(validationService domain.ValidationService, emailService domain.EmailService, cacheService domain.CacheService, cfg domain.Config) *service {
	return &service{
		validationService: validationService,
		emailService:      emailService,
		cacheService:      cacheService,
		cfg:               cfg,
	}
}

func (s *service) Initialize() {
	s.getActivities()
	s.getActivitiesDetails()
	s.checkActivityAvailability()
}

func makeRequest(endpoint string, method string, headers map[string]string, payload *strings.Reader) (int, []byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return 0, nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return resp.StatusCode, body, nil
}

func (s *service) getActivities() {
	url := "https://anc.ca.apm.activecommunities.com/burnaby/rest/onlinecalendar/multicenter/events?="
	method := "POST"
	payload := strings.NewReader(`{
    "calendar_id": 10,
    "center_ids": [
        ` + strconv.Itoa(s.cfg.RecreationalCentreId) + `
    ],
    "display_all": 0,
    "search_start_time": "",
    "search_end_time": "",
    "facility_ids": [],
    "activity_category_ids": [],
    "activity_sub_category_ids": [],
    "activity_ids": [],
    "activity_min_age": 20,
    "activity_max_age": null,
    "event_type_ids": []
}`)
	header := map[string]string{"page_info": "{\"page_number\":1,\"total_records_per_page\":50}"}
	header["Content-Type"] = "application/json"
	status, body, err := makeRequest(url, method, header, payload)
	if status != 200 {
		newErr := errors.New("Status: " + strconv.Itoa(status) + string(body))
		s.emailService.SendErrorEmailCache(newErr)
		log.Fatal(err)
	}
	if err != nil {
		s.emailService.SendErrorEmailCache(err)
		log.Fatal(err)
	}

	if err = json.Unmarshal([]byte(body), &s.postResponseBody); err != nil {
		s.emailService.SendErrorEmailCache(err)
		log.Fatal(err)
	}

}

func (s *service) getActivitiesDetails() {
	currentTime := time.Now()
	for _, act := range s.postResponseBody.Body.CenterEvents[0].Events {
		if !strings.Contains(act.Title, s.cfg.ActivityName) {
			continue
		}

		startTime := s.validationService.ParseTime(act.StartTime)
		// Get only future activities
		if currentTime.After(startTime) {
			continue
		}
		url := fmt.Sprintf("https://anc.ca.apm.activecommunities.com/burnaby/rest/onlinecalendar/activity-details/%d", act.EventItemID)
		method := "GET"
		status, body, err := makeRequest(url, method, nil, strings.NewReader(``))
		if status != 200 {
			newErr := errors.New("Status: " + strconv.Itoa(status) + string(body))
			s.emailService.SendErrorEmailCache(newErr)
			log.Fatal(err)
		}
		if err != nil {
			s.emailService.SendErrorEmailCache(err)
			log.Fatal(err)
		}

		var getResponseBody domain.GetResponseBody
		if err = json.Unmarshal([]byte(body), &getResponseBody); err != nil {
			s.emailService.SendErrorEmailCache(err)
			log.Fatal(err)
		}
		registrationTimeStr := getResponseBody.Body.ActivityDetail.RegistrationDate.EnrollmentDate[0].InternetTime
		registrationTime := s.validationService.ParseTime(registrationTimeStr)
		// Get only activities that registration has started
		if registrationTime.After(currentTime) {
			continue
		}

		newActivity := domain.Activity{
			Name:             act.Title,
			StartDate:        startTime,
			EndDate:          s.validationService.ParseTime(act.EndTime),
			RegistrationTime: s.validationService.ParseTime(getResponseBody.Body.ActivityDetail.RegistrationDate.EnrollmentDate[0].InternetTime),
			WeekDay:          startTime.Weekday().String(),
			Url:              act.ActivityDetailURL,
			AvailableSpots:   getResponseBody.Body.ActivityDetail.SpaceStatus,
			ActKeyCache:      strconv.Itoa(act.EventItemID),
		}

		s.activities = append(s.activities, newActivity)

	}

}

func (s *service) checkActivityAvailability() {
	var availableActivities []domain.Activity
	for _, act := range s.activities {
		if act.AvailableSpots != "Closed" && act.AvailableSpots != "Full" {
			availableActivities = append(availableActivities, act)
		}
		if act.AvailableSpots == "Full" {
			_, err := s.cacheService.DelKey(act.ActKeyCache)
			if err != nil && s.cfg.Redis.Enabled {
				s.emailService.SendErrorEmailCache(err)
				log.Fatal(err)
			}
		}
	}

	if len(availableActivities) > 0 {
		s.emailService.SendEmailCache(availableActivities)
	} else {
		log.Printf("%d activities were found!", len(availableActivities))
	}
}

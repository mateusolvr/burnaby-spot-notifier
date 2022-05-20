package email

import (
	"log"
	"net/smtp"
	"strconv"

	"github.com/mateusolvr/web-scraper-go/domain"
)

type service struct {
	cacheService domain.CacheService
	cfg          domain.Config
}

func NewService(cacheService domain.CacheService, cfg domain.Config) *service {
	return &service{
		cacheService: cacheService,
		cfg:          cfg,
	}
}

func (s *service) sendEmail(activities []domain.Activity) {

	from := s.cfg.Email.From
	pass := s.cfg.Email.Pass
	to := s.cfg.Email.To

	subject := "Available Activity - Burnaby"

	headers := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	htmlBody := s.buildHtmlBody(activities)

	msg := []byte(headers + mime + htmlBody)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Fatal(err)
	}
}

func (s *service) buildHtmlBody(activities []domain.Activity) string {
	var tableLines string
	for _, act := range activities {
		tableLines += `<tr>
		<td style="width: 20.7562%; border: 1px solid rgb(96, 96, 96);">` + act.CourseName + `</td>
		<td style="width: 14.9691%; border: 1px solid rgb(96, 96, 96);">` + act.WeekDay + `</td>
		<td style="width: 11.7525%; border: 1px solid rgb(96, 96, 96);">` + act.Times + `</td>
		<td style="width: 18.5258%; border: 1px solid rgb(96, 96, 96);">` + act.Date.Format("Jan-02-2006") + `</td>
		<td style="width: 20.2715%; border: 1px solid rgb(96, 96, 96);">` + act.ComplexName + `</td>
		<td style="width: 13.6791%; border: 1px solid rgb(96, 96, 96);">` + strconv.Itoa(act.AvailableSpaces) + `</td>
		</tr>`
	}
	htmlBody := `<p>Hello!</p>
	<p>The following activities are available:</p>
	<table style="width: 100%; border-collapse: collapse; border: 1px solid rgb(96, 96, 96);">
		<thead>
			<tr>
				<th style="width: 20.7562%; border: 1px solid rgb(96, 96, 96);">Course Name</th>
				<th style="width: 14.9691%; border: 1px solid rgb(96, 96, 96);">Week Day</th>
				<th style="width: 11.7525%; border: 1px solid rgb(96, 96, 96);">Time</th>
				<th style="width: 18.5258%; border: 1px solid rgb(96, 96, 96);">Date</th>
				<th style="width: 20.2715%; border: 1px solid rgb(96, 96, 96);">Complex Name</th>
				<th style="width: 13.6791%; border: 1px solid rgb(96, 96, 96);">Available Spaces</th>
			</tr>
		</thead>
		<tbody>
			` + tableLines + `
		</tbody>
	</table>
	<p>Enjoy! 😊</p>
	<p>🏐⚽️🏈🏋️🏊</p>`

	return htmlBody
}

func (s *service) SendErrorEmail(err error) {
	from := s.cfg.Email.From
	pass := s.cfg.Email.Pass
	to := s.cfg.Email.To

	log.Println("Email was sent for this error.")

	subject := "ERROR: Available Activity - Burnaby"

	htmlBody := `<p>Hello!</p>
	<p>The following error occurred while trying to fetch the activities:</p>
	<p>` + err.Error() + `</p>
	<p>Please, contact the administrator!</p>
	<p>🏐⚽️🏈🏋️🏊</p>`

	headers := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(headers + mime + htmlBody)

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Fatal(err)
		return
	}
}

func (s *service) SendEmailCache(activities []domain.Activity) {
	filteredActivities, err := s.cacheService.GetActivitiesWithoutCache(activities)
	if err != nil {
		s.SendErrorEmail(err)
		log.Fatal(err)
	}

	if len(filteredActivities) == 0 {
		log.Printf("%d activities were found but no email was sent!", len(activities))
		return
	}

	err = s.cacheService.AddActivitiesCache(filteredActivities)
	if err != nil {
		s.SendErrorEmail(err)
		log.Fatal(err)
	}

	s.sendEmail(filteredActivities)
	log.Printf("%d activities were found and %d activities were sent in the email!",
		len(activities), len(filteredActivities))
}

func (s *service) SendErrorEmailCache(err error) {
	errInCache, cacheErr := s.cacheService.CheckErrorCache(err)
	if cacheErr != nil {
		s.SendErrorEmail(cacheErr)
		log.Fatal(cacheErr)
	}
	if errInCache {
		log.Println("Email was not sent for this error due to cache.")
		return
	}

	s.cacheService.SetKey(err.Error(), "true", s.cfg.Redis.ExpireMinutes)
	s.SendErrorEmail(err)
}
